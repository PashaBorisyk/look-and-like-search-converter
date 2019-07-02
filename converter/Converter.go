package converter

import (
	"encoding/json"
	"fmt"
	"log"
)

func ConvertInterfaceToMap(val interface{}) (map[string]interface{}, error) {
	result, err := json.Marshal(val)
	content := make(map[string]interface{})
	err = json.Unmarshal(result, &content)
	if err != nil {
		log.Println("Error decoding interface to map: ", err)
		return nil, err
	}
	return content, nil
}

func MongoIDToString(id interface{}) string {
	return fmt.Sprintf("%v", id)[10:34]
}