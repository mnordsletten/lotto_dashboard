package main

import (
	"fmt"
	"martin/lotto_dashboard/server"
	"os"
	"strconv"
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
