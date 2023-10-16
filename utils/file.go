package utils

import (
	"os"

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