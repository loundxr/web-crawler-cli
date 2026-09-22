package model

type Node struct {
	Resource string  `json:"resource"`
	Title    string  `json:"title"`
	Links    []*Node `json:"links"`
}

func NewNode(resource, title string) *Node {
	return &Node{
		Resource: resource,
		Title:    title,
		Links:    make([]*Node, 0),
	}
}
