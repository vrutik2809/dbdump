package postgresql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type postgresql struct {
	username string
	password string
	host     string
	port     uint
	dbName   string
	client   *sql.DB
}

func NewPostgreSQL(username string, password string, host string, port uint, dbName string) *postgresql {
	return &postgresql{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dbName:   dbName,
		client:   nil,
	}
}

func (p *postgresql) GetURI() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", p.host, p.port, p.username, p.password, p.dbName)
}

func (p *postgresql) Connect() error {
	uri := p.GetURI()
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return err
	}
	p.client = db
	return nil
}

func (p *postgresql) Close() error {
	return p.client.Close()
}

func (p *postgresql) Ping() error {
	return p.client.Ping()
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

func (p *postgresql) FetchTables(schema string, dumpTables []string, excludeTables []string) ([]string, error) {
	in := "'" + strings.Join(dumpTables, "','") + "'"
	notIn := "'" + strings.Join(excludeTables, "','") + "'"
	query := fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", schema)
	if len(dumpTables) > 0 {
		query = fmt.Sprintf("%s AND table_name IN (%s)", query, in)
	}
	if len(excludeTables) > 0 {
		query = fmt.Sprintf("%s AND table_name NOT IN (%s)", query, notIn)
	}
	rows, err := p.client.Query(query)
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
			rowData[column] = val
		}

		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *postgresql) FetchAllRows(table string) ([]map[string]interface{}, error) {
	rows, err := p.client.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}

	return sqlRowToMap(rows)
}
