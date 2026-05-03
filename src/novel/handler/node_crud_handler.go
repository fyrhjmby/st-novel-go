// node_crud_handler.go — 为 settings/custom-plot/custom-analysis/custom-others 提供单项 CRUD
// 这些数据在后端以 JSON 树数组形式存储在 novel 的对应字段中。
// 前端使用单项 CRUD（POST/PATCH/DELETE）操作，本 handler 负责将单项操作转换为对 JSON 树的读写。
package handler

import (
	"encoding/json"
	"fmt"
	"st-novel-go/src/middleware"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
	"st-novel-go/src/novel/service"
	"st-novel-go/src/utils"

	"github.com/gin-gonic/gin"
)

// fieldName -> getter/setter 映射
type nodeCrudConfig struct {
	FieldName string
	Getter    func(string, uint) (interface{}, error)
}

var nodeCrudConfigs = map[string]nodeCrudConfig{
	"settings":        {FieldName: "settings_data", Getter: service.GetSettingsData},
	"custom-plot":     {FieldName: "plot_custom_data", Getter: service.GetPlotCustomData},
	"custom-analysis": {FieldName: "analysis_custom_data", Getter: service.GetAnalysisCustomData},
	"custom-others":   {FieldName: "others_custom_data", Getter: service.GetOthersCustomData},
}

// findNodeInTree 递归在 TreeNodeDTO 树中查找指定 ID 的节点，返回节点指针和父节点切片路径
func findNodeInTree(nodes []dto.TreeNodeDTO, targetID string) (*dto.TreeNodeDTO, []dto.TreeNodeDTO) {
	for i := range nodes {
		if result, path := findNodeRecursive(&nodes[i], targetID, nil); result != nil {
			return result, path
		}
	}
	return nil, nil
}

func findNodeRecursive(node *dto.TreeNodeDTO, targetID string, parent *dto.TreeNodeDTO) (*dto.TreeNodeDTO, []dto.TreeNodeDTO) {
	if node.ID == targetID {
		var path []dto.TreeNodeDTO
		if parent != nil {
			path = []dto.TreeNodeDTO{*parent}
		}
		return node, path
	}
	if node.Children != nil {
		for i := range node.Children {
			if result, path := findNodeRecursive(&node.Children[i], targetID, node); result != nil {
				return result, path
			}
		}
	}
	return nil, nil
}

// getNodeCrudConfig 从 URL path 中提取 type 参数
func getNodeCrudConfig(c *gin.Context) (nodeCrudConfig, string, error) {
	nodeType := c.Param("type")
	cfg, ok := nodeCrudConfigs[nodeType]
	if !ok {
		return nodeCrudConfig{}, "", fmt.Errorf("unsupported node type: %s", nodeType)
	}
	return cfg, nodeType, nil
}

// CreateNodeHandler POST /api/nodes/:type
// Body: { novelId, parentId?, title, content?, type, ... }
func CreateNodeHandler(c *gin.Context) {
	cfg, _, err := getNodeCrudConfig(c)
	if err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	var newNode dto.TreeNodeDTO
	if err := c.ShouldBindJSON(&newNode); err != nil {
		utils.FailWithBadRequest(c, "Invalid node data: "+err.Error())
		return
	}

	novelID := c.Query("novelId")
	if novelID == "" {
		var bodyMap map[string]interface{}
		json.Unmarshal(c.MustGetRawData(), &bodyMap)
		if id, ok := bodyMap["novelId"].(string); ok {
			novelID = id
		}
	}
	if novelID == "" {
		utils.FailWithBadRequest(c, "novelId is required")
		return
	}

	// 读取现有树，添加新节点，写回
	rawData, err := cfg.Getter(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to read data: "+err.Error())
		return
	}

	tree, err := convertToTree(rawData)
	if err != nil {
		utils.Fail(c, "Failed to parse existing data: "+err.Error())
		return
	}

	tree = append(tree, newNode)
	if err := service.UpdateNovelJSONField(novelID, userClaims.UserID, cfg.FieldName, tree); err != nil {
		utils.Fail(c, "Failed to save: "+err.Error())
		return
	}

	utils.Success(c, newNode)
}

// UpdateNodeHandler PATCH /api/nodes/:type/:nodeId
// Body: { title?, content?, ... }
func UpdateNodeHandler(c *gin.Context) {
	cfg, _, err := getNodeCrudConfig(c)
	if err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	nodeID := c.Param("nodeId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	novelID := c.Query("novelId")
	if novelID == "" {
		utils.FailWithBadRequest(c, "novelId query parameter is required")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.FailWithBadRequest(c, "Invalid update data: "+err.Error())
		return
	}

	rawData, err := cfg.Getter(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to read data: "+err.Error())
		return
	}

	tree, err := convertToTree(rawData)
	if err != nil {
		utils.Fail(c, "Failed to parse existing data: "+err.Error())
		return
	}

	node, _ := findNodeInTree(tree, nodeID)
	if node == nil {
		utils.Fail(c, "Node not found")
		return
	}

	if title, ok := updates["title"].(string); ok {
		node.Title = title
	}
	if content, ok := updates["content"].(string); ok {
		node.Content = content
	}

	if err := service.UpdateNovelJSONField(novelID, userClaims.UserID, cfg.FieldName, tree); err != nil {
		utils.Fail(c, "Failed to save: "+err.Error())
		return
	}

	utils.Success(c, node)
}

// DeleteNodeHandler DELETE /api/nodes/:type/:nodeId
func DeleteNodeHandler(c *gin.Context) {
	cfg, _, err := getNodeCrudConfig(c)
	if err != nil {
		utils.FailWithBadRequest(c, err.Error())
		return
	}

	nodeID := c.Param("nodeId")
	claims, _ := c.Get(middleware.UserClaimsKey)
	userClaims := claims.(*utils.Claims)

	novelID := c.Query("novelId")
	if novelID == "" {
		utils.FailWithBadRequest(c, "novelId query parameter is required")
		return
	}

	rawData, err := cfg.Getter(novelID, userClaims.UserID)
	if err != nil {
		utils.Fail(c, "Failed to read data: "+err.Error())
		return
	}

	tree, err := convertToTree(rawData)
	if err != nil {
		utils.Fail(c, "Failed to parse existing data: "+err.Error())
		return
	}

	// 从树中删除节点
	newTree := removeNodeFromTree(tree, nodeID)
	if len(newTree) == len(tree) {
		utils.Fail(c, "Node not found")
		return
	}

	if err := service.UpdateNovelJSONField(novelID, userClaims.UserID, cfg.FieldName, newTree); err != nil {
		utils.Fail(c, "Failed to save: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Node deleted")
}

func removeNodeFromTree(nodes []dto.TreeNodeDTO, targetID string) []dto.TreeNodeDTO {
	result := make([]dto.TreeNodeDTO, 0, len(nodes))
	for _, node := range nodes {
		if node.ID == targetID {
			continue
		}
		if node.Children != nil {
			node.Children = removeNodeFromTree(node.Children, targetID)
		}
		result = append(result, node)
	}
	return result
}

// convertToTree 将 interface{} 转换为 []dto.TreeNodeDTO
func convertToTree(rawData interface{}) ([]dto.TreeNodeDTO, error) {
	bytes, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}
	var tree []dto.TreeNodeDTO
	if err := json.Unmarshal(bytes, &tree); err != nil {
		return nil, err
	}
	return tree, nil
}

// verifyUserOwnership 通过 novelId 校验用户权限
func verifyNodeCrudOwnership(novelID string, userID uint) error {
	_, err := dao.FindNovelByID(novelID, userID)
	return err
}
