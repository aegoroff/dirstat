package scan

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type walker struct {
	fs         Filesystem
	wg         sync.WaitGroup
	mu         sync.RWMutex
	queue      []string
	restrictor chan struct{}
}

func newWalker(fs Filesystem, parallel int) *walker {
	bf := &walker{
		fs:         fs,
		queue:      make([]string, 0),
		restrictor: make(chan struct{}, parallel),
	}
	return bf
}

var linuxPrefixesToSkip = []string{
	"/proc",
	"/dev",
	"/sys",
}

func (bf *walker) dequeue() string {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	top := bf.queue[0]
	bf.queue = bf.queue[1:]
	return top
}

func (bf *walker) enqueue(s string) {
	bf.mu.Lock()
	bf.queue = append(bf.queue, s)
	bf.mu.Unlock()
}

func (bf *walker) len() int {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return len(bf.queue)
}

func (bf *walker) walk(d string, results chan<- *filesystemItem) {
	defer bf.wg.Done()

	// Process directory entries in chunks to reduce memory usage
	const chunkSize = 128
	var count int64
	var size int64

	// Open directory
	defer bf.releaseRestrict()
	f, err := bf.fs.Open(d)
	if err != nil {
		return
	}
	defer Close(f)

	// Process directory entries in chunks
	for {
		entries, err := f.Readdir(chunkSize)
		if err != nil {
			// Even if we have an error, we still need to send the directory event
			// to maintain consistency in the event stream
			break
		}

		// Process chunk of entries
		for _, entry := range entries {
			m := entry.Mode()
			// Skip symlinks
			if m&os.ModeSymlink != 0 {
				continue
			}

			path := filepath.Join(d, entry.Name())
			if bf.skip(path) {
				continue
			}
			// Queue subdirs to walk in a queue
			if m.IsDir() {
				bf.enqueue(path)
			} else if m.IsRegular() {
				// Send to channel
				fileEvent := filesystemItem{
					path:  path,
					event: fsEventFile,
					size:  entry.Size(),
				}
				results <- &fileEvent

				// update folder stat
				count++
				size += fileEvent.size
			}
		}

		// If we got less than chunkSize entries, we've reached the end of directory
		if len(entries) < chunkSize {
			break
		}
	}

	dirEvent := filesystemItem{
		path:  d,
		event: fsEventDir,
		count: count,
		size:  size,
	}
	results <- &dirEvent
}

func (bf *walker) skip(path string) bool {
	if runtime.GOOS != "linux" {
		return false
	}
	for _, prefix := range linuxPrefixesToSkip {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func (bf *walker) acquireRestrict() {
	bf.restrictor <- struct{}{}
}

func (bf *walker) releaseRestrict() {
	<-bf.restrictor
}

func (bf *walker) closeRestrict() {
	close(bf.restrictor)
}

func (bf *walker) start() {
	bf.wg.Add(1)
	bf.acquireRestrict()
}

func (bf *walker) wait() {
	bf.wg.Wait()
}
