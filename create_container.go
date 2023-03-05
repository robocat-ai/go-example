package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	robocat "github.com/robocat-ai/robocat/pkg/client"
)

type RobocatContainer struct {
	WebSocketAddress string
	VNCAddress       string
	Credentials      robocat.Credentials
}

func removeRobocatContainer() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	container, found := pool.ContainerByName("robocat-go-example")
	if found {
		if err := pool.Purge(container); err != nil {
			log.Fatalf("Could not purge container: %s", err)
		}
	}

}

func createRobocatContainer(username string, password string) *RobocatContainer {
	removeRobocatContainer()

	result := &RobocatContainer{}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	credentials := robocat.Credentials{
		Username: username,
		Password: password,
	}

	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "robocat-go-example",
		Repository:   "ghcr.io/robocat-ai/robocat",
		ExposedPorts: []string{"80/tcp", "5900/tcp"},
		Env: []string{
			"VNC_ENABLED=1",
			fmt.Sprintf("AUTH_USERNAME=%s", credentials.Username),
			fmt.Sprintf("AUTH_PASSWORD=%s", credentials.Password),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.Mounts = append(config.Mounts, docker.HostMount{
			Target: "/flow",
			Source: fmt.Sprintf("%s/flow", pwd),
			Type:   "bind",
		})
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	result.WebSocketAddress = fmt.Sprintf("ws://%s", container.GetHostPort("80/tcp"))
	result.VNCAddress = container.GetHostPort("5900/tcp")
	result.Credentials = credentials

	log.Print("Waiting for container to start...")

	if err := pool.Retry(func() error {
		client, err := robocat.Connect(result.WebSocketAddress, credentials)
		if err != nil {
			return err
		}
		defer client.Close()
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to the server: %s", err)
	}

	return result
}
