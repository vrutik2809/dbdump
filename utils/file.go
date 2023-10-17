package utils

import (
	"os"
	"compress/gzip"

	"go.mongodb.org/mongo-driver/bson"
)

func BsonDArrayToJsonFile(bsonDArray []bson.D, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("[")

	first := true

	for _, bsonD := range bsonDArray {
		res, err := BsonDToJson(bsonD)
		if err != nil {
			return err
		}

		if !first {
			file.WriteString(",")
		}
		first = false
		file.WriteString(string(res))
	}
	file.WriteString("]")

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