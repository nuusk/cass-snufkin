# cass-snufkin

<img src="https://raw.githubusercontent.com/pietersweter/github-images/master/pietersweter/cass-snufkin/index.png" width="300">

## Who is snufkin?

Snufkin is a philosophical vagabond who wanders the world fishing and playing the harmonica. He carries everything he needs in his pouch, as he believes that having too much stuff makes life overly complicated.

## How to play?

You play as a snufkin. You can explore the world and collect the artifacts hidden there. You can have only one item in your pouch at a time. When you collect the artifact, you can either sell it or exhibit a charity with it. Other snufkins (from alternative worlds) can participate in the auction to get your item, or they can receive this item from a charity (if you chose them as a recipient)

### Actions

On your turn, you can take one of the following actions:

1. *Pouch* - show your invenory and your account balance
2. *Place an auction* - sell your pouch item. Other snufkins will bid for it. At the end of an auction, the auction winner gets the item and pays the amount of money he declared. You receive the money, but lose the pouch item.
3. *Exhibit a charity* - put your pouch item on a chaity event and choose a snufkin. Other snufkins can donate their money for the cause of your choosing. All the given money will be transfered to the recipient. The winner gets your pouch item.
4. *Help* - show help info
5. *Quit* - exit the game
6. *Show snufkins* - display info about other snufkins and their account balance

## How it works

When a snufkin makes a bid, the bid value is stored in *bids* table. This table structure is the following:
***(auctionId text, itemName text, time timestamp, bid double, bidderId text, PRIMARY KEY((auctionId, bidderId, itemName), time))***.
As one can see, bids contain a timestamp, which helps to distinguish bids made by the same player on the same auction. The auction ends in a specific time. Those are the default values:

```
var NUM_SNUFKINS = 4
var BID_INTERVAL_MS time.Duration = 100
var REFRESH_TIMES_MS time.Duration = 500
var AUCTION_TIMER_S time.Duration = 5
```

***NUM_SNUFKINS*** corresponds to the number of snufkins in a simulation.
***BID_INTERVAL_MS*** is a time (in milliseconds) interval between each bid of every snufkin in a simulation.
***REFRESH_TIMES_MS*** is a time (in milliseconds) interval between debug statements. It ***does not*** affect the simulation.
***AUCTION_TIMER_S*** says how long each simulation will last.

When the auction ends, a snufkin with the biggest is declared a winner. He transfers the money to the vendor snufkin - that operation is stored in *transactions* table. This table structure is the following: ***(userId text, time timestamp, amount double, PRIMARY KEY(userId, time))***.
This structure works similarly to a distributed ledger which people can save data to.
The main difference between an auction and an exhibition is that auction only applies two transactions: one for the winner and one for the vendor. The exhibition naturally consists of many seperate transactions. This is a main reason for Cassandra in the project. It encourages many writes and a minimum number of reads.

## Prerequisites

- Golang
- Cassandra

## Getting started

Create keyspace for the project
```
create keyspace snufkin with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 2 };
```

Create tables using one of the following
```
CREATE TABLE bids (auctionId text, itemName text, time timestamp, bid double, bidderId text, PRIMARY KEY((auctionId, bidderId, itemName), time));
CREATE TABLE transactions (userId text, time timestamp, amount double, PRIMARY KEY(userId, time));
CREATE TABLE pouches (userId text, itemName text, itemStartingPrice double, PRIMARY KEY(userId));
```

```
go run helpers/init.go
```

Configure your *Cassandra cluster* and set up the connection in `bid.go` by changing *ip address* in ```cluster = gocql.NewCluster("127.0.0.1")```

After all this set up, run the program with
```
go run bid.go
```

## Authors

* **poe** - [pietersweter](https://github.com/pietersweter)
* **Pogoś** - [Siogop](https://github.com/siogop)
* **kamciokodzi** - [kamciokodzi](https://github.com/kamciokodzi)

## License

This project is licensed under the MIT License.
