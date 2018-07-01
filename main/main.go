package main

import (
	"fmt"
	"github.com/cunnie/u2date/u2date"
	"os"
	"bufio"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	os.Exit(0)
	fmt.Println(u2date.U2date(os.Stdin))
}
