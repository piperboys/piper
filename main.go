package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aerochrome/piper/internal/parser"
)

func main() {
	dat, err := os.ReadFile("./src.p")

	if err != nil {
		panic(err)
	}

	res, err := parser.Parse("", dat)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := json.MarshalIndent(res, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("output:\n%v\n", string(result))
}
