package mysql

import (
	"github.com/spf13/cobra"
)

var MysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "command for dumping MySQL database",
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func init() {
	MysqlCmd.Flags().StringP("username", "u", "root", "username of the database")
	MysqlCmd.Flags().String("password", "root", "password of the database")
	MysqlCmd.Flags().String("host", "localhost", "host of the database")
	MysqlCmd.Flags().UintP("port", "p", 3306, "port of the database")
	MysqlCmd.Flags().StringP("db-name", "d", "", "name of the database")
	MysqlCmd.Flags().String("dir", "dump", "name of the output directory")
	MysqlCmd.Flags().StringSliceP("tables", "t", []string{}, "name of the tables to dump")
	MysqlCmd.Flags().StringSliceP("exclude-tables", "e", []string{}, "name of the tables to exclude")
	MysqlCmd.Flags().StringP("output","o","json","output format of the dump (json, csv, tsv)")

	MysqlCmd.MarkFlagRequired("db-name")
}
