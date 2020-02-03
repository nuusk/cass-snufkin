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

func transaction(userId string, amount float64) {
	if err := session.Query(`INSERT INTO transactions (userId, time, amount) VALUES (?, ?, ?)`,
		userId, time.Now(), float64(amount)).Exec(); err != nil {
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

func showWallet(userId string) {
	var amount float64

	if err := session.Query(`SELECT SUM(amount) from transactions WHERE userId=?`, userId).Consistency(gocql.One).Scan(&amount); err != nil {
		log.Fatal(err)
	}

	fmt.Println("$:", amount)
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

	results := session.Query(`SELECT * FROM bids`).Iter()
	m := &map[string]interface{}{}
	ret := []map[string]interface{}{}
	var maxBid float64
	var winnerId string
	var currentBid float64
	var currentBidderId string
	var vendorId string

	for results.MapScan(*m) {
		ret := append(ret, *m)
		m = &map[string]interface{}{}
		currentBid, _ = strconv.ParseFloat(fmt.Sprint(ret[0]["bid"]), 64)
		currentBidderId = fmt.Sprint(ret[0]["bidderid"])
		if (currentBid > maxBid) {
			maxBid = currentBid
			winnerId = currentBidderId
			vendorId = fmt.Sprint(ret[0]["auctionid"])
		}
	}

	// and the winner is
	fmt.Println(winnerId, maxBid)

	transaction(winnerId, -maxBid)
	transaction(vendorId, maxBid)

	showWallet(vendorId)
}