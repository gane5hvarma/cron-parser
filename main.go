package main

import (
	"cron-parser/cron"
	"fmt"
	"os"
)

func main() {
	noOfArgs := len(os.Args)
	if noOfArgs != 2 {
		fmt.Println("Usage: cron '<cronexpression>'")
	}
	cp, err := cron.ParseCronExp(os.Args[1])
	if err != nil {
		fmt.Println("failed to parse cron", err)
	}
	cp.Print()
}
