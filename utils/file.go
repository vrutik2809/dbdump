package utils

import (
	"os"
	"encoding/json"
	"compress/gzip"

	"go.mongodb.org/mongo-driver/bson"
)

func BsonDArrayToJsonFile(bsonDArray []bson.D, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	result := []map[string]interface{}{}

	for _, bsonD := range bsonDArray {
		res, err := BsonDToMap(bsonD)
		if err != nil {
			return err
		}
		result = append(result, res)
	}

	jsonData, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		return err
	}

	file.Write(jsonData)

	return nil
}

func BsonDArrayToFile(bsonDArray []bson.D, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, bsonD := range bsonDArray {
		data, err := bson.Marshal(bsonD)
		if err != nil {
			return err
		}
		file.Write(data)
	}

	return nil
}

func BsonDArrayToGzipFile(bsonDArray []bson.D, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	for _, bsonD := range bsonDArray {
		data, err := bson.Marshal(bsonD)
		if err != nil {
			return err
		}
		gzipWriter.Write(data)
	}

	return nil
}