package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type mySql struct {
	username string
	password string
	host     string
	port     uint
	dbName   string
	client   *sql.DB
}

func NewMySQL(username string, password string, host string, port uint, dbName string) *mySql {
	return &mySql{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dbName:   dbName,
		client:   nil,
	}
}

func (m *mySql) GetURI() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.username, m.password, m.host, m.port, m.dbName)
}

func (m *mySql) Connect() error {
	db, err := sql.Open("mysql", m.GetURI())
	if err != nil {
		return err
	}
	m.client = db
	return nil
}

func (m *mySql) Close() error {
	return m.client.Close()
}

func (m *mySql) Ping() error {
	return m.client.Ping()
}

func sqlRowToString(rows *sql.Rows) ([]string, error) {
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

func (m *mySql) FetchTables(dumpTables []string, excludeTables []string) ([]string, error) {
	in := "'" + strings.Join(dumpTables, "','") + "'"
	notIn := "'" + strings.Join(excludeTables, "','") + "'"
	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", m.dbName)
	if len(dumpTables) > 0 {
		query = fmt.Sprintf("%s AND table_name IN (%s)", query, in)
	}
	if len(excludeTables) > 0 {
		query = fmt.Sprintf("%s AND table_name NOT IN (%s)", query, notIn)
	}
	rows, err := m.client.Query(query)
	if err != nil {
		return nil, err
	}
	return sqlRowToString(rows)
}

func sqlRowToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
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

func (m *mySql) FetchAllRows(tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := m.client.Query(query)
	if err != nil {
		return nil, err
	}
	return sqlRowToMap(rows)
}