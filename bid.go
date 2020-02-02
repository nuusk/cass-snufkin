package main

import (
	"log"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

var cluster *gocql.ClusterConfig
var session *gocql.Session

func bid(bidAmount int) {
	auctionId := "1"
	bidderId := "1"

	if err := session.Query(`INSERT INTO bids (auctionId, time, bid, bidderId) VALUES (?, ?, ?, ?)`,
		auctionId, time.Now(), float64(bidAmount), bidderId).Exec(); err != nil {
		log.Fatal(err)
	}
}

func refresh() {
	var maxBidAmount float64

	if err := session.Query(`SELECT MAX(bid) from bids`).Consistency(gocql.One).Scan(&maxBidAmount); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Max Bid Amount:", maxBidAmount)
}

func main() {
	refreshFinished := make(chan bool)
	// connect to the cluster
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ = cluster.CreateSession()
	defer session.Close()

	bid(4)
	bid(3)
	bid(7)
	bid(1)

	go func(refreshFinished chan bool) {
		for {
				time.Sleep(1 * time.Second)
				refresh()
		}
		refreshFinished <- true
	}(refreshFinished)

	<- refreshFinished

	
	// var id gocql.UUID
	// var text string

	// /* Search for a specific set of records whose 'timeline' column matches
	//  * the value 'me'. The secondary index that we created earlier will be
	//  * used for optimizing the search */
	// if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
	// 	"me").Consistency(gocql.One).Scan(&id, &text); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Tweet:", id, text)

	// // list all tweets
	// iter := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()

	// if err := iter.Close(); err != nil {
	// 	log.Fatal(err)
	// }
}