# Go's io.Reader

The `io.Reader` is defined as below:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

