package postgresql

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/vrutik2809/dbdump/utils/postgresql"
)

func run(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetUint("port")
	dbName, _ := cmd.Flags().GetString("db-name")
	outputDir, _ := cmd.Flags().GetString("dir")
	schema, _ := cmd.Flags().GetString("schema")
	dumpTables, _ := cmd.Flags().GetStringSlice("tables")
	excludeTables, _ := cmd.Flags().GetStringSlice("exclude-tables")

	os.RemoveAll(outputDir)
	os.Mkdir(outputDir, 0777)
	os.Chdir(outputDir)

	pg := postgresql.NewPostgreSQL(username, password, host, port, dbName)

	if err := pg.Connect(); err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	if err := pg.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to DB successfully | uri: ", pg.GetURI())

	tables, err := pg.FetchTables(schema, dumpTables, excludeTables)
	if err != nil {
		log.Fatal(err)
	}
	defer tables.Close()

	for tables.Next() {
		var tableName string
		err := tables.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("dumping table: ", tableName)

		rows, err := pg.FetchAllRows(tableName)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var result []map[string]interface{}

		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				log.Fatal(err)
			}

			values := make([]interface{}, len(columns))
			for i := range values {
				values[i] = new(interface{})
			}

			err = rows.Scan(values...)
			if err != nil {
				log.Fatal(err)
			}

			rowData := make(map[string]interface{})
			for i, column := range columns {
				val := *(values[i].(*interface{}))
				rowData[column] = val
			}

			result = append(result, rowData)
		}

		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		jsonData, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(tableName + ".json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.Write(jsonData)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = tables.Err(); err != nil {
		log.Fatal(err)
	}
}
