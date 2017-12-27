package main

import (
	"flag"
	"os"

	"github.com/TWinsnes/awscreds/console"
)

func main() {

	consoleCommand := flag.NewFlagSet("console", flag.ExitOnError)

	consoleProfileFlag := consoleCommand.String("profile", "Default", "AWS profile to log into console with")
	consoleNameFlag := consoleCommand.String("tokenname", "temptoken", "Name of federated token")

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "console":
		consoleCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := console.OpenConsole(*consoleNameFlag, *consoleProfileFlag)

	if err != nil {
		panic(err)
	}
}
