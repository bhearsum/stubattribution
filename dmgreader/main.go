package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/willdurand/go-dmg-reader/dmg"
)

func main() {
	file, err := dmg.OpenFile("/home/bhearsum/tmp/2023-09-11/target.dmg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := file.Parse()
	if err != nil {
		panic(err)
	}
	// why everything nil?!
	log.Println(data.Koly)

	b, err := json.MarshalIndent(data.Koly, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
