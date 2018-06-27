package u2date

import (
	"io"
	"strings"
)

func U2date(reader io.Reader) string {
	p := make([]byte, 9)
	_, err := reader.Read(p)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(p), "\000")
}
