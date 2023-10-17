package postgresql

import (
	"github.com/spf13/cobra"
)

var PostgresqlCmd = &cobra.Command{
	Use:   "pg",
	Short: "command for dumping postgresql database",
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func init() {
	PostgresqlCmd.Flags().StringP("username", "u", "postgres", "username of the database")
	PostgresqlCmd.Flags().String("password", "123456", "password of the database")
	PostgresqlCmd.Flags().String("host", "localhost", "host of the database")
	PostgresqlCmd.Flags().UintP("port", "p", 5432, "port of the database")
	PostgresqlCmd.Flags().StringP("db-name", "d", "", "name of the database")
	PostgresqlCmd.Flags().String("dir", "dump", "name of the output directory")
	PostgresqlCmd.Flags().StringP("schema", "s", "public", "name of the schema")

	PostgresqlCmd.MarkFlagRequired("db-name")
}
