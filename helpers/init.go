package main

import (
	"log"

	"github.com/gocql/gocql" 
)

var cluster *gocql.ClusterConfig
var session *gocql.Session

func initBidTable() {
	query := "CREATE TABLE bids (auctionId text, time timestamp, bid double, bidderId text, PRIMARY KEY((auctionId, bidderId), time))"
	if err := session.Query(query).Exec(); err != nil {
		log.Fatal(err)
	}
}

func initTransactionTable() {
	query := "CREATE TABLE transactions(userId text, time timestamp, amount double, PRIMARY KEY(userId, time));"
	if err := session.Query(query).Exec(); err != nil {
		log.Fatal(err)
	}
}

func initItemTable() {
	query := "CREATE TABLE items(userId text, itemName text, itemStartingPrice double, PRIMARY KEY(userId, itemName))"
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

	initBidTable()
	initTransactionTable()
	initItemTable()
}