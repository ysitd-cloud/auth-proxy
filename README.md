# OAuth Proxy

Proxy with OAuth 2 authorize

## Requirement

- Golang 1.10.2
- [dep](https://github.com/golang/dep)

## Config
| Name | Usage |
|------|-------|
| PORT | Port number for serving (default: 8080) |
| DB_URL | Data Source Url for connecting Postgres |
| VERBOSE | Enable Verbose loggging (default: false) |
| LOG_FILE | Filename for log output (optional) |
| SESSION_SECRET | Secret for session |