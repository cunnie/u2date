package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	wrap := ""
	flag.StringVar(&wrap, "wrap", "", `string to wrap converted dates; for JSON compatibility, set to double-quote: -wrap=\"`)
	flag.Parse()
	input := bufio.NewScanner(os.Stdin)

	reTimeSeconds := regexp.MustCompile(`([^\d]|^)1\d{9}\.\d+`)          // seconds since the Epoch, includes subseconds
	reTimeNanoseconds := regexp.MustCompile(`([^\d]|^)1\d{18}([^\d]|$)`) // nanoseconds since the Epoch, no subseconds

	for input.Scan() {
		s := reTimeSeconds.ReplaceAllStringFunc(input.Text(), func(s string) string { return secondsToTime(s, wrap) })
		s = reTimeNanoseconds.ReplaceAllStringFunc(s, func(s string) string { return nanosecondsToTime(s, wrap) })
		fmt.Println(s) // Println will add back the final '\n'
	}
	check(input.Err())
}

func secondsToTime(rawSecondsString string, wrapper string) string {
	// The regex might pick up an optional non-decimal leading character, "prefix"
	prefix := ""
	if rawSecondsString[0] != uint8('1') {
		prefix += string(rawSecondsString[0])
		rawSecondsString = rawSecondsString[1:]
	}
	secondsAndNanoseconds := strings.Split(rawSecondsString, ".")
	seconds, err := strconv.ParseInt(secondsAndNanoseconds[0], 10, 64)
	check(err)
	nanoseconds, err := strconv.ParseInt((secondsAndNanoseconds[1] + "000000000")[:9], 10, 64)
	check(err)
	return wrapper + prefix + time.Unix(seconds, nanoseconds).String() + wrapper
}

// static variable, declaring at global scope seems like the least evil
// https://stackoverflow.com/questions/30558071/static-local-variable-in-go
var reDigits = regexp.MustCompile(`\d`)

func nanosecondsToTime(rawNanosecondsString string, wrapper string) string {
	// The regex might pick up optional non-decimal leading & trailing characters, "prefix" and "suffix"
	prefix := ""
	suffix := ""
	if rawNanosecondsString[0] != uint8('1') {
		prefix = string(rawNanosecondsString[0])
		rawNanosecondsString = rawNanosecondsString[1:]
	}
	lastCharIfDecimal := reDigits.Find([]uint8(string(rawNanosecondsString[len(rawNanosecondsString)-1])))
	if lastCharIfDecimal == nil {
		suffix = string(rawNanosecondsString[len(rawNanosecondsString)-1])
		rawNanosecondsString = rawNanosecondsString[:len(rawNanosecondsString)-1]
	}
	nanoseconds, err := strconv.ParseInt(rawNanosecondsString, 10, 64)
	check(err)
	return wrapper + prefix + time.Unix(0, nanoseconds).String() + suffix + wrapper
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
