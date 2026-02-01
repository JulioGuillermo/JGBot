package tools

import "github.com/fsnotify/fsnotify"

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	OnChange func(event fsnotify.Event)
	OnError  func(err error)
}

func NewFileWatcher(filePath string) (*FileWatcher, error) {
	watcher := &FileWatcher{}
	var err error

	watcher.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.watcher.Add(filePath)
	if err != nil {
		return nil, err
	}

	go watcher.run()

	return watcher, nil
}

func (w *FileWatcher) run() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				w.Close()
				return
			}
			if w.OnChange != nil {
				w.OnChange(event)
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				w.Close()
				return
			}
			if w.OnError != nil {
				w.OnError(err)
			}
		}
	}
}

func (w *FileWatcher) Close() {
	w.watcher.Close()
}
