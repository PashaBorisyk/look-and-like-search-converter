package indexer

import (
	"log"
	"look-and-like-search-converter/converter"
	"look-and-like-search-converter/models"
	"look-and-like-search-converter/web"
	"time"
)

var convertTime = time.Now().Format(time.RFC3339)

func IndexProduct(product models.Product) error {

	product.ID = converter.MongoIDToString(product.ID)
	log.Println("Processing product with id ", product.ID)
	content, err := converter.ConvertInterfaceToMap(product)
	if err != nil {
		log.Println("Will not continue processing due to error while converting to map: ", err)
		return err
	}
	convertTimeToString(&content)
	content["@search.action"] = "upload"
	delete(content,"id")

	searchUploadModel := models.NewSearchUploadModel(&content)
	err = web.UploadModelToIndex(searchUploadModel)

	return err

}

func convertTimeToString(content *map[string]interface{}) {
	switch metaInformation := (*content)["metaInformation"].(type) {
	case map[string]interface{}:
		switch metaInformation["insertDate"].(type) {
		case map[string]interface{}:
			log.Println("Old time format of metaInformation.insertDate detected. Will perform convert")
			metaInformation["insertTime"] = convertTime
		}
	}

}
