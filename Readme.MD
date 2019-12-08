# Work in Progress
# Reverse proxy with logging
A reverse proxy with logging capabilities made for development purposses

# How to use it
## Command Line
Redirect traffic from port 8080 to port 6000
```shell
./go-redirect-proxy -source=8080 -target=6000
```
## Config File
```shell
./go-redirect-proxy -file=config.json
```

## Config file format

```json
{
"rules":[{
       "source":"8080",
       "target":"7090"
      }]
}
```

# View Last Requests Logs
View request which path match the regex pattern from query and sort "desc" or "asc" by timestamp.
```shell
http://localhost:8001/?path=<regex>&sort=desc
```

