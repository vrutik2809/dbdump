package mongodb

import (
	"github.com/spf13/cobra"
)

var MongodbCmd = &cobra.Command{
	Use:   "mongodb",
	Short: "command for dumping mongodb database",
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func init(){
	MongodbCmd.Flags().StringP("username", "u", "", "username of the database")
	MongodbCmd.Flags().String("password", "", "password of the database")
	MongodbCmd.Flags().String("host", "localhost", "host of the database")
	MongodbCmd.Flags().UintP("port", "p", 0, "port of the database")
	MongodbCmd.Flags().StringP("db-name", "d", "", "name of the database")
	MongodbCmd.Flags().String("dir", "dump", "name of the output directory")
	MongodbCmd.Flags().Bool("srv", false, "use SRV connection format")
	MongodbCmd.Flags().StringSliceP("collections", "c", []string{}, "name of the collections to dump")
	MongodbCmd.Flags().StringSliceP("exclude-collections", "e", []string{}, "name of the collections to exclude")
	MongodbCmd.Flags().StringP("output", "o", "json", "output format of the dump (json, bson, gzip)")
	MongodbCmd.Flags().Bool("test-mode",false, "run the command in test mode (not showing progress bar)")

	MongodbCmd.MarkFlagRequired("db-name")
	MongodbCmd.MarkFlagsRequiredTogether("username", "password")
}