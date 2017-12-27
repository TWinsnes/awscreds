package main

import (
	"github.com/TWinsnes/awscreds/console"
)

func main() {

	err := console.OpenConsole("temptoken", "vibrato-tom")

	if err != nil {
		panic(err)
	}
}
