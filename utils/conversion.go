package utils

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func BsonDToMap(b bson.D) (map[string]interface{}, error){
	data, err := bson.Marshal(b)
	if err != nil {
		return nil,err
	}

	var mp map[string]interface{}
	err = bson.Unmarshal(data,&mp)
	if err != nil {
		return nil, err
	}

	return mp,err
}

func BsonDToJson(b bson.D) (string, error){
	mp,err := BsonDToMap(b)
	if err != nil {
		return "",err
	}

	res,err := json.Marshal(mp)
	return string(res),err
}