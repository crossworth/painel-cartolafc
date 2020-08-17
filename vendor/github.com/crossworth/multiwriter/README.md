## Synced file aware MultiWriter

Example usage

```go
mw := multiwriter.MultiWriter{
    IO1: os.Stdout,
    IO2: myFile,
}

log.Logger = log.Output(mw)
```
