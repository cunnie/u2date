## u2date [![u2date](https://ci.nono.io/api/v1/pipelines/u2date/jobs/integration/badge)](https://ci.nono.io/teams/main/pipelines/u2date)

`u2date` is a command-line filter which converts [UNIX Epoch
time](https://en.wikipedia.org/wiki/Unix_time) (e.g. `1530473256.452262878`
(seconds) or `1530473256452262878` (nanoseconds)) to human-readable time (e.g.
`2018-07-01 12:27:36.452262878 -0700 PDT`).

### Linux Quick Start

```shell
curl -o u2date -L https://github.com/cunnie/u2date/releases/download/1.3.0/u2date-linux-amd64
chmod +x u2date
./u2date < /var/vcap/sys/log/atc/atc.stdout.log | less
```

sample input (notice the timestamps are in UNIX Epoch time):
```json
{"timestamp":"1530473256.452262878","source":"atc","message":"atc.db.failed-to-open-db-retrying","log_level":2,"data":{"error":"dial tcp 10.128.0.4:5432: connect: connection refused","session":"3"}}
{"timestamp":"1530473262.311251402","source":"atc","message":"atc.build-tracker.track.start","log_level":0,"data":{"session":"34.1"}}
{"timestamp":"1530473262.311525583","source":"atc","message":"atc.listening","log_level":1,"data":{"debug":"127.0.0.1:8079","http":"0.0.0.0:8080","https":"0.0.0.0:443"}}
```

sample output (notice the timestamps are human-readable):
```json
{"timestamp":"2018-07-01 12:27:36.452262878 -0700 PDT","source":"atc","message":"atc.db.failed-to-open-db-retrying","log_level":2,"data":{"error":"dial tcp 10.128.0.4:5432: connect: connection refused","session":"3"}}
{"timestamp":"2018-07-01 12:27:42.311251402 -0700 PDT","source":"atc","message":"atc.build-tracker.track.start","log_level":0,"data":{"session":"34.1"}}
{"timestamp":"2018-07-01 12:27:42.311525583 -0700 PDT","source":"atc","message":"atc.listening","log_level":1,"data":{"debug":"127.0.0.1:8079","http":"0.0.0.0:8080","https":"0.0.0.0:443"}}
```

### Technical Notes

On Linux, you can use environment variable `TZ` to set the output timezone.
For example, say the server's time is set to UTC (a common practice), but you'd
prefer to see the output in Toronto time ([complete list of `TZ`
values](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)), you'd
set the `TZ` variable as in the following example:

```bash
TZ=America/Toronto ./u2date < /var/vcap/sys/log/atc/atc.stdout.log | less
```

*[My personal favorite is `TZ=America/Los_Angeles`]*

You may want to wrap the converted time in double-quotes (") if, for example,
the input is valid JSON and the the timestamp is a number and not a string. This
preserves JSON-compatibility. Use `u2date -wrap=\"` to surround the converted
date with quotes.

The UNIX Epoch time pattern-match is unsophisticated:

- When searching for timestamps in seconds, it looks for a number between 1
  billion and 2 billion seconds (which is an
  overly-broad search spanning from Sunday, September 9, 2001 01:46:40 UTC to
  Wednesday, May 18, 2033 03:33:20 UTC)
  - which does not have commas (e.g. it would _not_ match "1,530,473,256.4")
  - which has a decimal point followed by at least one number (e.g. it would
    _not_ match "1530473256", but it would match "1530473256.4")
- Similarly for the nanosecond timestamps, it searches for a number between 1
  quadrillion and 2 quadrillion
  - unlike the seconds-based timestamp search, the nanosecond search does _not_ look for
    a decimal point

### Developer Notes

Developer notes can be found [here](docs/DEVELOPER.md).
