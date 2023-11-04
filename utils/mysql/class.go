package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vrutik2809/dbdump/utils"
)

type MySQL struct {
	username string
	password string
	host     string
	port     uint
	dbName   string
	client   *sql.DB
}

func NewMySQL(username string, password string, host string, port uint, dbName string) *MySQL {
	return &MySQL{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dbName:   dbName,
		client:   nil,
	}
}

func (m *MySQL) GetURI() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.username, m.password, m.host, m.port, m.dbName)
}

func (m *MySQL) Connect() error {
	db, err := sql.Open("mysql", m.GetURI())
	if err != nil {
		return err
	}
	m.client = db
	return nil
}

func (m *MySQL) Close() error {
	return m.client.Close()
}

func (m *MySQL) Ping() error {
	return m.client.Ping()
}

func (m *MySQL) FetchTables(dumpTables []string, excludeTables []string) ([]string, error) {
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

func (m *MySQL) FetchAllRows(tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := m.client.Query(query)
	if err != nil {
		return nil, err
	}
	return utils.SqlRowToMap(rows)
}

func (p *MySQL) ExecuteQuery(query string) error {
	_, err := p.client.Exec(query)
	return err
}
