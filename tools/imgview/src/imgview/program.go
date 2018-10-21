package main

import (
	"fmt"
	"math"
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

	rels := make(map[string][]string)
	nodeIds := make(map[string]*dot.Node)
	nodeLayers := make(map[string][]string)

	for i, image := range images {
		var id = image.ID[7:19]
		var size = humanize.Bytes(uint64(image.Size))
		var tag = image.RepoTags[0]
		fmt.Println(id, tag, size)

		var info, _, _ = cli.ImageInspectWithRaw(context.Background(), image.ID)
		var layers = info.RootFS.Layers
		for _, layer := range layers {
			if v, found := rels[layer]; found {
				rels[layer] = append(v, image.ID)
			} else {
				rels[layer] = []string{image.ID}
			}
		}
		nodeLayers[image.ID] = layers

		var imgID = "image" + fmt.Sprintf("%d", i)
		var imgNode = dot.NewNode(imgID)
		imgNode.Set("shape", "ellipsis")
		imgNode.Set("label", strings.Replace(tag, ":", " #", -1))
		g.AddNode(imgNode)
		nodeIds[image.ID] = imgNode

		// imgEdge := dot.NewEdge(hostNode, imgNode)
		// g.AddEdge(imgEdge)
	}

	for id, layers := range nodeLayers {
		var myRels = make(map[string]int)
		var firstNode = nodeIds[id]
		for _, layer := range layers {
			var others = rels[layer]
			for _, other := range others {
				if other == id {
					continue
				}
				if v, found := myRels[other]; found {
					myRels[other] = v + 1
				} else {
					myRels[other] = 1
				}
			}
		}

		var bestScore = float64(-1)
		var bestMatch = float64(-1)
		var bestParent = "?"

		for k, v := range myRels {
			var count = len(nodeLayers[k])
			var match = float64(v) / float64(count)
			var score = float64(v) + match
			if math.Max(bestScore, score) == score {
				bestScore = score
				bestMatch = match
				bestParent = k
			}
		}

		var distantNode = nodeIds[bestParent]
		distantEdge := dot.NewEdge(firstNode, distantNode)
		if bestMatch < float64(1) {
			distantEdge.Set("style", "dashed")
			var text = fmt.Sprintf("%.0f", bestMatch*100)
			distantEdge.Set("label", text+" %")
		}
		g.AddEdge(distantEdge)
	}

	fmt.Println()
	var text = fmt.Sprint(g)
	r := &renderer.PNGRenderer{OutputFile: filepath.Join(pwd, "report.png")}
	r.Render(text)
}
