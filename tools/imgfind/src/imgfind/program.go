package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

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
	var searchTerm = args[0]
	if len(args) == 1 {
		host, _ = os.Hostname()
		remote = false
	} else {
		host = args[1]
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

	fmt.Println("Searching for ", searchTerm, "...")
	images, err := cli.ImageSearch(context.Background(), searchTerm, types.ImageSearchOptions{Limit: 100})
	if err != nil {
		panic(err)
	}

	var showErrors = false
	var wg sync.WaitGroup
	wg.Add(len(images))

	for _, image := range images {
		go func(imageName string) {
			inspectCmd := exec.Command("./linux/skopeo", "inspect", "docker://"+imageName)
			inspectOut, err := inspectCmd.CombinedOutput()
			if err != nil {
				if showErrors {
					fmt.Println(strings.TrimSpace(string(inspectOut)))
				}
				defer wg.Done()
				return
			}
			res := myImageInfo{}
			json.Unmarshal(inspectOut, &res)
			var isLinux = strings.Contains(res.Os, "linux")
			var isNano = false
			var isWin = false
			fmt.Print(" * ")
			for _, tag := range res.RepoTags {
				if !isNano {
					isNano = strings.Contains(tag, "nano")
				}
				if !isWin {
					isWin = strings.Contains(tag, "windows")
				}
			}
			var linuxStr = "[     ]"
			if isLinux {
				linuxStr = "[linux]"
			}
			var nanoStr = "[    ]"
			if isNano {
				nanoStr = "[nano]"
			}
			var winStr = "[       ]"
			if isWin {
				winStr = "[windows]"
			}
			fmt.Print(" " + linuxStr + " " + nanoStr + " " + winStr + " ")
			fmt.Println("'" + res.Name + "'")
			defer wg.Done()
		}(image.Name)
	}

	wg.Wait()
	fmt.Println("Done.")
}
