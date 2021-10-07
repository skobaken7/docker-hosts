package main

import (
	"context"
	"io"
	"log"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return errors.Wrap(err, "failed to instantiate Listners")
	}

	listners := newListeners(dockerClient)

	wait := make(chan error)
	go (func() {
		log.Println("Listening")

		events, errs := dockerClient.Events(context.Background(), types.EventsOptions{})
		for {
			select {
			case event := <-events:
				err := listners.receive(event)
				if err != nil {
					log.Printf("listner returns error: %+v\n", err)
				}

			case err := <-errs:
				wait <- err
				return

			default:
			}
		}
	})()

	if err := listners.init(); err != nil {
		log.Printf("listners.init returns error: %+v\n", err)
	}

	if err := <-wait; err != io.EOF {
		log.Printf("events stream catch the error: %+v\n", err)
	}

	return nil
}
