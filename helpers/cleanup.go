package main

import (
	"log"

	"github.com/gocql/gocql"
)

var cluster *gocql.ClusterConfig
var session *gocql.Session

func dropBidTable() {
	query := "DROP TABLE bids"
	if err := session.Query(query).Exec(); err != nil {
		log.Fatal(err)
	}
}

func dropItemTable() {
	query := "DROP TABLE items"
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

	dropBidTable()
	dropItemTable()
}
