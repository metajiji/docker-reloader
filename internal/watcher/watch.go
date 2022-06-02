package watcher

import (
	"log"
	"time"

	"docker-reloader/internal/config"
	"docker-reloader/internal/docker"

	"github.com/docker/docker/client"
	"github.com/fsnotify/fsnotify"
)

func Watch(file config.WatchedFile, dClient client.APIClient, retried bool) {
	// File created, but no reload
	defer func() {
		if r := recover(); r != nil {
			log.Printf(`Restart inotify watcher for file %s with error: %s`, file.Path, r)
			time.Sleep(time.Second)
			Watch(file, dClient, true)
		}
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	err = watcher.Add(file.Path)
	if err != nil {
		panic(err)
	}

	if retried {
		log.Printf(`File "%s" has been appeared, trigger docker exec command for container "%s"`, file.Path, file.Container)
		docker.ExecCommand(dClient, file.Container, file.Command)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				panic(`File has been removed`)
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf(`File "%s" has been modified, trigger docker exec command for container "%s"`, file.Path, file.Container)
				docker.ExecCommand(dClient, file.Container, file.Command)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			panic(err)
		}
	}
}
