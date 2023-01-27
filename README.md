# Go heap dump inspector

[![GoDoc](https://godoc.org/github.com/alexey-medvedchikov/go-heapview?status.svg)](https://pkg.go.dev/github.com/alexey-medvedchikov/go-heapview)

*This is a PoC software*

View contents of the heap dump file:

```shell
go run ./cmd/heapview/... dump heapdump.dat
```

Find pointers that are rooted in the stack frames and walk owned pointers to get the whole graph with total object count:

```shell
go run ./cmd/heapview/... owned heapdump.dat
```
