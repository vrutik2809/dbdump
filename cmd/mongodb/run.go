package mongodb

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cheggaaa/pb/v3"
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

func dumpToFile(bsonDArray []bson.D, bar *pb.ProgressBar, collection string, output string) error {
	filename := collection + getFileExtension(output)
	switch output {
		case JSON:
			return utils.BsonDArrayToJsonFile(bsonDArray, bar, filename)
		case BSON:
			return utils.BsonDArrayToFile(bsonDArray, bar, filename)
		case GZIP:
			return utils.BsonDArrayToGzipFile(bsonDArray, bar, filename)
	default:
		return nil
	}
}

func dumpCollection(wg *sync.WaitGroup, bar *pb.ProgressBar, mongo *mongodb.MongoDB, collection string, output string) {
	defer wg.Done()

	bsonDArray, err := mongo.FetchAllDocuments(collection)
	if err != nil {
		log.Fatal(err)
	}

	if err := dumpToFile(bsonDArray, bar, collection, output); err != nil {
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
	isSRV, _ := cmd.Flags().GetBool("srv")
	collections, _ := cmd.Flags().GetStringSlice("collections")
	collectionsExclude, _ := cmd.Flags().GetStringSlice("exclude-collections")
	output, _ := cmd.Flags().GetString("output")
	testMode, _ := cmd.Flags().GetBool("test-mode")

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

	var wg sync.WaitGroup
	
	bars := utils.GetBars(collections, "collection",testMode)

	var barPool *pb.Pool

	if testMode {
		barPool, err = nil, nil
	} else {
		barPool, err = pb.StartPool(bars...)
	}

	if err != nil {
		log.Fatal(err)
	}

	for idx, collection := range collections {
		wg.Add(1)
		go dumpCollection(&wg, bars[idx], mongo, collection, output)
	}

	wg.Wait()

	if barPool != nil {
		barPool.Stop()
	}

	fmt.Println("dumped collections successfully")
}
