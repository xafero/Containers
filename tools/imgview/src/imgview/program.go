package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	pwd, _ := os.Getwd()
	fmt.Println("Current root :=", pwd)

	args := os.Args[1:]
	var host = args[0]
	const version = "v1.30"
	fmt.Println("Connecting to ", host, " with ", version, "...")

	cli, _ := client.NewClientWithOpts(client.WithHost(host), client.WithVersion(version))
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID)
		fmt.Println(image.RepoTags)
	}
}
