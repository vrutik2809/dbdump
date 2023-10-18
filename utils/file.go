package utils

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"strconv"
	"strings"

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

func MapArrayToJSONFile(mp []map[string]interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(mp, "", "\t")
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	
	return err
}

func interfaceToString(itf interface{}) string {
	switch v := itf.(type) {
		case int64:
			return strconv.FormatInt(v, 10)
		case string:
			return v
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(v)
		default:
			return ""
	}
}

func stringToCSVRow(str []string) string {
	return `"` + strings.Join(str, `","`) + `"` + "\n"
}

func MapArrayToCSVFile(mp []map[string]interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	header := []string{}
	for key, _ := range mp[0] {
		header = append(header, key)
	}
	file.WriteString(stringToCSVRow(header))

	for _, row := range mp {
		rowData := []string{}
		for _, key := range header {
			str := interfaceToString(row[key])
			rowData = append(rowData, str)
		}
		file.WriteString(stringToCSVRow(rowData))
	}

	return nil
}

func stringToTSVRow(str []string) string {
	return `"` + strings.Join(str, "\"\t\"") + `"` + "\n"
}

func MapArrayToTSVFile(mp []map[string]interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	header := []string{}
	for key, _ := range mp[0] {
		header = append(header, key)
	}
	file.WriteString(stringToTSVRow(header))

	for _, row := range mp {
		rowData := []string{}
		for _, key := range header {
			str := interfaceToString(row[key])
			rowData = append(rowData, str)
		}
		file.WriteString(stringToTSVRow(rowData))
	}

	return nil
}
