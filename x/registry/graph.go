//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"fmt"
	"hash/fnv"

	"github.com/emicklei/dot"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// GraphResourceNode creates a node for a resource.
func GraphResourceNode(g *dot.Graph, r ResourceYaml) dot.Node {
	color := fmt.Sprintf("#%x", hash(r.Type)&0x00FFFFFF)
	node := g.Node(string(r.ID)).Attr("shape", "record").Attr("style", "").Attr("color", color)

	nodeLabel := fmt.Sprintf("%s | %s", string(r.ID)[:18], r.Type)
	if resourceLabel, ok := r.Attributes["label"].(string); ok {
		nodeLabel = fmt.Sprintf("%s | %s", nodeLabel, resourceLabel)
	}

	node.Attr("label", nodeLabel)

	for linkID, linkData := range r.Links {
		linkLabel := ""
		if linkAttrs, ok := linkData.(map[string]interface{}); ok {
			if labelAttr, ok := linkAttrs["label"].(string); ok {
				linkLabel = labelAttr
			}
		}

		g.Edge(
			g.Node(string(r.ID)),
			g.Node(string(linkID)),
			linkLabel,
		)
	}

	return node
}
