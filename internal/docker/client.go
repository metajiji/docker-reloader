package docker

import (
	"context"
	"log"

	"github.com/docker/docker/client"
)

func DockerClient() client.APIClient {
	dClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	i, err := dClient.Info(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(`Connected to dockerd version "%s" on machine "%s" with OS "%s"`, i.ServerVersion, i.Name, i.OperatingSystem)

	// TODO:Checkif we are connected successfully?

	return dClient
}
