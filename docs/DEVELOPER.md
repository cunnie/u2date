## Developer Notes

- Pull requests (PRs) are welcome
- PRs should include test coverage
- PRs should be squashed in as few commits as possible
- PRs should have a good [commit message](https://chris.beams.io/posts/git-commit/)

### Running tests

`u2date` uses the Golang BDD testing framework, [Ginkgo](https://github.com/onsi/ginkgo):

```
go get github.com/cunnie/u2date
cd $GOPATH/src/github.com/cunnie/u2date
ginkgo -r .
```

### TODO

We'd like to be able to parse nanosecond output (i.e. the number of nanoseconds
elapsed since January 1, 1970 UTC, e.g. the output from Golang's
[`time.UnixNano()](https://golang.org/pkg/time/#Time.UnixNano)
"1530473256452262878")

We'd like to remove the decimal-point requirement (so that "1530473256" would
match properly).

We'd like to narrow the matching window (a 32-year span is rather broad).