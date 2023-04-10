package main

import "sort"

const (
	TYPE_LEAF     = 0
	TYPE_INTERIOR = 1
)

type BTree struct {
	root    *Node
	maxKeys int
}

type Node struct {
	tree   *BTree
	parent *Node
	typ    int
	keys   []int
	childs []*Node
}

func (t *BTree) Insert(key int) {
	if t.root == nil {
		t.root = NewNode(t, TYPE_LEAF)
	}

	t.root.insert(key)
}

func NewNode(tree *BTree, typ int) *Node {
	n := &Node{tree: tree, typ: typ}
	n.keys = []int{}
	n.childs = []*Node{}

	return n
}

func (n *Node) insert(key int) {
	if n.typ == TYPE_INTERIOR {
		n.insertInterior(key)
	} else {
		n.insertLeaf(key)
	}
}

func (n *Node) insertInterior(key int) {
	// find the child to insert into
	var child *Node
	for i, k := range n.keys {
		if key < k {
			child = n.childs[i]
			break
		}
	}

	if child == nil {
		// insert into the last child
		child = n.childs[len(n.childs)-1]
	}

	child.insert(key)

	if len(n.keys) > n.tree.maxKeys {
		n.splitInterior()
	}
}

func (n *Node) splitInterior() {
	if n.parent == nil {
		createParent(n)
	}

	// split the node
	// current node will contain the first half of the keys
	// new node will contain the second half of the keys
	midKeys := len(n.keys) / 2
	midChilds := len(n.childs) / 2
	splitKey := n.keys[midKeys]

	// set right node
	newNode := NewNode(n.tree, TYPE_INTERIOR)
	newNode.keys = n.keys[midKeys+1:]
	newNode.childs = n.childs[midChilds:]
	newNode.parent = n.parent

	for _, childNode := range newNode.childs {
		childNode.parent = newNode
	}

	// update left node
	n.keys = n.keys[:midKeys]
	n.childs = n.childs[:midChilds]

	index := sort.SearchInts(n.parent.keys, splitKey)
	n.parent.keys = append(n.parent.keys[:index], append([]int{splitKey}, n.parent.keys[index:]...)...)
	n.parent.childs = append(n.parent.childs[:index+1], append([]*Node{newNode}, n.parent.childs[index+1:]...)...)
}

func (n *Node) insertLeaf(key int) {
	// check if the key already exists
	for _, k := range n.keys {
		if k == key {
			return
		}
	}

	index := sort.SearchInts(n.keys, key)
	newKeys := make([]int, len(n.keys)+1)
	copy(newKeys[:index], n.keys[:index])
	newKeys[index] = key
	copy(newKeys[index+1:], n.keys[index:])
	n.keys = newKeys

	// check if we need to split
	if len(n.keys) > n.tree.maxKeys {
		n.split()
	}
}

func (n *Node) split() {
	if n.parent == nil {
		createParent(n)
	}

	mid := len(n.keys) / 2

	newNode := NewNode(n.tree, TYPE_LEAF)
	newNode.keys = n.keys[mid:]
	newNode.parent = n.parent
	if len(n.childs) > 0 {
		newNode.childs = n.childs[mid:]
		for _, childNode := range newNode.childs {
			childNode.parent = newNode
		}
	}

	n.keys = n.keys[:mid]
	if len(n.childs) > 0 {
		n.childs = n.childs[:mid]
	}

	splitKey := newNode.keys[0]
	index := sort.SearchInts(n.parent.keys, splitKey)
	n.parent.keys = append(n.parent.keys[:index], append([]int{splitKey}, n.parent.keys[index:]...)...)
	n.parent.childs = append(n.parent.childs[:index+1], append([]*Node{newNode}, n.parent.childs[index+1:]...)...)
}

func createParent(n *Node) {
	parent := NewNode(n.tree, TYPE_INTERIOR)
	parent.childs = append(parent.childs, n)

	n.parent = parent
	n.tree.root = parent
}

func main() {
	tree := BTree{maxKeys: 2}

	tree.Insert(1)
	tree.Insert(6)
	tree.Insert(7)
	tree.Insert(8)
	//tree.Insert(9)
	//tree.Insert(2)
	//tree.Insert(3)

	printTree(tree)
}
