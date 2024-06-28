package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aerochrome/piper/internal/interpreter"
	"github.com/aerochrome/piper/internal/parser"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func main() {
	forkPtr := flag.Bool("repl", false, "Starts a repl")
	flag.Parse()

	if !*forkPtr {
		handleFileRead()
	} else {
		handleRepl()
	}
}

func handleRepl() {
	// Create a channel to handle 'Ctrl+C' (SIGINT)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nExiting...")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("--- The Piper Language REPL ---\nType an expression and instantly see the result\nPress q to quit")

	interpreterStruct := interpreter.NewInterpreter()

	for {
		fmt.Print(">")
		input := scanner.Scan()
		if !input {
			return
		}

		if scanner.Text() == "q" {
			fmt.Println("Quitting...")
			break
		}

		// Parse
		res, err := parser.Parse("", scanner.Bytes())
		if err != nil {
			fmt.Println(err)
			break
		}

		// Transform
		for idx, v := range res.([]any) {
			res.([]any)[idx] = interpreter.TransformToReversePolishNotation(v.([]any))
		}

		// Evaluate
		for _, v := range res.([]any) {
			evalResult := interpreterStruct.Evaluate(v.([]any), nil)
			fmt.Println(evalResult)
		}

	}
}

func handleFileRead() {
	dat, err := os.ReadFile("./src.p")

	if err != nil {
		panic(err)
	}

	res, err := parser.Parse("", dat)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Transform
	for idx, v := range res.([]any) {
		res.([]any)[idx] = interpreter.TransformToReversePolishNotation(v.([]any))
	}

	result, err := json.MarshalIndent(res, "", "   ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("output:\n%+v\n", string(result))

	printDebug(res.([]any))

	interpreterStruct := interpreter.NewInterpreter()

	for _, v := range res.([]any) {
		evalResult := interpreterStruct.Evaluate(v.([]any), nil)

		fmt.Println(evalResult)
	}
}

func printDebug(res []any) {
	fmt.Println("Debug:")

	for idx, v := range res {
		if idx > 0 {
			fmt.Print(",\n")
		}
		printAnySliceRecursive(v.([]any))
	}

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
