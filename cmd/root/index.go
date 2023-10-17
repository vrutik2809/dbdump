package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vrutik2809/dbdump/cmd/mongodb"
	"github.com/vrutik2809/dbdump/cmd/postgresql"
)

var RootCmd = &cobra.Command{
	Use:   "dbdump",
	Short: "CLI tool for dumping various databases",
	Long:  "CLI tool for dumping various databases in various formats like tsv,csv,json,gzip",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("welcome to dbdump")
	},
}

// the init() function in Go is automatically invoked when the package is initialized.
func init() {
	RootCmd.AddCommand(mongodb.MongodbCmd)
	RootCmd.AddCommand(postgresql.PostgresqlCmd)
}