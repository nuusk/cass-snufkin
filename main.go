package main

import (
	"bufio"
	"fmt"
	"strings"
	"os"
)

func main() {
	fmt.Println("Welcome in 5 Seconds Auctions!!!")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(" ls - lists auctions\n join [id] - joins auction\n q - exit\n")
	
	app_state := true

	for app_state {
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSuffix(cmd, "\n")
		command := strings.Split(cmd, " ")

		switch command[0] {
		case "ls":
			fmt.Println("ls -> lists auctions")
		case "join":
			fmt.Println("joined auction")
		case "q":
			app_state = false
		default:
			fmt.Println("Syntax error")
		}

	}

}