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

	reTimeSeconds := regexp.MustCompile(`([^\d]|^)1\d{9}\.\d+`)          // Seconds since the Epoch, includes subseconds
	reTimeNanoseconds := regexp.MustCompile(`([^\d]|^)1\d{18}([^\d]|$)`) // Seconds since the Epoch, includes subseconds

	for scanner.Scan() {
		s := reTimeSeconds.ReplaceAllStringFunc(scanner.Text(), secondsToTime)
		s = reTimeNanoseconds.ReplaceAllStringFunc(s, nanosecondsToTime)
		fmt.Println(s) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input: %s", err.Error())
	}
}

func secondsToTime(s string) string {
	prefix := ""
	if s[0] != uint8('1') {
		prefix = string(s[0])
		s = s[1:]
	}
	ses := strings.Split(s, ".")
	seconds, err := strconv.ParseInt(ses[0], 10, 64)
	check(err)
	nanoseconds, err := strconv.ParseInt((ses[1] + "000000000")[:9], 10, 64)
	check(err)
	return prefix + time.Unix(seconds, nanoseconds).String()
}

// static variable, declaring at global scope seems like the least evail
// https://stackoverflow.com/questions/30558071/static-local-variable-in-go
var reDigits = regexp.MustCompile(`\d`)

func nanosecondsToTime(s string) string {
	prefix := ""
	suffix := ""
	if s[0] != uint8('1') {
		prefix = string(s[0])
		s = s[1:]
	}
	a := reDigits.Find([]uint8(string(s[len(s)-1])))
	if a == nil {
		suffix = string(s[len(s)-1])
		s = s[:len(s)-1]
	}
	nanoseconds, err := strconv.ParseInt(s, 10, 64)
	check(err)
	return prefix + time.Unix(0, nanoseconds).String() + suffix
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
