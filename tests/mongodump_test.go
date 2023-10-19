package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/vrutik2809/dbdump/utils"
	"github.com/vrutik2809/dbdump/utils/mongodb"
)

type MongoDumpTestSuite struct {
	suite.Suite
	mongo            *mongodb.MongoDB
	dummyCollections []string
}

func (suite *MongoDumpTestSuite) SetupSuite() {
	mongo := mongodb.NewMongoDB("admin", "admin123", "localhost", 27019, "test", false)
	if err := mongo.Connect(); err != nil {
		suite.Error(err)
	}

	if err := mongo.Ping(); err != nil {
		suite.Error(err)
	}

	suite.T().Log("Connected to MongoDB | uri: " + mongo.GetURI())

	dummyCollections := []string{"users", "photos", "profiles", "roles"}
	dummyDocuments := []interface{}{
		map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
		map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
		map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}

	for _, collection := range dummyCollections {
		if err := mongo.CreateCollection(collection); err != nil {
			suite.Error(err)
		}
	}

	for _, collection := range dummyCollections {
		if err := mongo.InsertDocuments(collection, dummyDocuments); err != nil {
			suite.Error(err)
		}
	}

	suite.mongo = mongo
	suite.dummyCollections = dummyCollections
}

func (suite *MongoDumpTestSuite) TearDownSuite() {
	defer suite.mongo.Close()
	for _, collection := range suite.dummyCollections {
		if err := suite.mongo.DropCollection(collection); err != nil {
			suite.Error(err)
		}
	}
	suite.T().Log("Disconnected from MongoDB | uri: " + suite.mongo.GetURI())
}

func (suite *MongoDumpTestSuite) TestJsonDump() {
	dumpDir := "../dump/"
	ext := ".json"
	cmd := exec.Command("go", "run", "../main.go", "mongodb", "--username", "admin", "--password", "admin123", "--host", "localhost", "--port", "27019", "--db-name", "test", "--dir", dumpDir)
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mongodump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, collection := range suite.dummyCollections {
		_, err = os.Stat(dumpDir + collection + ext)
		assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+collection+ext)
	}
}
func (suite *MongoDumpTestSuite) TestBsonDump() {
	dumpDir := "../dump/"
	ext := ".bson"
	cmd := exec.Command("go", "run", "../main.go", "mongodb", "--username", "admin", "--password", "admin123", "--host", "localhost", "--port", "27019", "--db-name", "test", "--dir", dumpDir, "--output", "bson")
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mongodump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, collection := range suite.dummyCollections {
		_, err = os.Stat(dumpDir + collection + ext)
		assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+collection+ext)
	}
}
func (suite *MongoDumpTestSuite) TestGZIPDump() {
	dumpDir := "../dump/"
	ext := ".gz"
	cmd := exec.Command("go", "run", "../main.go", "mongodb", "--username", "admin", "--password", "admin123", "--host", "localhost", "--port", "27019", "--db-name", "test", "--dir", dumpDir, "--output", "gzip")
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mongodump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, collection := range suite.dummyCollections {
		_, err = os.Stat(dumpDir + collection + ext)
		assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+collection+ext)
	}
}

func (suite *MongoDumpTestSuite) TestCollectionDump() {
	dumpDir := "../dump/"
	ext := ".json"
	dumpCollections := []string{"users", "photos"}
	cmd := exec.Command("go", "run", "../main.go", "mongodb", "--username", "admin", "--password", "admin123", "--host", "localhost", "--port", "27019", "--db-name", "test", "--dir", dumpDir, "--collections", strings.Join(dumpCollections, ","))
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mongodump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, collection := range suite.dummyCollections {
		if utils.Contains(dumpCollections, collection) {
			_, err = os.Stat(dumpDir + collection + ext)
			assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+collection+ext)
		} else {
			_, err = os.Stat(dumpDir + collection + ext)
			assert.Error(suite.T(), err, "dump directory contains file with name: "+collection+ext)
		}
	}
}

func (suite *MongoDumpTestSuite) TestExcludeCollectionDump() {
	dumpDir := "../dump/"
	ext := ".json"
	excludeCollections := []string{"users"}
	cmd := exec.Command("go", "run", "../main.go", "mongodb", "-u", "admin", "--password", "admin123", "--host", "localhost", "-p", "27019", "-d", "test", "--dir", dumpDir, "-e", strings.Join(excludeCollections, ","))
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mongodump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, collection := range suite.dummyCollections {
		if utils.Contains(excludeCollections, collection) {
			_, err = os.Stat(dumpDir + collection + ext)
			assert.Error(suite.T(), err, "dump directory contains file with name: "+collection+ext)
		} else {
			_, err = os.Stat(dumpDir + collection + ext)
			assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+collection+ext)
		}
	}
}

func (suite *MongoDumpTestSuite) TestAggregatedDump() {
	dumpDir := "../dump/"
	ext := ".json"
	dumpCollections := []string{"users", "photos"}
	excludeCollections := []string{"users"}
	cmd := exec.Command("go", "run", "../main.go", "mongodb", "-u", "admin", "--password", "admin123", "--host", "localhost", "-p", "27019", "-d", "test", "--dir", dumpDir, "-c", strings.Join(dumpCollections, ","), "-e", strings.Join(excludeCollections, ","))
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mongodump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, collection := range suite.dummyCollections {
		if utils.Contains(dumpCollections, collection) && !utils.Contains(excludeCollections, collection) {
			_, err = os.Stat(dumpDir + collection + ext)
			assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+collection+ext)
		} else {
			_, err = os.Stat(dumpDir + collection + ext)
			assert.Error(suite.T(), err, "dump directory contains file with name: "+collection+ext)
		}
	}
}

func TestMongoDumpTestSuite(t *testing.T) {
	suite.Run(t, new(MongoDumpTestSuite))
}
