# cass-snufkin

<img src="https://raw.githubusercontent.com/pietersweter/github-images/master/pietersweter/cass-snufkin/index.png" width="440">

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

## Prerequisites

- Golang
- Cassandra

## Getting started

Create keyspace for the project
```
create keyspace snufkin with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 2 };
```

Create tables using one of the following
in cqlsh```
    CREATE TABLE bids (auctionId text, itemName text, time timestamp, bid double, bidderId text, PRIMARY KEY((auctionId, bidderId, itemName), time));
    CREATE TABLE transactions (userId text, time timestamp, amount double, PRIMARY KEY(userId, time));
    CREATE TABLE pouches (userId text, itemName text, itemStartingPrice double, PRIMARY KEY(userId));
```

or```
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
