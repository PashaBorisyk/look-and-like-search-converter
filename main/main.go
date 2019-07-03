package main

import (
	"log"
	"look-and-like-search-converter/indexer"
	"look-and-like-search-converter/logger"
	"look-and-like-search-converter/models"
	"look-and-like-search-converter/mongoClient"
	"look-and-like-search-converter/queue"
	"sync"
)

var totalSuccess = 0
var totalWarnings = 0
var totalFails = 0

var collection = mongoClient.GetCollection("products")

func processProduct(product models.Product, err error) error {

	if err != nil {
		log.Println("Error not nil at the start of processProduct loop")
		return err
	}

	mongoID := product.ID
	err = indexer.IndexProduct(product)
	if err != nil {
		totalFails++
		log.Println("Error while uploading document with id ", product.ID, " occurred: ", err)
	} else {
		err = collection.SetIndexed(mongoID)
		if err != nil {
			totalWarnings++
			log.Println("Error: document was indexed, but 'metaInformation.indexed' field was not set to 'true' : ", err)
			return err
		}
		totalSuccess++
		log.Println("Document with id ", mongoID, "successfully indexed (index id is ", product.ID, ")")
	}
	log.Println("Total successes: ", totalSuccess, "; Total warnings: ", totalWarnings, "; Total fails: ", totalFails)

	return err
}

func init() {
	logger.Init()
}

func main() {
	//RunIndex()

	log.Println("Search Converter starting...")
	wg := sync.WaitGroup{}
	wg.Add(1)
	queue.InitConsumer(&wg, func(id string) {
		product, err := collection.GetByID(id)
		if err != nil {
			log.Println("Error getting product by ID :", err)
		} else {
			err = processProduct(*product, err)
			if err != nil {
				log.Println("Error processing product :", err)
			}
		}
	})

	log.Println("Search converter started")
	wg.Wait()
	log.Println("Search converter finished")
}

func RunIndex() {
	err := collection.GetNotIndexedDocuments(processProduct)
	if err != nil {
		log.Println("Foreach loop finished with error: ", err)
	}
}
