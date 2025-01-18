package ioutil_test

import (
	"sync"
	"testing"

	"go.abhg.dev/io/ioutil"
)

func TestSyncWriter_race(t *testing.T) {
	const (
		NumWorkers = 10
		NumWrites  = 100
	)

	// This writer isn't thread-safe.
	// SyncWriter will make it so.
	var wroteBytes int
	counter := writerFunc(func(p []byte) {
		wroteBytes += len(p)
	})

	syncWriter := ioutil.SyncWriter(counter)

	var ready, done sync.WaitGroup
	ready.Add(NumWorkers)
	done.Add(NumWorkers)

	giveBytes := []byte("hello")

	for i := 0; i < NumWorkers; i++ {
		go func() {
			defer done.Done()

			ready.Done() // I'm ready.
			ready.Wait() // Are others ready?

			for i := 0; i < NumWrites; i++ {
				if _, err := syncWriter.Write(giveBytes); err != nil {
					t.Errorf("Write failed: %v", err)
					return
				}
			}
		}()
	}

	done.Wait()

	expectedBytes := NumWorkers * NumWrites * len(giveBytes)
	if wroteBytes != expectedBytes {
		t.Errorf("wrote %d bytes, expected %d", wroteBytes, expectedBytes)
	}
}

type writerFunc func([]byte)

func (f writerFunc) Write(p []byte) (int, error) {
	f(p)
	return len(p), nil
}
