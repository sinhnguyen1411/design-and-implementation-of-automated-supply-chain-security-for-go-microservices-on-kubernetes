package main

import (
	"github.com/sinhnguyen1411/stock-trading-be/services/user-service/cmd"
	"os"
)

func main() {
	appCli := cmd.AppCommandLineInterface()
	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}

}
