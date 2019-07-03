package mongoClient

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"look-and-like-search-converter/models"
	"time"
)

var database *mongo.Database
var collectionsMap map[string]*Collection

type Collection struct {
	collection *mongo.Collection
}

func init() {

	mongoUri := "mongodb://look-and-like-test:unE8DZr3T7yA6SLDPjknaT8Bj0MzLD4O4604EDq0OE44Lv9BxAslwWXTqLvJFzvqLCBoCDshGgUUJuKoahpT6w==@look-and-like-test.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	productDatabaseName := "look-and-like-test"

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal("Unable to create Mongo client with address: ", mongoUri, "; ", err)
		return
	}

	ctx := createContext()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Unable to connect to Mongo database with address: ", mongoUri, "; ", err)
		return
	}

	database = client.Database(productDatabaseName)
	collectionsMap = make(map[string]*Collection)
}

func GetCollection(name string) *Collection {
	if collectionsMap[name] == nil {
		collectionsMap[name] = &Collection{database.Collection(name)}
	}
	return collectionsMap[name]
}

func decodeMultipleResult(cursor *mongo.Cursor, foreach func(product models.Product, err error) error) error {
	ctx := createContext()
	var product models.Product
	for cursor.Next(ctx) {
		err := cursor.Decode(&product)
		if err != nil {
			log.Println("Unable to decode document: ", err)
		}
		//product.ID = cursor.Current.Lookup("_id").ObjectID()
		if foreach(product, err) != nil {
			println("Error in foreach loop: ", err)
		}
	}
	_ = cursor.Close(ctx)
	err := cursor.Err()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (holder *Collection) GetNotIndexedDocuments(foreach func(message models.Product, err error) error) error {
	ctx := createContext()
	queryResult, err := holder.collection.Find(ctx, bson.M{"metaInformation.indexed": nil})
	if err != nil {
		log.Println("Unable to get result from FindUnsentByReceiverUserID: ", err)
	}
	defer queryResult.Close(ctx)
	return decodeMultipleResult(queryResult, foreach)
}

func (holder *Collection) GetByID(id string) (product *models.Product,err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error decoding ObjectID from hex string ",id," :",err)
		return nil, err
	}
	ctx := createContext()
	queryResult := holder.collection.FindOne(ctx,bson.M{"_id":objectID})
	err = queryResult.Decode(product)
	return product,err
}

func (holder *Collection) SetIndexed(productID interface{}) error {
	ctx := createContext()
	result := holder.collection.FindOneAndUpdate(ctx,
		bson.M{"_id": productID},
		bson.D{{"$set", bson.D{{"metaInformation.indexed", true}}}})
	_, err := result.DecodeBytes()
	if err != nil {
		log.Println("Error updating Products document with id", productID," : ", err)
		return err
	}
	return nil
}

func createContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 24*time.Hour)
	return ctx
}
