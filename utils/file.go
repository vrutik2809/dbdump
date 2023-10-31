package utils

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func BsonDArrayToJsonFile(bsonDArray []bson.D, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	result := []map[string]interface{}{}

	count := len(bsonDArray)

	bar := pb.New(count).SetTemplateString(GetBarTemplate("collection", GetName(filename))).SetMaxWidth(80).Start()

	for _, bsonD := range bsonDArray {
		res, err := BsonDToMap(bsonD)
		if err != nil {
			return err
		}
		result = append(result, res)
		bar.Increment()
	}

	bar.Finish()

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

	count := len(bsonDArray)

	bar := pb.New(count).SetTemplateString(GetBarTemplate("collection", GetName(filename))).SetMaxWidth(80).Start()

	for _, bsonD := range bsonDArray {
		data, err := bson.Marshal(bsonD)
		if err != nil {
			return err
		}
		file.Write(data)
		bar.Increment()
	}

	bar.Finish()

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

	count := len(bsonDArray)

	bar := pb.New(count).SetTemplateString(GetBarTemplate("collection", GetName(filename))).SetMaxWidth(80).Start()

	for _, bsonD := range bsonDArray {
		data, err := bson.Marshal(bsonD)
		if err != nil {
			return err
		}
		gzipWriter.Write(data)
		bar.Increment()
	}

	bar.Finish()

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
			str := InterfaceToString(row[key])
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
			str := InterfaceToString(row[key])
			rowData = append(rowData, str)
		}
		file.WriteString(stringToTSVRow(rowData))
	}

	return nil
}
