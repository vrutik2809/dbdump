package mongodb

import (
	"fmt"

	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetUint("port")
	dbName, _ := cmd.Flags().GetString("db-name")
	
	mongo := mongoDB{
		username: username,
		password: password,
		host: host,
		port: port,
		dbName: dbName,
	}

	output, err := mongo.dump()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(output)
}