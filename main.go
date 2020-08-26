package main

import (
    "fmt"
    "github.com/nielsing/yar/internal/args"
)

func main() {
    parsedArgs := args.ParseArgs()
    fmt.Println("Args:", parsedArgs)
}
