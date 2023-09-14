package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/willdurand/go-dmg-reader/dmg"
)

func main() {
	file, err := dmg.OpenFile("/home/bhearsum/repos/stubattribution/dmgreader/tests/fixtures/attributable.dmg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := file.Parse()
	if err != nil {
		panic(err)
	}

	rd, _ := data.Resources.GetResourceDataByName("blkx")

	blkx, err := dmg.ParseBlkxData(rd[3].Data)
	if err != nil {
		log.Printf("err: %s", err)
	}
	b, err := json.MarshalIndent(blkx, "", "  ")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	fmt.Println(string(b))
}
