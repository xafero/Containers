package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dustin/go-humanize"
	"github.com/tmc/dot"
	"github.com/windler/dotgraph/renderer"
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

	fmt.Println()
	var g = dot.NewGraph("DockerImages")
	g.SetType(dot.DIGRAPH)

	hostNode := dot.NewNode("host")
	hostNode.Set("shape", "box")
	var simpleHost = strings.Split(strings.Split(host, "//")[1], ":")[0]
	hostNode.Set("label", "<<B>"+simpleHost+"</B>>")
	g.AddNode(hostNode)

	for i, image := range images {
		var id = image.ID[7:19]
		var size = humanize.Bytes(uint64(image.Size))
		var tag = image.RepoTags[0]
		fmt.Println(id, tag, size)

		var imgNode = dot.NewNode("image" + fmt.Sprintf("%d", i))
		imgNode.Set("shape", "ellipsis")
		imgNode.Set("label", strings.Replace(tag, ":", " #", -1))
		g.AddNode(imgNode)

		imgEdge := dot.NewEdge(hostNode, imgNode)
		g.AddEdge(imgEdge)
	}

	fmt.Println()
	var text = fmt.Sprint(g)
	r := &renderer.PNGRenderer{OutputFile: filepath.Join(pwd, "report.png")}
	r.Render(text)
}
