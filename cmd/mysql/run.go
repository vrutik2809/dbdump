package mysql

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/vrutik2809/dbdump/utils"
	"github.com/vrutik2809/dbdump/utils/mysql"
)

const (
	JSON string = "json"
	CSV  string = "csv"
	TSV  string = "tsv"
)

func isOutputValid(output string) bool {
	validTypes := []string{JSON, CSV, TSV}	
	for _, validType := range validTypes {
		if output == validType {
			return true
		}
	}
	return false
}

func getFileExtension(output string) string {
	switch output {
		case JSON:
			return ".json"
		case CSV:
			return ".csv"
		case TSV:
			return ".tsv"
		default:
			return ""
	}
}

func writeToOutputFile(output string, result []map[string]interface{}, tableName string) error {
	switch output {
		case JSON:
			return utils.MapArrayToJSONFile(result, tableName + getFileExtension(output))
		case CSV:
			return utils.MapArrayToCSVFile(result, tableName + getFileExtension(output))
		case TSV:
			return utils.MapArrayToTSVFile(result, tableName + getFileExtension(output))
	}
	return nil
}

func dumpTable(wg *sync.WaitGroup, msq *mysql.MySQL, table string, output string) {
	defer wg.Done()
	fmt.Println("dumping table: ", table)
	rows, err := msq.FetchAllRows(table)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeToOutputFile(output, rows, table); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetUint("port")
	dbName, _ := cmd.Flags().GetString("db-name")
	outputDir, _ := cmd.Flags().GetString("dir")
	dumpTables, _ := cmd.Flags().GetStringSlice("tables")
	excludeTables, _ := cmd.Flags().GetStringSlice("exclude-tables")
	output, _ := cmd.Flags().GetString("output")

	if !isOutputValid(output) {
		log.Fatal("invalid output type. valid types are: json, csv, tsv")
	}

	os.RemoveAll(outputDir)
	os.Mkdir(outputDir, 0777)
	os.Chdir(outputDir)

	msq := mysql.NewMySQL(username, password, host, port, dbName)

	if err := msq.Connect(); err != nil {
		log.Fatal(err)
	}
	defer msq.Close()

	if err := msq.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL | uri: ", msq.GetURI())

	tables, err := msq.FetchTables(dumpTables, excludeTables)

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, table := range tables {
		wg.Add(1)
		go dumpTable(&wg, msq, table, output)
	}

	wg.Wait()

	fmt.Println("dumped tables successfully")

}
