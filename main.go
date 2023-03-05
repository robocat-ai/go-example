package main

import (
	"log"
	"time"

	robocat "github.com/robocat-ai/robocat/pkg/client"
)

var flowName = "google"
var searchQuery = "Cute llamas"

func main() {
	container := createRobocatContainer("robocat", "demo")
	defer removeRobocatContainer()

	client, err := robocat.Connect(
		container.WebSocketAddress, container.Credentials,
	)
	if err != nil {
		log.Fatal(err)
	}

	client.SetSizeLimit("5M")

	// Push file "search_query" with contents of searchQuery byte slice.
	err = client.Input("search_query", []byte(searchQuery))
	if err != nil {
		log.Fatal(err)
	}

	// Start the flow asynchrously and get a reference to it.
	flow := client.Flow(flowName).WithTimeout(15 * time.Second).Run()

	flow.Log().Watch(func(line string) {
		log.Println("log:", line)
	})

	flow.Files().Watch(func(file *robocat.File) {
		log.Printf("got file: %s (%s)\n", file.Path, file.MimeType)
		writeFile(file) // Write file to the "output" directory.
	})

	// Wait for the flow to finish. This is a blocking operation.
	err = flow.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
