package u2date

import (
	"io"
	"strings"
	"time"
	"fmt"
)

type DateScan struct {
	Now   time.Time
	regex string
}

func U2date(reader io.Reader) string {
	p := make([]byte, 9)
	_, err := reader.Read(p)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(p), "\000")
}

func (ds DateScan) Initialize(t time.Time) {
	ds.Now = t
	ds.regex = fmt.Sprintf("%d", ds.Now.Unix())
}