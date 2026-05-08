package llm

type _Heap struct {
	tg map[string]*Tree
}

type _HeapSpace string

var (
	_HeapSpaceTg _HeapSpace = "HEAP_SPACE_TG"
)

var Heap = _Heap{ tg: make(map[string]*Tree) }

func getTree(heapSpace _HeapSpace, treeId string) (*Tree, bool) {
	if treeId == "" {
		panic("tree id is empty")
	}

	var heapMap map[string]*Tree

	switch heapSpace {
	case _HeapSpaceTg: heapMap = Heap.tg
	}

	if tree, ok := heapMap[treeId]; !ok {
		return nil, false
	} else {
		return tree, true
	}
}

func createTree(heapSpace _HeapSpace, treeId string, node *TreeNode) *Tree {
	tree := NewTree(node)

	var heapMap map[string]*Tree
	switch heapSpace {
	case heapSpace: heapMap = Heap.tg
	default: panic("unreachable")
	}

	heapMap[treeId] = tree

	return tree
}

func (*_Heap) collectMessages(heapSpace _HeapSpace, treeId string, nodeId string) ([]Message, bool) {
	if nodeId == "" {
		panic("node id is empty")
	}
	tree, ok := getTree(heapSpace, treeId)
	if !ok {
		return nil, false
	}

	node, ok := tree.nodes[nodeId]
	if !ok {
		return nil, false
	}

	messages := tree.CollectMessages(node)

	return messages, true
}

func (*_Heap) containsNode(heapSpace _HeapSpace, treeId string, nodeId string) bool {
	if nodeId == "" {
		panic("node id is empty")
	}
	if tree, ok := getTree(heapSpace, treeId); ok {
		_, ok := tree.nodes[nodeId]
		return ok
	} else {
		return false
	}
}
