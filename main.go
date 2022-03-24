package main

import (
	"context"
	"io"
	"log"
	"os"
	"runtime/pprof"

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

	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)

	ctx, cancel := context.WithCancel(context.Background())
	events, errs := dockerClient.Events(ctx, types.EventsOptions{})
	defer cancel()

	if err := listners.init(); err != nil {
		log.Printf("listners.init returns error: %+v\n", err)
		return err
	}

	log.Println("Listening")

	for {
		select {
		case event := <-events:
			err := listners.receive(event)
			if err != nil {
				log.Printf("listner returns error: %+v\n", err)
			}

		case err := <-errs:
			if err == io.EOF {
				return nil
			}
			log.Printf("docker client receives error: %+v\n", err)
			return err
		}
	}
}
