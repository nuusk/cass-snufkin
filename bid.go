package main

import (
	"log"
	"fmt"
	"time"
	"math/rand"
	"strconv"
	"os"
	"strings"

	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
	"github.com/manifoldco/promptui"
)

var cluster *gocql.ClusterConfig
var session *gocql.Session
var simulationActive bool

func bid(auctionId string, itemName string, bidderId string, bidAmount int) {
	if err := session.Query(`INSERT INTO bids (auctionId, itemName, time, bid, bidderId) VALUES (?, ?, ?, ?, ?)`,
		auctionId, itemName, time.Now(), float64(bidAmount), bidderId).Exec(); err != nil {
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

func getBalance(userId string) float64 {
	var amount float64

	if err := session.Query(`SELECT SUM(amount) from transactions WHERE userId=?`, userId).Consistency(gocql.One).Scan(&amount); err != nil {
		fmt.Println("$:", amount)
	}

	return amount
}

func getPouch(userId string) string {
	var item string
	
	iter := session.Query(`SELECT itemName from pouches WHERE userId=? LIMIT 1`, userId).Consistency(gocql.Quorum).Iter()
	
	for iter.Scan(&item) {
		if item != "" {return item}
	}

	return "<empty>"
}

func clearScreen() {
	print("\033[H\033[2J")
}

type action struct {
	Description string
	Id int
}

func isActive() bool {
	return simulationActive
}

func setActive(toggle bool) {
	simulationActive = toggle
	fmt.Println("setting to", simulationActive)
}

func simulateUserInAuction(auctionId string, itemName string, bidFinished chan bool) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	maxRandomBidAmount := 10000
	id := uuid.Must(uuid.NewV4()).String()
	var simulationActive bool = isActive()

	for simulationActive {
			simulationActive = isActive()
			bidAmount := r1.Intn(maxRandomBidAmount)

			bid(auctionId, itemName, id, bidAmount)
			fmt.Println(id, "biddig", bidAmount)
			time.Sleep(200 * time.Millisecond)
	}
	bidFinished <- true
}

func endAuction(itemName string) {
	fmt.Println("Auction ended!")

	results := session.Query(`SELECT * FROM bids WHERE itemName=? ALLOW FILTERING`, itemName).Iter()
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
	fmt.Println("The winner is:", winnerId, maxBid)

	transaction(winnerId, -maxBid)
	transaction(vendorId, maxBid)
}

func main() {
	refreshFinished := make(chan bool)
	bidFinished := make(chan bool)
	// s1 := rand.NewSource(time.Now().UnixNano())
	// r1 := rand.New(s1)
	// maxRandomBidAmount := 10000
	mainUserId := uuid.Must(uuid.NewV4()).String()

	// connect to the cluster
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ = cluster.CreateSession()
	defer session.Close()

	actions := []action{
		{Description: "Pouch", Id: 0},
		{Description: "Place an auction", Id: 1},
		{Description: "Exhibit a charity", Id: 2},
		{Description: "Explore", Id: 3},
		{Description: "Help", Id: 4},
		{Description: "Quit", Id: 5},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F300 {{ .Description | cyan }} ({{ .Id | red }})",
		Inactive: "  {{ .Description | white }} ({{ .Id | red }})",
		Selected: "\U0001F300 {{ .Description | red | cyan }}",
	}

	for {
		prompt := promptui.Select{
			Label: "What do you want to do?",
			Items: actions,
			Templates: templates,
			Size:      6,
		}
	
		actionId, _, err := prompt.Run()
	
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
	
		clearScreen()

		switch actionId {
		case 0:
			item := getPouch(mainUserId)
			fmt.Println("Your pouch:", item)
			balance := getBalance(mainUserId)
			fmt.Println("$:", balance)
		case 1:
			item := getPouch(mainUserId)
			if (item == "<empty>") {
				fmt.Println("Do you want to sell:", item, "?")
				
				prompt := promptui.Prompt{
					Label:  "(Y/n)",
				}

				result, err := prompt.Run()

				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}

				if (result == "" || strings.ToLower(result) == "y") {
					// simulate an auction
					fmt.Println("Simulating an auction...")


					go func() {
						setActive(true)
						time.Sleep(400 * time.Millisecond)
						setActive(false)
					}()

					go func(refreshFinished chan bool) {
						for simulationActive {
								refresh()
								time.Sleep(1 * time.Second)
						}
						refreshFinished <- true
					}(refreshFinished)

					go simulateUserInAuction(mainUserId, item, bidFinished)
					go simulateUserInAuction(mainUserId, item, bidFinished)
					go simulateUserInAuction(mainUserId, item, bidFinished)
					go simulateUserInAuction(mainUserId, item, bidFinished)
					go simulateUserInAuction(mainUserId, item, bidFinished)

					<- refreshFinished
					<- bidFinished
					
					endAuction(item)

					// go func(bidFinished chan bool) {
					// 	id := uuid.Must(uuid.NewV4()).String()
						
					// 	for simulationActive {
					// 			bidAmount := r1.Intn(maxRandomBidAmount)
				
					// 			bid(id, bidAmount)
					// 			fmt.Println(id, "bidding", bidAmount)
					// 			time.Sleep(200 * time.Millisecond)
					// 	}
					// 	bidFinished <- true
					// }(bidFinished)
				
					// go func(bidFinished chan bool) {
					// 	id := uuid.Must(uuid.NewV4()).String()
						
					// 	for simulationActive {
					// 			bidAmount := r1.Intn(maxRandomBidAmount)
				
					// 			bid(id, bidAmount)
					// 			fmt.Println(id, "bidding", bidAmount)
					// 			time.Sleep(200 * time.Millisecond)
					// 	}
					// 	bidFinished <- true
					// }(bidFinished)
				
					// go func(bidFinished chan bool) {
					// 	id := uuid.Must(uuid.NewV4()).String()
						
					// 	for simulationActive {
					// 			bidAmount := r1.Intn(maxRandomBidAmount)
				
					// 			bid(id, bidAmount)
					// 			fmt.Println(id, "bidding", bidAmount)
					// 			time.Sleep(200 * time.Millisecond)
					// 	}
					// 	bidFinished <- true
					// }(bidFinished)

				} else {
					break;
				}
			} else {
				fmt.Println("Your pouch is empty!")
				break;
			}
		case 5:
			os.Exit(3)
		}
	}
}