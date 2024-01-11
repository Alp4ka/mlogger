# MLogger
Logger with additional contact points and message templates based on slog for golang 1.21.1+

## Issues
1) Performance improvements. Current version is the latest revision of logger over slog. I have pretty much questions about its structure.
2) Render of templates at the gateway's level. Looks like a good idea until the moment you realize that all SMs just don't give a fuck at what the Markdown actually is and how it works.
3) Tests. At the moment my tests actually are a bunch of examples. 
4) Documentation. To be more concrete - the lack of it.

## Contact points implementations:
1) Matrix
2) Telegram

## Parse Modes
1) Markdown

# Usage 
## Install
```go get -u github.com/Alp4ka/mlogger```
## Example
See examples in **test_example.go** file.