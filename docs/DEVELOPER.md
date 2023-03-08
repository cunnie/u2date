## Developer Notes

- Pull requests (PRs) are welcome
- PRs should include test coverage
- PRs should be squashed in as few commits as possible
- PRs should have a good [commit message](https://chris.beams.io/posts/git-commit/)

### Running tests

`u2date` uses the Golang BDD testing framework, [Ginkgo](https://github.com/onsi/ginkgo):

```sh
cd u2date/u2date
go install github.com/onsi/ginkgo/v2/ginkgo # installs the v2 ginkgo CLI
go get github.com/onsi/gomega/... # fetches the matcher library
ginkgo -r .
```

### TODO

We'd like to be able to parse nanosecond output (i.e. the number of nanoseconds
elapsed since January 1, 1970 UTC, e.g. the output from Golang's
[`time.UnixNano()`](https://golang.org/pkg/time/#Time.UnixNano)
"1530473256452262878")

We'd like to narrow the matching window (a 32-year span is rather broad).
