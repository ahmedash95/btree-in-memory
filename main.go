package main

import (
	"math/rand"
	"sort"
	"time"
)

const (
	TYPE_LEAF     = 0
	TYPE_INTERIOR = 1
)

type BTree struct {
	root        *Node
	maxKeys     int
	callOnSplit func()
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

	if n.tree.callOnSplit != nil {
		n.tree.callOnSplit()
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

	// splitLeaf the node
	// current node will contain the first half of the keys
	// new node will contain the second half of the keys
	midKey := n.keys[len(n.keys)/2]
	n.parent.addKey(midKey)

	// set right node
	newNode := NewNode(n.tree, TYPE_INTERIOR)
	newNode.parent = n.parent

	n.keys, newNode.keys = n.splitTwoKeys()
	n.childs, newNode.childs = n.splitTwoChilds()

	// splitting internal node is a bit different
	// we take one key out of the sibling node and put it in the parent node
	if len(newNode.keys) > 1 {
		newNode.keys = newNode.keys[1:]
	}

	for _, childNode := range newNode.childs {
		childNode.parent = newNode
	}

	newNode.parent.addChild(newNode)
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

	// check if we need to splitLeaf
	if len(n.keys) > n.tree.maxKeys {
		n.splitLeaf()
	}
}

func (n *Node) splitLeaf() {
	if n.parent == nil {
		createParent(n)
	}

	newNode := NewNode(n.tree, TYPE_LEAF)
	newNode.parent = n.parent

	n.keys, newNode.keys = n.splitTwoKeys()

	splitKey := newNode.keys[0]
	n.parent.addKey(splitKey)
	n.parent.addChild(newNode)
}

func (n *Node) splitTwoKeys() ([]int, []int) {
	mid := len(n.keys) / 2
	leftKeys := make([]int, mid)
	rightKeys := make([]int, len(n.keys)-mid)

	copy(leftKeys, n.keys[:mid])
	copy(rightKeys, n.keys[mid:])

	return leftKeys, rightKeys
}

func (n *Node) splitTwoChilds() ([]*Node, []*Node) {

	var mid int
	if len(n.childs)%2 != 0 { // odd
		mid = (len(n.childs) / 2) + 1
	} else {
		mid = len(n.childs) / 2
	}

	leftChilds := make([]*Node, mid)
	rightChilds := make([]*Node, len(n.childs)-mid)

	copy(leftChilds, n.childs[:mid])
	copy(rightChilds, n.childs[mid:])

	return leftChilds, rightChilds
}

func (n *Node) addKey(key int) {
	newKeys := make([]int, len(n.keys)+1)
	index := sort.SearchInts(n.keys, key)
	copy(newKeys, n.keys[:index])
	newKeys[index] = key
	copy(newKeys[index+1:], n.keys[index:])

	n.keys = newKeys
}

func (n *Node) addChild(node *Node) {
	newChilds := make([]*Node, len(n.childs)+1)
	// search for index by child first key
	index := sort.Search(len(n.childs), func(i int) bool { return n.childs[i].keys[0] >= node.keys[0] })
	copy(newChilds, n.childs[:index])
	newChilds[index] = node
	copy(newChilds[index+1:], n.childs[index:])

	n.childs = newChilds
}

func (n *Node) removeKey(key int) {
	newKeys := make([]int, len(n.keys)-1)
	for i := 0; i < len(n.keys); i++ {
		if n.keys[i] == key {
			continue
		}

		newKeys = append(newKeys, n.keys[i])
	}

	n.keys = newKeys
}

func createParent(n *Node) {
	parent := NewNode(n.tree, TYPE_INTERIOR)
	parent.childs = append(parent.childs, n)

	n.parent = parent
	n.tree.root = parent
}

func main() {
	tree := BTree{maxKeys: 5}

	var html []string
	tree.callOnSplit = func() {
		newTree := MermaidHtml(tree)
		if len(html) > 0 && html[len(html)-1] == newTree {
			return
		}

		html = append(html, newTree)
	}

	var nums []int
	for i := 1; i <= 25; i++ {
		nums = append(nums, i)
	}

	nums = shuffle(nums)

	for _, num := range nums {
		tree.Insert(num)
	}

	mermaidToHtml(html)
}

func shuffle(nums []int) []int {
	// shuffle the numbers
	rand.Seed(time.Now().UnixNano())
	for i := range nums {
		j := rand.Intn(i + 1)
		nums[i], nums[j] = nums[j], nums[i]
	}

	return nums
}
