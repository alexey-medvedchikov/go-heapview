# 01-leaky-slice

In this example we have a leak via slice buried in a some kind of client (maybe a cache).

To produce heap dump, you can increase number of top-level objects allocated by supplying integer argument, the default is 1000.
If you increase the number of objects you need to raise sleep times or do this manually.

```shell
go build -gcflags='all=-N -l' -o bin/01-leaky-slice ./example/01-leaky-slice/...
./bin/01-leaky-slice &
pid=$!
sleep 1
kill -s USR1 "$pid"
sleep 5
kill -s TERM "$pid"
```

After the heap dump is produced you can run `heapview` to inspect the file:

```plain
❯ go run ./cmd/heapview/ owned heapdump.dat
{"Address":824633745824,"OwnSize":416,"OwnedSize":208,"OwnedCount":3,"Frames":[{"Address":824634154328,"FuncName":"runtime.chanrecv"}]}
{"Address":824634221024,"OwnSize":96,"OwnedSize":16,"OwnedCount":1,"Frames":[{"Address":824634425216,"FuncName":"github.com/alexey-medvedchikov/go-heapview.SetupHandler.func1"}]}
{"Address":824636440576,"OwnSize":10240,"OwnedSize":1048000,"OwnedCount":2000,"Frames":[{"Address":824634154512,"FuncName":"main.main"}]}
{"Address":824633750400,"OwnSize":416,"OwnedSize":384,"OwnedCount":4,"Frames":[{"Address":824634281464,"FuncName":"runtime.selectgo"}]}
{"Address":824634220544,"OwnSize":96,"OwnedSize":528,"OwnedCount":3,"Frames":[{"Address":824634154328,"FuncName":"runtime.chanrecv"}]}
```

The last object is the one we're interested in. It is an array that backs `Client.cache` slice.