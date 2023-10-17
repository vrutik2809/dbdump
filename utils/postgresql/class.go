package postgresql

import (
	"database/sql"
	"fmt"

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

func (p *postgresql) FetchAllTables(schema string) (*sql.Rows, error) {
	rows, err := p.client.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = $1", schema)

	return rows, err
}

func (p *postgresql) FetchAllRows(table string) (*sql.Rows, error) {
	rows, err := p.client.Query("SELECT * FROM " + table)

	return rows, err
}
