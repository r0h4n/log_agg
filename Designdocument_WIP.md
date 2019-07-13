Requirement: Simple Log aggregator service (only supports http api)

Modules:
1) Api, Http
2) Input, logs collector (currently only http, syslog support being added https://github.com/r0h4n/log_agg/issues/1)
3) Transform, transforms logs and channels them towards the Output/Archive module
4) Output, log file Archiver (boltdb, can publish to syslog or other products like datadog etc)

TODO:
1) Redact, implement  in the transform module in form of user defined config.redact_regex

