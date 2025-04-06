package core

import "crypto/sha256"

// MerkleTree represent a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	for range len(data) / 2 {
		var newLevel []MerkleNode

		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(nodes, *node)
		}

		nodes = newLevel
	}

	mTree := MerkleTree{RootNode: &nodes[0]}

	return &mTree
}

// MerkleNode represent a Merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left *MerkleNode, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{Left: left, Right: right}

	hash := sha256.Sum256(data)
	if left != nil || right != nil {
		prevHashes := append(left.Data, right.Data...)
		hash = sha256.Sum256(prevHashes)
	}

	mNode.Data = hash[:]

	return &mNode
}
