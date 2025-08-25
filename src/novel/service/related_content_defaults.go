package service

import (
	"encoding/json"
	"gorm.io/datatypes"
	"st-novel-go/src/novel/dto"
)

func GetDefaultSettingsData() datatypes.JSON {
	settings := []dto.TreeNodeDTO{
		{
			"id":    "setting",
			"title": "设定",
			"type":  "root",
			"icon":  "fa-solid fa-book-journal-whills",
			"children": []dto.TreeNodeDTO{
				{
					"id":    "characters",
					"title": "角色",
					"type":  "group",
					"icon":  "fa-regular fa-folder",
					"children": []dto.TreeNodeDTO{
						{
							"id":         "characters_overview",
							"title":      "角色总览",
							"type":       "characters_overview",
							"icon":       "fa-solid fa-users",
							"isOverview": true,
							"isReadOnly": true,
							"content":    "<h1>角色总览</h1><p>此分类下暂无内容，请添加条目。</p>",
						},
					},
				},
				{
					"id":    "locations",
					"title": "地点",
					"type":  "group",
					"icon":  "fa-regular fa-folder",
					"children": []dto.TreeNodeDTO{
						{
							"id":         "locations_overview",
							"title":      "地点总览",
							"type":       "locations_overview",
							"icon":       "fa-solid fa-map-location-dot",
							"isOverview": true,
							"isReadOnly": true,
							"content":    "<h1>地点总览</h1><p>此分类下暂无内容，请添加条目。</p>",
						},
					},
				},
				{
					"id":    "worldview",
					"title": "世界观",
					"type":  "group",
					"icon":  "fa-regular fa-folder",
					"children": []dto.TreeNodeDTO{
						{
							"id":         "worldview_overview",
							"title":      "世界观总览",
							"type":       "worldview_overview",
							"icon":       "fa-solid fa-book-atlas",
							"isOverview": true,
							"isReadOnly": true,
							"content":    "<h1>世界观总览</h1><p>此分类下暂无内容，请添加条目。</p>",
						},
					},
				},
				{
					"id":    "items",
					"title": "物品",
					"type":  "group",
					"icon":  "fa-regular fa-folder",
					"children": []dto.TreeNodeDTO{
						{
							"id":         "items_overview",
							"title":      "物品总览",
							"type":       "items_overview",
							"icon":       "fa-solid fa-box-archive",
							"isOverview": true,
							"isReadOnly": true,
							"content":    "<h1>物品总览</h1><p>此分类下暂无内容，请添加条目。</p>",
						},
					},
				},
			},
		},
	}
	data, _ := json.Marshal(settings)
	return data
}

func GetDefaultCustomData() datatypes.JSON {
	data, _ := json.Marshal([]dto.ItemNodeDTO{})
	return data
}
