package heapview

import (
	"log"
	"os"
	"os/signal"
	"runtime/debug"
)

// SetupHandler sets up a handler to dump heap on certain signal
func SetupHandler(dumpFileFn func() string, signum ...os.Signal) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signum...)

	go func() {
		for range sigCh {
			writeHeapDump(dumpFileFn())
		}
	}()
}

func writeHeapDump(fpath string) {
	fp, err := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("could not create heap dump file %s: %+v", fpath, err)
	}
	defer func() {
		if closeErr := fp.Close(); closeErr != nil {
			log.Printf("could not close heap dump file %s: %+v", fpath, closeErr)
		}
	}()
	debug.WriteHeapDump(fp.Fd())
	log.Printf("heap dumped to %s", fpath)
}
