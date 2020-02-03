package main

import (
	"log"

	"github.com/gocql/gocql" 
)

var cluster *gocql.ClusterConfig
var session *gocql.Session


func populateItemTable() {
	query := `
		BEGIN BATCH
		INSERT INTO pouches (userId, itemName, itemStartingPrice) values ('qwe', 'item rgd gfhwr', 0);
		INSERT INTO pouches (userId, itemName, itemStartingPrice) values ('qwe', 'item 4rwkgfsfd', 0);
		INSERT INTO pouches (userId, itemName, itemStartingPrice) values ('qwe', 'item fdmvbrwrf', 0);
		INSERT INTO pouches (userId, itemName, itemStartingPrice) values ('asd', 'item 5teddfkdd', 0);
		APPLY BATCH;
	`

	if err := session.Query(query).Exec(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// connect to the cluster
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ = cluster.CreateSession()
	defer session.Close()

	populateItemTable()
}
