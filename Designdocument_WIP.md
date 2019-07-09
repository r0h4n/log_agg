Requirement: Simple Log aggregator service

Modules:
1) Api, Http
2) Input, logs collector (currently only http)
3) Transform, transforms logs and channels them towards the Output/Archive module
4) Output, log file Archiver (boltdb, can publish to syslog or other products like datadog etc)

TODO:
1) Redact, implement  in the transform module in form of user defined config.redact_regex

