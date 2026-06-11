package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"WinUpdBrief/internal/app"
)

func main() {
	detail := flag.Bool("detail", false, "show detailed update information")
	flag.Parse()

	if err := app.Run(context.Background(), app.Options{Detail: *detail}, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
