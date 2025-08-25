package service

import (
	"encoding/json"
	"st-novel-go/src/novel/dao"
	"st-novel-go/src/novel/dto"
)

func UpdateNovelJSONField(novelID string, userID uint, fieldName string, data interface{}) error {
	return dao.UpdateNovelJSONField(novelID, userID, fieldName, data)
}

func getNovelJSONField(novelID string, userID uint, fieldName string, target interface{}) error {
	novel, err := dao.FindNovelByID(novelID, userID)
	if err != nil {
		return err
	}

	var rawData json.RawMessage
	switch fieldName {
	case "settings_data":
		rawData = json.RawMessage(novel.SettingsData)
	case "plot_custom_data":
		rawData = json.RawMessage(novel.PlotCustomData)
	case "analysis_custom_data":
		rawData = json.RawMessage(novel.AnalysisCustomData)
	case "others_custom_data":
		rawData = json.RawMessage(novel.OthersCustomData)
	}

	if len(rawData) > 0 && string(rawData) != "null" {
		if err := json.Unmarshal(rawData, target); err != nil {
			return err
		}
	}
	return nil
}

func GetSettingsData(novelID string, userID uint) (interface{}, error) {
	var result []dto.TreeNodeDTO
	err := getNovelJSONField(novelID, userID, "settings_data", &result)
	if result == nil {
		return make([]dto.TreeNodeDTO, 0), err
	}
	return result, err
}

func GetPlotCustomData(novelID string, userID uint) (interface{}, error) {
	var result []dto.ItemNodeDTO
	err := getNovelJSONField(novelID, userID, "plot_custom_data", &result)
	if result == nil {
		return make([]dto.ItemNodeDTO, 0), err
	}
	return result, err
}

func GetAnalysisCustomData(novelID string, userID uint) (interface{}, error) {
	var result []dto.ItemNodeDTO
	err := getNovelJSONField(novelID, userID, "analysis_custom_data", &result)
	if result == nil {
		return make([]dto.ItemNodeDTO, 0), err
	}
	return result, err
}

func GetOthersCustomData(novelID string, userID uint) (interface{}, error) {
	var result []dto.ItemNodeDTO
	err := getNovelJSONField(novelID, userID, "others_custom_data", &result)
	if result == nil {
		return make([]dto.ItemNodeDTO, 0), err
	}
	return result, err
}
