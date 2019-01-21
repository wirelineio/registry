//
// Copyright 2019 Wireline, Inc.
//

package utxo

import (
	"fmt"
	"strings"

	"github.com/emicklei/dot"
)

// LabelLength is leading number of chars from the hash used in the dot graph.
const LabelLength int = 8

// TxNodeLabel returns the node label for a transaction.
func TxNodeLabel(txID Hash) string {
	return txID.String()[:LabelLength]
}

// TxInNodeID generates a dot graph ID for a Tx input.
func TxInNodeID(txID Hash, in TxIn) string {
	return fmt.Sprintf("TXIN_%s_%d", txID, in.Input.Index)
}

// TxNode creates a node for a transaction.
func TxNode(g *dot.Graph, txNodeID Hash, tx Tx) dot.Node {
	node := g.Node(txNodeID.String()).Attr("shape", "record").Attr("style", "").Attr("color", "red")

	var txLabel []string
	txLabel = append(txLabel, "TXN")
	txLabel = append(txLabel, TxNodeLabel(txNodeID))

	for _, txIn := range tx.TxIn {
		g.Edge(
			g.Node(txNodeID.String()),
			g.Node(txIn.Input.Hash.String()),
			TxInLabel(txIn),
		)
	}

	for txOutputIndex, txOut := range tx.TxOut {
		txLabel = append(txLabel, TxOutLabel(txOutputIndex, txOut))
	}

	node.Attr("label", strings.Join(txLabel, " | "))

	return node
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
	return fmt.Sprintf("A/C OUT | %s | VAL = %d", accOut.ID.String()[:LabelLength], accOut.Value)
}

// AccOutNode creates a node for an account output.
func AccOutNode(g *dot.Graph, accOut AccOutput) dot.Node {
	return g.Node(accOut.ID.String()).Attr("shape", "record").Attr("color", "blue").Attr("style", "").Attr("label", AccOutLabel(accOut))
}

// UnspentOutputNode creates a node for an UTXO.
func UnspentOutputNode(g *dot.Graph, utxo OutPoint) dot.Node {
	utxoNodeID := fmt.Sprintf("UTXO_%s_%d", utxo.Hash, utxo.Index)

	node := g.Node(utxoNodeID)
	node.Attr("color", "black").Attr("fillcolor", "green").Attr("shape", "oval").Attr("label", "UTXO").Attr("style", "filled")

	g.Edge(
		g.Node(utxoNodeID),
		g.Node(utxo.Hash.String()),
		fmt.Sprintf("REF(OUT #%d)", utxo.Index),
	)

	return node
}
