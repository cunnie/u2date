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
	"unicode/utf8"
)

func main() {
	wrap := ""
	flag.StringVar(&wrap, "wrap", "", `string to wrap converted dates; for JSON compatibility, set to double-quote: -wrap=\"`)
	flag.Parse()
	input := bufio.NewScanner(os.Stdin)

	reTimeSeconds := regexp.MustCompile(`([^\d]|^)1\d{9}\.\d+`)              // seconds since the Epoch, includes subseconds
	reTimeNanoseconds := regexp.MustCompile(`([^\d]|^)1\d{18}([^\d]|$)`)     // nanoseconds since the Epoch, no subnanoseconds
	reTimeSecondsNoDecimal := regexp.MustCompile(`([^\d]|^)1\d{9}([^\d]|$)`) // seconds since the Epoch, no subseconds

	for input.Scan() {
		s := reTimeSeconds.ReplaceAllStringFunc(input.Text(), func(s string) string { return secondsToTime(s, wrap) })
		s = reTimeNanoseconds.ReplaceAllStringFunc(s, func(s string) string { return nanosecondsToTime(s, wrap) })
		s = reTimeSecondsNoDecimal.ReplaceAllStringFunc(s, func(s string) string { return secondsNoDecimalToTime(s, wrap) })
		fmt.Println(s) // Println will add back the final '\n'
	}
	check(input.Err())
}

func secondsToTime(rawSecondsString string, wrapper string) string {
	// The regex might pick up an optional non-decimal leading character, "prefix"
	prefix := ""
	eldritchRune, width := utf8.DecodeRuneInString(rawSecondsString)
	if eldritchRune != '1' {
		prefix = string(eldritchRune)
		rawSecondsString = rawSecondsString[width:]
	}
	secondsAndNanoseconds := strings.Split(rawSecondsString, ".")
	seconds, err := strconv.ParseInt(secondsAndNanoseconds[0], 10, 64)
	check(err)
	nanoseconds, err := strconv.ParseInt((secondsAndNanoseconds[1] + "000000000")[:9], 10, 64)
	check(err)
	return prefix + wrapper + time.Unix(seconds, nanoseconds).String() + wrapper
}

// static variable, declaring at global scope seems like the least evil
// https://stackoverflow.com/questions/30558071/static-local-variable-in-go
var reDigits = regexp.MustCompile(`\d`)

func nanosecondsToTime(rawNanosecondsString string, wrapper string) string {
	// The regex might pick up optional non-decimal leading & trailing characters, "prefix" and "suffix"
	prefix := ""
	suffix := ""
	eldritchRune, width := utf8.DecodeRuneInString(rawNanosecondsString)
	if eldritchRune != '1' {
		prefix = string(eldritchRune)
		rawNanosecondsString = rawNanosecondsString[width:]
	}
	eldritchRune, width = utf8.DecodeLastRuneInString(rawNanosecondsString)
	if !reDigits.MatchString(string(eldritchRune)) {
		suffix = string(eldritchRune)
		rawNanosecondsString = rawNanosecondsString[:len(rawNanosecondsString)-width]
	}
	nanoseconds, err := strconv.ParseInt(rawNanosecondsString, 10, 64)
	check(err)
	return prefix + wrapper + time.Unix(0, nanoseconds).String() + wrapper + suffix
}

func secondsNoDecimalToTime(rawSecondsString string, wrapper string) string {
	// The regex might pick up optional non-decimal leading & trailing characters, "prefix" and "suffix"
	prefix := ""
	suffix := ""
	eldritchRune, width := utf8.DecodeRuneInString(rawSecondsString)
	if eldritchRune != '1' {
		prefix = string(eldritchRune)
		rawSecondsString = rawSecondsString[width:]
	}
	eldritchRune, width = utf8.DecodeLastRuneInString(rawSecondsString)
	if !reDigits.MatchString(string(eldritchRune)) {
		suffix = string(eldritchRune)
		rawSecondsString = rawSecondsString[:len(rawSecondsString)-width]
	}
	seconds, err := strconv.ParseInt(rawSecondsString, 10, 64)
	check(err)
	return prefix + wrapper + time.Unix(seconds, 0).String() + wrapper + suffix
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
