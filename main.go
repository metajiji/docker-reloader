package main

import (
	"log"

	"docker-reloader/internal/config"
	"docker-reloader/internal/docker"
	"docker-reloader/internal/options"
	"docker-reloader/internal/watcher"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	// Parse command line arguments
	args := options.ParseArgs(version, commit, date)

	// Load config file
	cfg := config.LoadConfig(string(args.Config))

	// Connect do docker daemon
	dClient := docker.DockerClient()

	log.Println(`Service started...`)

	ch := make(chan int, len(cfg.Watch))
	for _, wFile := range cfg.Watch {
		log.Println(`Spawn inotify watcher for file:`, wFile.Path)
		go watcher.Watch(wFile, dClient, false)
	}

	<-ch
}
