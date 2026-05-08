package llm

import (
	"slices"

	"github.com/d0ctr/bldbr-bot-go/llm/types"
)

type Tree struct {
	nodes map[string]*TreeNode
}

func NewTree(root *TreeNode) *Tree {
	return &Tree{ map[string]*TreeNode{root.id(): root} }
}

func (t *Tree) AddNode(node *TreeNode) {
	t.nodes[node.id()] = node
}

func (t *Tree) AppendNode(prev *TreeNode, node *TreeNode) bool {
	if prev == nil {
		panic("empty prev node")
	} else if node == nil {
		panic("empty node")
	}

	// wire node to a parent node
	prev.append(node)

	// add the node to the node map
	t.nodes[node.id()] = node

	return true
}

func (t *Tree) GetNode(nodeId string) *TreeNode {
	if node, ok := t.nodes[nodeId]; ok {
		return node
	} else {
		return nil
	}
}

func (t *Tree) CollectMessages(node *TreeNode) []types.Message {
	if node == nil {
		panic("empty node")
	}
	messages := []types.Message{ node.message }

	for node.prev != nil {
		node = node.prev
		messages = append(messages, node.message)
	}

	slices.Reverse(messages)

	return messages
}


type TreeNode struct {
	message types.Message
	prev *TreeNode
	// next []*TreeNode
}

func NewTreeNode(message types.Message) *TreeNode {
	return &TreeNode{ message, nil }
}

func (n TreeNode) id() string {
	return n.message.Id()
}

func (n *TreeNode) append(next *TreeNode) {
	next.prev = n
}

