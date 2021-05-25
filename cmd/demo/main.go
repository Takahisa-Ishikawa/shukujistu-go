package main

import (
	"fmt"

	shukujitsu "github.com/Takahisa-Ishikawa/shukujistu-go"
)

func main() {
	entries, err := shukujitsu.AllEntries()
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		fmt.Printf("%s = %s\n", e.YMD, e.Name)
	}
}
