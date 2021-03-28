package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goku321/chart-fetcher/chart"
)

func usage() {
	fmt.Printf("Usage:\n./chart-fetcher <url> <count>\n")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	count, err := strconv.ParseInt(os.Args[2], 10, 32)
	if err != nil {
		fmt.Printf("invalid argument: count should be an integer\n")
		os.Exit(1)
	}
	f := chart.NewFetcher(os.Args[1], int(count))
	f.Init()
	if err = f.Start(); err != nil {
		fmt.Printf("cannot fetch chart: %s", err)
		os.Exit(1)
	}
	f.Chart.PrintJSON()
}
