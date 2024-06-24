package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

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
	fmt.Printf("output:\n%+v\n", string(result))

	printDebug(res.([]any))
}

func printDebug(res []any) {
	fmt.Println("Debug:")

	printAnySliceRecursive(res)

	fmt.Printf("\n")
}

func printAnySliceRecursive(slice []any) {
	for idx, item := range slice {
		if idx > 0 {
			fmt.Print(", ")
		}

		// Is any slice?
		itemSlice, ok := item.([]any)
		if ok {
			fmt.Print("\n[")
			printAnySliceRecursive(itemSlice)
			fmt.Print("]")
			continue
		}

		// is struct?
		if reflect.ValueOf(item).Kind().String() == "struct" {
			fmt.Printf("%s%+v", reflect.TypeOf(item), item)
			continue
		}

		// everything else
		fmt.Printf("%v", item)
	}
}
