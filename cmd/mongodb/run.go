package mongodb

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/vrutik2809/dbdump/utils/mongodb"
	"github.com/vrutik2809/dbdump/utils"
)

func run(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetUint("port")
	dbName, _ := cmd.Flags().GetString("db-name")
	outputDir, _ := cmd.Flags().GetString("dir")
	isSRV, _ := cmd.Flags().GetBool("srv")
	collections, _ := cmd.Flags().GetStringSlice("collections")
	
	mongo := mongodb.NewMongoDB(username, password, host, port, dbName, isSRV)

	if err := mongo.Connect(); err != nil {
		log.Fatal(err)
	}

	defer mongo.Close()

	if err := mongo.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB | uri: " + mongo.GetURI())

	filter := mongodb.CollectionFilter(collections)

	collections, err := mongo.FetchCollections(filter)

	if err != nil {
		log.Fatal(err)
	}

	os.RemoveAll(outputDir)
	os.Mkdir(outputDir, 0777)
	os.Chdir(outputDir)

	for _, collection := range collections {
		fmt.Println("dumping collection:", collection)

		bsonDArray, err := mongo.FetchAllDocuments(collection)
		if err != nil {
			log.Fatal(err)
		}

		if err := utils.BsonDArrayToJsonFile(bsonDArray, collection+".json"); err != nil {
			log.Fatal(err)
		}
	}
}