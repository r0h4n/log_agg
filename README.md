# Log_agg
[![GoDoc](https://godoc.org/github.com/r0h4n/log_agg?status.svg)](https://godoc.org/github.com/r0h4n/log_agg)

log aggregation service.

## Quickstart

```sh
# start server (may require commented flags)
log_agg # -d /tmp/log_agg.db -u 0.0.0.0:6361

# add a log via http
curl http://0.0.0.0:6360/logs \
     -d '{"id":"log-test", "type":"log", "message":"my first log"}'

# view log via http
curl "http://0.0.0.0:6360/logs?type=log"

```

## Usage
```
  log_agg [flags]
```

Flags:
```
  -c, --config-file string    config file location for server
  -d, --db-address string     Log storage address (default "boltdb:///var/db/log_agg.bolt")
  -a, --listen-http string    API listen address (same endpoint for http log collection) (default "0.0.0.0:6360")
  -k, --log-keep string       Age or number of logs to keep per type '{"app":"2w", "deploy": 10}'' (int or X(m)in, (h)our,  (d)ay, (w)eek, (y)ear) (default "{\"app\":\"2w\"}")
  -l, --log-level string      Level at which to log (default "info")
  -L, --log-type string       Default type to apply to incoming logs (commonly used: app|deploy) (default "app")
  -v, --version               Print version info and exit
```

Config File: (takes precedence over cli flags)
```json
// log_agg.json
{
  "listen-http": "0.0.0.0:6360",
  "db-address": "boltdb:///var/db/log_agg.bolt",
  "log-keep": "{\"app\":\"2w\"}",
  "log-type": "app",
  "log-level": "info",
}
```

#### Adding|Viewing Logs
See http examples [here](./api/README.md)  




Design:

Requirement: Simple Log aggregator service (only supports http api)

Modules:

    Api, Http
    Input, logs collector (currently only http, syslog support being added https://github.com/r0h4n/log_agg/issues/1)
    Transform, transforms logs and channels them towards the Output/Archive module
    Output, log file Archiver (boltdb, can publish to syslog or other products like datadog etc)

TODO:

    Redact, implement in the transform module in form of user defined config.redact_regex
