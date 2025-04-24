package gv

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileCallback func(path string)

type Watcher struct {
	watcher     *fsnotify.Watcher
	rootPath    string
	callback    FileCallback
	ignorePaths []string
	closeOnce   sync.Once
	done        chan struct{}
}

func NewWatcher(rootPath string, ignorePaths []string, cb FileCallback) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher:     w,
		rootPath:    rootPath,
		ignorePaths: ignorePaths,
		callback:    cb,
		done:        make(chan struct{}),
	}, nil
}

func (w *Watcher) isIgnored(path string) bool {
	for _, ignore := range w.ignorePaths {
		if strings.HasPrefix(path, ignore) {
			return true
		}
	}
	return false
}

func (w *Watcher) watchRecursive(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && !w.isIgnored(path) {
			err := w.watcher.Add(path)
			if err != nil {
				log.Println("Watcher add error:", err)
				return err
			}
		}
		return nil
	})
}

func (w *Watcher) listen() {
	defer log.Println("File watcher stopped")

	for {
		select {
		case <-w.done:
			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if w.isIgnored(event.Name) {
				log.Printf("Ignoring event for path: %s", event.Name)
				continue
			}

			// Handle directory creation
			if event.Op&fsnotify.Create == fsnotify.Create {
				// Wait a bit for the file system to stabilize
				time.Sleep(100 * time.Millisecond)

				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					log.Printf("Adding new directory to watch: %s", event.Name)
					if err := w.watchRecursive(event.Name); err != nil {
						log.Printf("Error adding directory to watch: %s - %v", event.Name, err)
					}
				}
			}

			// Handle file events
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				info, err := os.Stat(event.Name)
				if err == nil && !info.IsDir() {
					// Use a separate goroutine to avoid blocking the event loop
					go func(path string) {
						defer func() {
							if r := recover(); r != nil {
								log.Printf("Panic in callback: %v", r)
							}
						}()
						w.callback(path)
					}(event.Name)
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

func (w *Watcher) Start() error {
	if err := w.watchRecursive(w.rootPath); err != nil {
		return err
	}

	// Start listening in a separate goroutine
	go w.listen()

	return nil
}

func (w *Watcher) Close() error {
	var err error
	w.closeOnce.Do(func() {
		close(w.done)
		err = w.watcher.Close()
	})
	return err
}
