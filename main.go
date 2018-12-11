package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mnordsletten/lotto_dashboard/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Must provide port to run on")
		os.Exit(1)
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("error converting port number to int: ", err)
		os.Exit(1)
	}
	server.Serve(port)
}
