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

var NUM_SNUFKINS = 4
var snufkins []string
var items []string

var cluster *gocql.ClusterConfig
var session *gocql.Session
var simulationActive bool

func initSnufkins() {
	for i := 0; i < NUM_SNUFKINS; i++ {
		newId := uuid.Must(uuid.NewV4()).String()
		snufkins = append(snufkins, newId)
	}
}

func initItems() {
	items = []string{
		"Legendary Sword Mionsz++",
		"Ancient Belt for Loosers +2",
		"Common Fire Bow",
		"Rare Chest Armor",
		"Epic Axe of Piter Bird",
		"Distributed Ice Spell",
	}
}

func printSnufkins() {
	var amount float64

	for i := 0; i < NUM_SNUFKINS; i++ {
		if err := session.Query(`SELECT SUM(amount) from transactions WHERE userId=?`, snufkins[i]).Consistency(gocql.One).Scan(&amount); err != nil {
			log.Fatal(err)
		}

		fmt.Println(i, snufkins[i], amount)
	}
}

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
		log.Fatal(err)
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
	fmt.Println("Simulation state changed to", simulationActive)
}

func simulateUserInAuction(snufkinId string, auctionId string, itemName string, bidFinished chan bool) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	maxRandomBidAmount := 10000
	var simulationActive bool = isActive()

	for simulationActive {
			simulationActive = isActive()
			bidAmount := r1.Intn(maxRandomBidAmount)

			bid(auctionId, itemName, snufkinId, bidAmount)
			fmt.Println(snufkinId, "biddig", bidAmount)
			time.Sleep(200 * time.Millisecond)
	}
	bidFinished <- true
}

func simulateUserInExhibition(snufkinId string, auctionId string, itemName string, bidFinished chan bool) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	maxRandomBidAmount := 10000
	var simulationActive bool = isActive()

	for simulationActive {
			simulationActive = isActive()
			bidAmount := r1.Intn(maxRandomBidAmount)

			bid(auctionId, itemName, snufkinId, bidAmount)
			fmt.Println(snufkinId, "biddig", bidAmount)
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
	deleteFromPouch(snufkins[0])
}

func endExhibition(itemName string, chosenSnufkin string) {
	fmt.Println("Exhibition ended!")

	results := session.Query(`SELECT * FROM bids WHERE itemName=? ALLOW FILTERING`, itemName).Iter()
	m := &map[string]interface{}{}
	ret := []map[string]interface{}{}
	var currentBid float64
	var currentBidderId string

	for results.MapScan(*m) {
		ret := append(ret, *m)
		m = &map[string]interface{}{}
		currentBid, _ = strconv.ParseFloat(fmt.Sprint(ret[0]["bid"]), 64)
		currentBidderId = fmt.Sprint(ret[0]["bidderid"])
		transaction(currentBidderId, -currentBid)
		transaction(chosenSnufkin, currentBid)
	}

	deleteFromPouch(snufkins[0])
}

func findItem(snufkinId string) {
	rand.Seed(time.Now().UnixNano())
    foundItem := items[rand.Intn(len(items))]
	fmt.Println("Congrats! You found", foundItem)

	if err := session.Query(`INSERT INTO pouches (userId, itemName, itemStartingPrice) VALUES (?, ?, ?)`,
		snufkinId, foundItem, 0.0).Exec(); err != nil {
		log.Fatal(err)
	}
}

func deleteFromPouch(snufkinId string) {
	if err := session.Query(`DELETE FROM pouches WHERE id=? IF EXISTS;`,
		snufkinId).Exec(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initSnufkins()
	initItems()
	refreshFinished := make(chan bool)
	bidFinished := make(chan bool)
	mainUserId := snufkins[0]

	// connect to the cluster
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "snufkin"
	cluster.Consistency = gocql.One
	session, _ = cluster.CreateSession()
	defer session.Close()

	actions := []action{
		{Description: "Pouch", Id: 0},
		{Description: "Place an auction", Id: 1},
		{Description: "Exhibit a charity", Id: 2},
		{Description: "Explore", Id: 3},
		{Description: "Help", Id: 4},
		{Description: "Quit", Id: 5},
		{Description: "Show snufkins", Id: 6},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F300 {{ .Description | cyan }} ({{ .Id | red }})",
		Inactive: "  {{ .Description | white }} ({{ .Id | red }})",
		Selected: "\U0001F300 {{ .Description | red | cyan }}",
	}

	clearScreen()
	
	for {

		prompt := promptui.Select{
			Label: "What do you want to do?",
			Items: actions,
			Templates: templates,
			Size:      7,
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
			clearScreen()
			item := getPouch(mainUserId)
			if (item != "<empty>") {
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

					for i := 1; i < NUM_SNUFKINS; i++ {
						go simulateUserInAuction(snufkins[i], mainUserId, item, bidFinished)
					}

					<- refreshFinished
					<- bidFinished
					
					endAuction(item)
				} else {
					break;
				}
			} else {
				fmt.Println("Your pouch is empty!")
				break;
			}
		case 2:
			clearScreen()
			item := getPouch(mainUserId)
			if (item != "<empty>") {
				fmt.Println("Who is this exhibition for?")
				
				printSnufkins()

				prompt := promptui.Prompt{
					Label:  "Choose index (or c to cancel):",
				}

				result, err := prompt.Run()

				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}

				fmt.Println(result)

				if (strings.ToLower(result) == "c") {
					break;
				} else {
					chosenSnufkinIndex, _ := strconv.Atoi(result)
					chosenSnufkin := snufkins[chosenSnufkinIndex]
					// simulate an auction
					fmt.Println("Simulating an exhibition...")

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

					for i := 1; i < NUM_SNUFKINS; i++ {
						go simulateUserInExhibition(snufkins[i], mainUserId, item, bidFinished)
					}

					<- refreshFinished
					<- bidFinished
					
					endExhibition(item, chosenSnufkin)
				}
			} else {
				fmt.Println("Your pouch is empty!")
				break;
			}
		case 3:
			clearScreen()
			findItem(snufkins[0])
		case 4:
			clearScreen()
			fmt.Println("Ask the developers for help")
		case 5:
			os.Exit(3)
		case 6:
			printSnufkins()
		}
	}
}