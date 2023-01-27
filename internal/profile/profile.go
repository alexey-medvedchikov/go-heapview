package profile

import (
	"log"
	"os"
	"runtime/pprof"
)

// CPU starts CPU profiling and saves the statistics into a file specified by profilePath. It overwrites the file
// if it already exists. To finalize profiling caller needs to call finish function returned.
func CPU(profilePath string) (finish func(), err error) {
	fp, err := os.OpenFile(profilePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	if err := pprof.StartCPUProfile(fp); err != nil {
		return nil, err
	}

	return func() {
		pprof.StopCPUProfile()

		if err := fp.Close(); err != nil {
			log.Printf("could not close CPU profile file %s: %v", profilePath, err)
		}
	}, nil
}

// Heap starts heap profiling and saves the statistics into a file specified by profilePath. It overwrites the file
// if it already exists. To finalize profiling caller needs to call finish function returned.
func Heap(profilePath string) (finish func(), err error) {
	fp, err := os.OpenFile(profilePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return func() {
		if err := pprof.WriteHeapProfile(fp); err != nil {
			log.Printf("could not write heap profile to %s: %v", profilePath, err)
		}

		if err := fp.Close(); err != nil {
			log.Printf("could not close heap profile file %s: %v", profilePath, err)
		}
	}, nil
}
