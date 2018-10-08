package main

import (
	"os"
	"fmt"
	"github.com/domdom82/go-backpressure/server"
)

func main() {
	var sconfig *server.ServerConfig
	//var cconfig *client.ClientConfig
	if len(os.Args) > 1 {
		fmt.Printf("Loading sconfig from %q\n", os.Args[1])
		s, err := server.NewServerConfigFromFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		sconfig = s
	} else {
		fmt.Printf("ERROR: No configuration file provided.\nUsage: %s <-client|-server> <configfile>", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf("Starting the server.\n")
	sconfig.NewServer().Run()
}
