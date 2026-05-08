package llm

var heap struct {
	tg map[string]*Tree
}

func init() {
	heap.tg = make(map[string]*Tree)
}

