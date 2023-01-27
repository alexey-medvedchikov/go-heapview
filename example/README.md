# 01-leaky-slice

In this example we have a leak via slice buried in a some kind of client (maybe a cache).

To produce heap dump, you can increase number of top-level objects allocated by supplying integer argument, the default is 1000.
If you increase the number of objects you need to raise sleep times or do this manually.

```shell
go build -o bin/01-leaky-slice ./example/01-leaky-slice/...
./bin/01-leaky-slice &
pid=$!
sleep 1
kill -s USR1 "$pid"
sleep 5
kill -s TERM "$pid"
```

After the heap dump is produced you can run heapview to inspec the file:

```plain
❯ go run ./cmd/heapview/ owned heapdump.dat
pointer 0xc0000061a0
  own size 416
  owned size 208
  owned count 3
  found in frames
    0: c000104df808x runtime.chanrecv
pointer 0xc00007a000
  own size 96
  owned size 528
  owned count 3
  found in frames
    0: c000104df808x runtime.chanrecv
pointer 0xc000007380
  own size 416
  owned size 384
  owned count 4
  found in frames
    0: c000100e1808x runtime.selectgo
pointer 0xc00007a1e0
  own size 96
  owned size 16
  owned count 1
  found in frames
    0: c00009bfa008x github.com/alexey-medvedchikov/go-heapview.SetupHandler.func1
pointer 0xc000308000
  own size 10240
  owned size 1048000
  owned count 2000
  found in frames
    0: c000104eb008x main.main
```

The last object is the one we're interested in. It is an array that backs `Client.cache` slice.