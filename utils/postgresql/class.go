package postgresql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/vrutik2809/dbdump/utils"
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
	return utils.SqlRowToString(rows)
}

func (p *postgresql) FetchAllRows(table string) ([]map[string]interface{}, error) {
	rows, err := p.client.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}

	return utils.SqlRowToMap(rows)
}
