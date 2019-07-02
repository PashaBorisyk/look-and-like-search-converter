package models

type SearchUploadModel struct {
	Value []map[string]interface{} `json:"value"`
}

func NewSearchUploadModel(content *map[string]interface{}) *SearchUploadModel {
	value := make([]map[string]interface{}, 1)
	value[0] = *content
	return &SearchUploadModel{
		Value: value,
	}
}
