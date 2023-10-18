package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vrutik2809/dbdump/utils"
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
	return utils.SqlRowToString(rows)
}

func (m *mySql) FetchAllRows(tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := m.client.Query(query)
	if err != nil {
		return nil, err
	}
	return utils.SqlRowToMap(rows)
}
