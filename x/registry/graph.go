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

// GraphResourceNode creates a node for a record.
func GraphResourceNode(g *dot.Graph, r Record) dot.Node {
	color := fmt.Sprintf("#%x", hash(r.Type)&0x00FFFFFF)
	node := g.Node(string(r.ID)).Attr("shape", "record").Attr("style", "").Attr("color", color)

	nodeLabel := fmt.Sprintf("%s | %s", string(r.ID)[:18], r.Type)
	if resourceLabel, ok := r.Attributes["label"].(string); ok {
		nodeLabel = fmt.Sprintf("%s | %s", nodeLabel, resourceLabel)
	}

	node.Attr("label", nodeLabel)

	// for _, link := range r.Links {
	// 	linkID := ""
	// 	linkLabel := ""

	// 	if idAttr, ok := link["id"].(string); ok {
	// 		linkID = idAttr
	// 	}

	// 	if labelAttr, ok := link["label"].(string); ok {
	// 		linkLabel = labelAttr
	// 	}

	// 	if linkID != "" {
	// 		g.Edge(
	// 			g.Node(string(r.ID)),
	// 			g.Node(string(linkID)),
	// 			linkLabel,
	// 		)
	// 	}
	// }

	return node
}
