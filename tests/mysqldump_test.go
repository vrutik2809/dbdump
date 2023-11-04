package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/vrutik2809/dbdump/utils"
	"github.com/vrutik2809/dbdump/utils/mysql"
)

type MySQLDumpTestSuite struct {
	suite.Suite
	mysql     *mysql.MySQL
	tables []map[string]interface{}
}

func (suite *MySQLDumpTestSuite) SetupSuite() {
	mysql := mysql.NewMySQL("root", "admin123", "localhost", 3307, "test")
	if err := mysql.Connect(); err != nil {
		suite.T().Error(err)
	}
	if err := mysql.Ping(); err != nil {
		suite.T().Error(err)
	}

	suite.T().Log("Connected to MySQL | uri: " + mysql.GetURI())

	tables := []map[string]interface{}{
		{
			"name":        "users",
			"createQuery": "CREATE TABLE users (id INT NOT NULL, name VARCHAR (50) NOT NULL, email VARCHAR (50) NOT NULL);",
			"insertQuery": "INSERT INTO users (id, name, email) VALUES (1, 'user1', 'user1@gmail.com'), (2, 'user2', 'user2@gmail.com'), (3, 'user3','user3@gmail.com');",
			"dropQuery":   "DROP TABLE users;",
		},
		{
			"name":        "photos",
			"createQuery": "CREATE TABLE photos (id INT NOT NULL, url VARCHAR (50) NOT NULL, user_id INT NOT NULL);",
			"insertQuery": "INSERT INTO photos (id, url, user_id) VALUES (1, 'https://user1.com', 1), (2, 'https://user2.com', 2), (3, 'https://user3.com', 3);",
			"dropQuery":   "DROP TABLE photos;",
		},
		{
			"name":        "profiles",
			"createQuery": "CREATE TABLE profiles (id INT NOT NULL, user_id INT NOT NULL, age INT NOT NULL);",
			"insertQuery": "INSERT INTO profiles (id, user_id, age) VALUES (1, 1, 20), (2, 2, 21), (3, 3, 22);",
			"dropQuery":   "DROP TABLE profiles;",
		},
	}

	for _, table := range tables {
		if err := mysql.ExecuteQuery(table["createQuery"].(string)); err != nil {
			suite.T().Error(err)
		}
		if err := mysql.ExecuteQuery(table["insertQuery"].(string)); err != nil {
			suite.T().Error(err)
		}
	}

	suite.mysql = mysql
	suite.tables = tables

}

func (suite *MySQLDumpTestSuite) TearDownSuite() {
	defer suite.mysql.Close()
	for _, table := range suite.tables {
		if err := suite.mysql.ExecuteQuery(table["dropQuery"].(string)); err != nil {
			suite.T().Error(err)
		}
	}
	suite.T().Log("Disconnected from PostgreSQL | uri: " + suite.mysql.GetURI())
}

func (suite *MySQLDumpTestSuite) TestJsonDump() {
	dumpDir := "../dump/"
	ext := ".json"
	cmd := exec.Command("go", "run", "../main.go", "mysql", "-u", "root", "--password", "admin123", "--host", "localhost", "-p", "3307", "-d", "test", "--dir", dumpDir)
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mysql dump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, table := range suite.tables {
		tableName := table["name"].(string)
		_, err = os.Stat(dumpDir + tableName + ext)
		assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+tableName+ext)
	}
}

func (suite *MySQLDumpTestSuite) TestCSVDump() {
	dumpDir := "../dump/"
	ext := ".csv"
	cmd := exec.Command("go", "run", "../main.go", "mysql", "-u", "root", "--password", "admin123", "--host", "localhost", "-p", "3307", "-d", "test", "--dir", dumpDir, "-o", "csv")
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mysql dump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, table := range suite.tables {
		tableName := table["name"].(string)
		_, err = os.Stat(dumpDir + tableName + ext)
		assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+tableName+ext)
	}
}

func (suite *MySQLDumpTestSuite) TestTSVDump() {
	dumpDir := "../dump/"
	ext := ".tsv"
	cmd := exec.Command("go", "run", "../main.go", "mysql", "-u", "root", "--password", "admin123", "--host", "localhost", "-p", "3307", "-d", "test", "--dir", dumpDir, "-o", "tsv")
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mysql dump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, table := range suite.tables {
		tableName := table["name"].(string)
		_, err = os.Stat(dumpDir + tableName + ext)
		assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+tableName+ext)
	}
}

func (suite *MySQLDumpTestSuite) TestTablesDump() {
	dumpDir := "../dump/"
	ext := ".json"
	dumpTables := []string{"users", "photos"}
	cmd := exec.Command("go", "run", "../main.go", "mysql", "-u", "root", "--password", "admin123", "--host", "localhost", "-p", "3307", "-d", "test", "--dir", dumpDir, "-t", strings.Join(dumpTables, ","))
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mysql dump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, table := range suite.tables {
		tableName := table["name"].(string)
		if utils.Contains(dumpTables, tableName) {
			_, err = os.Stat(dumpDir + tableName + ext)
			assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+tableName+ext)
		} else {
			_, err = os.Stat(dumpDir + tableName + ext)
			assert.Error(suite.T(), err, "dump directory contains file with name: "+tableName+ext)
		}
	}
}

func (suite *MySQLDumpTestSuite) TestExcludeTablesDump() {
	dumpDir := "../dump/"
	ext := ".json"
	excludeTables := []string{"photos"}
	cmd := exec.Command("go", "run", "../main.go", "mysql", "-u", "root", "--password", "admin123", "--host", "localhost", "-p", "3307", "-d", "test", "--dir", dumpDir, "-e", strings.Join(excludeTables, ","))
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mysql dump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, table := range suite.tables {
		tableName := table["name"].(string)
		if utils.Contains(excludeTables, tableName) {
			_, err = os.Stat(dumpDir + tableName + ext)
			assert.Error(suite.T(), err, "dump directory contains file with name: "+tableName+ext)
		} else {
			_, err = os.Stat(dumpDir + tableName + ext)
			assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+tableName+ext)
		}
	}
}

func (suite *MySQLDumpTestSuite) TestAggregatedDump() {
	dumpDir := "../dump/"
	ext := ".json"
	dumpTables := []string{"users", "photos"}
	excludeTables := []string{"photos"}
	cmd := exec.Command("go", "run", "../main.go", "mysql", "-u", "root", "--password", "admin123", "--host", "localhost", "-p", "3307", "-d", "test", "--dir", dumpDir, "-t", strings.Join(dumpTables, ","), "-e", strings.Join(excludeTables, ","))
	err := cmd.Run()
	assert.NoError(suite.T(), err, "mysql dump failed")

	_, err = os.Stat(dumpDir)
	assert.NoError(suite.T(), err, "dump directory not created")

	for _, table := range suite.tables {
		tableName := table["name"].(string)
		if utils.Contains(dumpTables, tableName) && !utils.Contains(excludeTables, tableName) {
			_, err = os.Stat(dumpDir + tableName + ext)
			assert.NoError(suite.T(), err, "dump directory does not contain file with name: "+tableName+ext)
		} else {
			_, err = os.Stat(dumpDir + tableName + ext)
			assert.Error(suite.T(), err, "dump directory contains file with name: "+tableName+ext)
		}
	}
}

func TestMySQLDumpTestSuite(t *testing.T) {
	suite.Run(t, new(MySQLDumpTestSuite))
}
