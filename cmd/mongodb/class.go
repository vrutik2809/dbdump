package mongodb

import (
	"fmt"
	"os/exec"
)

type mongoDB struct {
	username string
	password string
	host string
	port uint
	dbName string
}

func (mongo mongoDB) print(){
	fmt.Printf("username: %s, password: %s, host: %s, port: %d, db-name: %s\n", mongo.username, mongo.password, mongo.host, mongo.port, mongo.dbName)
}

func (mongo mongoDB) dump() (string,error){
	fmt.Println("dumping the database ....")

	cmd := exec.Command("./bin/mongodump", "--help")

	output, err := cmd.CombinedOutput()
	
	return string(output), err
}