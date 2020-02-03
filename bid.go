package main

import (
	"log"
	"fmt"
	"time"
	"math/rand"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
)

var cluster *gocql.ClusterConfig
var session *gocql.Session

func bid(bidderId string, bidAmount int) {
	auctionId := "1"

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
	bidFinished := make(chan bool)
	auctionActive := true
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	maxRandomBidAmount := 10000

	// connect to the cluster
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ = cluster.CreateSession()
	defer session.Close()


	go func() {
		time.Sleep(400 * time.Millisecond)
		auctionActive = false
	}()

	go func(refreshFinished chan bool) {
		for auctionActive {
				refresh()
				time.Sleep(1 * time.Second)
		}
		refreshFinished <- true
	}(refreshFinished)

	go func(bidFinished chan bool) {
		id := uuid.Must(uuid.NewV4()).String()

		for auctionActive {
				bidAmount := r1.Intn(maxRandomBidAmount)

				bid(id, bidAmount)
				fmt.Println(id, "bidding", bidAmount)
				time.Sleep(200 * time.Millisecond)
		}
		bidFinished <- true
	}(bidFinished)

	go func(bidFinished chan bool) {
		id := uuid.Must(uuid.NewV4()).String()
		
		for auctionActive {
				bidAmount := r1.Intn(maxRandomBidAmount)

				bid(id, bidAmount)
				fmt.Println(id, "bidding", bidAmount)
				time.Sleep(200 * time.Millisecond)
		}
		bidFinished <- true
	}(bidFinished)

	go func(bidFinished chan bool) {
		id := uuid.Must(uuid.NewV4()).String()
		
		for auctionActive {
				bidAmount := r1.Intn(maxRandomBidAmount)

				bid(id, bidAmount)
				fmt.Println(id, "bidding", bidAmount)
				time.Sleep(200 * time.Millisecond)
		}
		bidFinished <- true
	}(bidFinished)

	go func(bidFinished chan bool) {
		id := uuid.Must(uuid.NewV4()).String()
		
		for auctionActive {
				bidAmount := r1.Intn(maxRandomBidAmount)

				bid(id, bidAmount)
				fmt.Println(id, "bidding", bidAmount)
				time.Sleep(200 * time.Millisecond)
		}
		bidFinished <- true
	}(bidFinished)

	<- refreshFinished
	<- bidFinished

	fmt.Println("end")

	// var winnerBid float64
	// var winnerId string
	// if err := session.Query(`SELECT MAX(bid) from bids`).Consistency(gocql.Quorum).Scan(&winnerBid); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Max Bid:", winnerBid)

	// var bidderId string
	// var bid float64
	results := session.Query(`SELECT * FROM bids`).Iter()
	m := &map[string]interface{}{}
	ret := []map[string]interface{}{}
	// var winnerId string
	// var winnerBid float64
	var maxBid float64
	var winnerId string
	var currentBid float64
	var currentBidderId string

	for results.MapScan(*m) {
		ret := append(ret, *m)
		m = &map[string]interface{}{}
		currentBid, _ = strconv.ParseFloat(fmt.Sprint(ret[0]["bid"]), 64)
		currentBidderId = fmt.Sprint(ret[0]["bidderid"])
		if (currentBid > maxBid) {
			maxBid = currentBid
			winnerId = currentBidderId
		}
		// fmt.Println(ret[0]["bid"])
	}

	fmt.Println(winnerId, maxBid)
	

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