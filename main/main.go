package main

import (
	"fmt"
	"github.com/cunnie/u2date/u2date"
	"os"
	"bufio"
	"io"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	buffer := make([]byte, 64*1024) // 64 kiB, admittedly arbitrary

	for {
		n, err := stdin.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Errorf("couldn't read stdin: %s", err.Error()))
		}
		stdout := bufio.NewWriter(os.Stdout)
		o, err := stdout.Write(buffer[:n])
		if err != nil {
			panic(fmt.Errorf("couldn't write stdout: %s", err.Error()))
		}
		if (n != o) {
			panic(fmt.Errorf("I expected to write %d, but instead wrote %d", n, o))
		}
	}
	os.Exit(0)
	fmt.Println(u2date.U2date(os.Stdin))
}
