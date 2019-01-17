//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"fmt"

	"github.com/emicklei/dot"
	"github.com/wirelineio/wirechain/x/utxo/utils"
)

// TxInNodeID generates a dot graph ID for a Tx input.
func TxInNodeID(txID []byte, in TxIn) string {
	return fmt.Sprintf("TXIN_%s_%d", utils.BytesToBase64(txID), in.Input.Index)
}

// TxNode creates a node for a transaction.
func TxNode(g *dot.Graph, txID string) dot.Node {
	return g.Node(txID).Attr("shape", "record").Attr("style", "").Attr("color", "red")
}

// TxInLabel returns the edge label for a transaction input.
func TxInLabel(txIn TxIn) string {
	return fmt.Sprintf("IN(OUT #%d)", txIn.Input.Index)
}

// TxOutLabel returns the label for a transaction output.
func TxOutLabel(txOutIndex int, txOut TxOut) string {
	return fmt.Sprintf("OUT #%d = %d", txOutIndex, txOut.Value)
}

// AccOutLabel returns the label for an account output.
func AccOutLabel(accOut AccOutput) string {
	return fmt.Sprintf("A/C OUT | %s | VAL = %d", utils.BytesToBase64(accOut.ID), accOut.Value)
}

// AccOutNode creates a node for an account output.
func AccOutNode(g *dot.Graph, accOut AccOutput) dot.Node {
	accOutNodeID := utils.BytesToBase64(accOut.ID)
	return g.Node(accOutNodeID).Attr("shape", "record").Attr("color", "blue").Attr("style", "").Attr("label", AccOutLabel(accOut))
}

// UnspentOutputNode creates a node for an UTXO.
func UnspentOutputNode(g *dot.Graph, utxo OutPoint) dot.Node {
	utxoNodeID := fmt.Sprintf("UTXO_%s_%d", utils.BytesToBase64(utxo.Hash), utxo.Index)

	node := g.Node(utxoNodeID)
	node.Attr("color", "black").Attr("fillcolor", "green").Attr("shape", "oval").Attr("label", "UTXO").Attr("style", "filled")

	g.Edge(
		g.Node(utxoNodeID),
		g.Node(utils.BytesToBase64(utxo.Hash)),
		fmt.Sprintf("REF(OUT #%d)", utxo.Index),
	)

	return node
}
