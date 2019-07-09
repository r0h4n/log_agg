# Log_agg

log aggregation service.

## Routes:

| Route | Description | Payload | Output |
| --- | --- | --- | --- |
| **Post** / | Post a log | json Log object | success message string |
| **Get** / | List all services | None | json array of Log objects |

### Query Parameters:
| Parameter | Description |
| --- | --- |
| **id** | Filter by id |
| **tag** | Filter by tag |
| **type** | Filter by type |
| **start** | Start time (unix epoch(nanoseconds)) at which to view logs older than (defaults to now) |
| **end** | End time (unix epoch(nanoseconds)) at which to view logs newer than (defaults to 0) |
| **limit** | Number of logs to read (defaults to 100) |
| **level** | Severity of logs to view (defaults to 'trace') |
`?id=my-app&tag=apache%5Berror%5D&type=deploy&start=0&limit=5`

## Data types:
### Log:
```json
{
  "id": "my-app",
  "tag": "build-1234",
  "type": "deploy",
  "priority": "4",
  "message": "$ mv r0h4n/.htaccess .htaccess\n[✓] SUCCESS"
}
```
| Field | Description |
| --- | --- |
| **time** | Timestamp of log (`time.Now()` on post) |
| **id** | Id or hostname of sender |
| **tag** | Tag for log |
| **type** | Log type (commonly 'app' or 'deploy'. default value configured via `log-type`) |
| **priority** | Severity of log (0(trace)-5(fatal)) |
| **message*** | Log data |
Note: * = required on submit


## Usage

publish log - success
```
$ curl -i http://localhost:6360 -d '{"id":"my-app","type":"deploy","message":"$ mv r0h4n/.htaccess .htaccess\n[✓] SUCCESS"}'
sucess!
HTTP/1.1 200 OK
```

get deploy logs
```
$ curl http://localhost:6360?kind=deploy 
[{"time":"2016-03-07T15:48:57.668893791-07:00","id":"my-app","tag":"","type":"deploy","priority":0,"message":"$ mv r0h4n/.htaccess .htaccess\n[✓] SUCCESS"}]
```

get app logs
```
$ curl http://localhost:6360 
[]
```
