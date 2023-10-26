package postgresql

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/vrutik2809/dbdump/utils"
)

type PostgreSQL struct {
	username string
	password string
	host     string
	port     uint
	dbName   string
	client   *sql.DB
}

func NewPostgreSQL(username string, password string, host string, port uint, dbName string) *PostgreSQL {
	return &PostgreSQL{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dbName:   dbName,
		client:   nil,
	}
}

func (p *PostgreSQL) GetURI() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", p.host, p.port, p.username, p.password, p.dbName)
}

func (p *PostgreSQL) Connect() error {
	uri := p.GetURI()
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return err
	}
	p.client = db
	return nil
}

func (p *PostgreSQL) Close() error {
	return p.client.Close()
}

func (p *PostgreSQL) Ping() error {
	return p.client.Ping()
}

func (p *PostgreSQL) FetchTables(schema string, dumpTables []string, excludeTables []string) ([]string, error) {
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

func (p *PostgreSQL) FetchAllRows(table string) ([]map[string]interface{}, error) {
	rows, err := p.client.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}

	return utils.SqlRowToMap(rows)
}


func (p *PostgreSQL) ExecuteQuery(query string) error {
	_, err := p.client.Exec(query)
	return err
}

// func (p *PostgreSQL) CreateTable(tableQuery string) error {
// 	_, err := p.client.Exec(tableQuery)
// 	return err
// }

// func (p *PostgreSQL) DropTable(table string) error {
// 	_, err := p.client.Exec("DROP TABLE " + table)
// 	return err
// }

// func (p *PostgreSQL) getColumnNames(table string) ([]string, error) {
// 	rows, err := p.client.Query("SELECT column_name FROM information_schema.columns WHERE table_name = '" + table + "'")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return utils.SqlRowToString(rows)
// }

// func (p *PostgreSQL) getInsertQuery(table string, rows []interface{}) (string, error) {
// 	columns, err := p.getColumnNames(table)
// 	if err != nil {
// 		return "", err
// 	}
// 	values := []string{}
// 	for _, row := range rows {
// 		values = append(values, utils.InterfaceToSQLInsertRowString(row))
// 	}
// 	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), strings.Join(values, ",")), nil
// }

// func (p *PostgreSQL) InsertRow(table string, row []interface{}) error {
// 	query, err := p.getInsertQuery(table, row)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = p.client.Exec(query)
// 	return err
// }
