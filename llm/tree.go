package llm

import (
	"slices"

	"d0ctr/bldbr-bot/llm/types"

	"github.com/google/uuid"
)

type Tree struct {
	nodes map[string]*TreeNode
	id string
}

func NewTree(root *TreeNode) *Tree {
	id := uuid.New().String()
	return &Tree{ map[string]*TreeNode{root.id(): root}, id }
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

func (t *Tree) UpdateCursor(node *TreeNode, cursor string) {
	node.cursor = cursor
	node = node.prev
	for node != nil {
		node.cursor = ""
	}
}

func (t *Tree) GetId() string {
	return t.id
}

type TreeNode struct {
	message types.Message
	prev *TreeNode
	// next []*TreeNode
	cursor string // this is effectively the response id that got us this node, if there is none, then the full context should be sent
}

func NewTreeNode(message types.Message) *TreeNode {
	return &TreeNode{ message, nil, "" }
}

func (n TreeNode) id() string {
	return n.message.Id()
}

func (n *TreeNode) append(next *TreeNode) {
	next.prev = n
}
