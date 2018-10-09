package main

import (
	"fmt"
	"os"

	"github.com/domdom82/go-backpressure/client"
	"github.com/domdom82/go-backpressure/server"
)

func main() {

	usage := func() {
		fmt.Printf("Usage: %s <-client|-server> <configfile>\n", os.Args[0])
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		if os.Args[1] == "-server" {
			fmt.Printf("Loading server config from %q\n", os.Args[2])
			s, err := server.NewServerConfigFromFile(os.Args[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Starting the server.\n")
			s.NewServer().Run()
		} else if os.Args[1] == "-client" {
			fmt.Printf("Loading client config from %q\n", os.Args[2])
			c, err := client.NewClientConfigFromFile(os.Args[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("Starting the client.\n")
			c.NewClient().Run()
		}
	} else {
		usage()
	}

}
