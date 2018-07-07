package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	reTime := regexp.MustCompile(`1[0-9]{9}\.[0-9]+`) // Seconds since the Epoch, includes subseconds

	for scanner.Scan() {
		s := reTime.ReplaceAllStringFunc(scanner.Text(), toTime)
		fmt.Println(s) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input: %s", err.Error())
	}
}

func toTime(s string) string {
	ses := strings.Split(s, ".")
	seconds, err := strconv.ParseInt(ses[0], 10, 64)
	check(err)
	nanoseconds, err := strconv.ParseInt((ses[1] + "000000000")[:9], 10, 64)
	check(err)
	return time.Unix(seconds, nanoseconds).String()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
