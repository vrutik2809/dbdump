package utils

import (
	"encoding/json"
	"database/sql"

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

func SqlRowToString(rows *sql.Rows) ([]string, error) {
	defer rows.Close()

	var result []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		result = append(result, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func SqlRowToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, column := range columns {
			val := *(values[i].(*interface{}))
			switch v := val.(type) {
				case []byte:
					rowData[column] = string(v)
				default:
					rowData[column] = v
			}
		}

		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}