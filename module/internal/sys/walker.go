package sys

import (
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"sync"
)

type walker struct {
	fs         afero.Fs
	wg         sync.WaitGroup
	mu         sync.RWMutex
	queue      []string
	restrictor chan struct{}
}

func newWalker(fs afero.Fs, parallel int) *walker {
	bf := &walker{
		fs:         fs,
		queue:      make([]string, 0),
		restrictor: make(chan struct{}, parallel),
	}
	return bf
}

func (bf *walker) pop() string {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	top := bf.queue[0]
	bf.queue = bf.queue[1:]
	return top
}

func (bf *walker) push(s string) {
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

	entries := bf.readdir(d)

	// Folder stat
	var count int64
	var size int64

	for _, entry := range entries {
		// Queue subdirs to walk in a queue
		if entry.isDir {
			bf.push(filepath.Join(d, entry.name))
		} else {
			// Send to channel
			fileEvent := filesystemItem{
				dir:   d,
				name:  entry.name,
				event: fsEventFile,
				count: 1,
				size:  entry.size,
			}
			results <- &fileEvent

			// update folder stat
			count++
			size += entry.size
		}
	}

	dirEvent := filesystemItem{
		dir:   d,
		event: fsEventDir,
		count: count,
		size:  size,
	}
	results <- &dirEvent
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

func (bf *walker) addOne() {
	bf.wg.Add(1)
	bf.acquireRestrict()
}

func (bf *walker) wait() {
	bf.wg.Wait()
}

func (bf *walker) readdir(path string) []*filesysEntry {
	defer bf.releaseRestrict()
	f, err := bf.fs.Open(path)
	if err != nil {
		return []*filesysEntry{}
	}
	defer Close(f)

	entries, err := f.Readdir(-1)
	if err != nil {
		return []*filesysEntry{}
	}

	result := make([]*filesysEntry, 0, len(entries))
	for _, e := range entries {
		// dont follow symlinks
		if e.Mode()&os.ModeSymlink == 0 {
			fi := filesysEntry{name: e.Name(), size: e.Size(), isDir: e.IsDir()}
			result = append(result, &fi)
		}
	}

	return result
}
