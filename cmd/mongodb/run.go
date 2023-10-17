package mongodb

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/vrutik2809/dbdump/utils"
	"github.com/vrutik2809/dbdump/utils/mongodb"
)

const (
	JSON string = "json"
	BSON string = "bson"
	GZIP string = "gzip"
)

func isOutputTypeValid(output string) bool {
	validTypes := []string{JSON, BSON, GZIP}
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
		case BSON:
			return ".bson"
		case GZIP:
			return ".gz"
		default:
			return ""
	}
}

func dumpToFile(bsonDArray []bson.D, collection string, output string) error {
	filename := collection + getFileExtension(output)
	switch output {
		case JSON:
			return utils.BsonDArrayToJsonFile(bsonDArray, filename)
		case BSON:
			return utils.BsonDArrayToFile(bsonDArray, filename)
		case GZIP:
			return utils.BsonDArrayToGzipFile(bsonDArray, filename)
		default:
			return nil
	}
}


func run(cmd *cobra.Command, args []string) {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetUint("port")
	dbName, _ := cmd.Flags().GetString("db-name")
	outputDir, _ := cmd.Flags().GetString("dir")
	isSRV, _ := cmd.Flags().GetBool("srv")
	collections, _ := cmd.Flags().GetStringSlice("collections")
	collectionsExclude, _ := cmd.Flags().GetStringSlice("exclude-collections")
	output, _ := cmd.Flags().GetString("output")

	if !isOutputTypeValid(output) {
		log.Fatal("invalid output type. valid types are: json, bson, gzip")
	}
	
	mongo := mongodb.NewMongoDB(username, password, host, port, dbName, isSRV)

	if err := mongo.Connect(); err != nil {
		log.Fatal(err)
	}

	defer mongo.Close()

	if err := mongo.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB | uri: " + mongo.GetURI())

	filter := mongodb.CollectionFilter(collections, collectionsExclude)

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

		if err := dumpToFile(bsonDArray, collection, output); err != nil {
			log.Fatal(err)
		}
	}
}