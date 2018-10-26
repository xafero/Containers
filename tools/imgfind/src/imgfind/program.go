package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type myImageInfo struct {
	Name          string            `json:"Name"`
	Digest        string            `json:"Digest"`
	RepoTags      []string          `json:"RepoTags"`
	Created       string            `json:"Created"`
	DockerVersion string            `json:"DockerVersion"`
	Labels        map[string]string `json:"Labels"`
	Architecture  string            `json:"Architecture"`
	Os            string            `json:"Os"`
	Layers        []string          `json:"Layers"`
}

func main() {
	pwd, _ := os.Getwd()
	fmt.Println("Current root :=", pwd)

	args := os.Args[1:]
	var host string
	var remote bool
	if len(args) == 0 {
		host, _ = os.Hostname()
		remote = false
	} else {
		host = args[0]
		remote = true
	}
	const version = "v1.30"
	fmt.Println("Connecting to ", host, " with ", version, "...")

	var cli *client.Client
	if remote {
		cli, _ = client.NewClientWithOpts(client.WithHost(host), client.WithVersion(version))
	} else {
		cli, _ = client.NewClientWithOpts(client.WithVersion(version))
	}
	images, err := cli.ImageSearch(context.Background(), "jenkins", types.ImageSearchOptions{Limit: 100})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		go func(imageName string) {
			inspectCmd := exec.Command("./linux/skopeo", "inspect", "docker://"+imageName)
			inspectOut, err := inspectCmd.CombinedOutput()
			if err != nil {
				fmt.Println(strings.TrimSpace(string(inspectOut)))
				return
			}
			res := myImageInfo{}
			json.Unmarshal(inspectOut, &res)
			fmt.Println(res)
		}(image.Name)
	}

	fmt.Scanln()
	fmt.Println("Done.")
}
