package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func (t *BTree) Mermaid() string {
	out := fmt.Sprintln("graph TD;")
	out += t.root.Mermaid("", "Tree")

	return out
}

func (n *Node) Mermaid(oldPrefix, prefix string) string {
	var output string
	var nodeType string
	if n.typ == TYPE_LEAF {
		nodeType = "L"
	} else {
		nodeType = "I"
	}

	//assign pointer id to nodetype
	//nodeType = fmt.Sprintf("%p", n)[8:]
	nodeType = "" // keep it empty for now

	if len(nodeType) > 0 {
		nodeType = fmt.Sprintf("-%s-", nodeType)
	}

	keysStr := ""
	for _, key := range n.keys {
		keysStr += fmt.Sprintf("%d ", key)
	}
	output += fmt.Sprintf("%s(%s)\n", prefix, keysStr)
	if oldPrefix != "" {
		output += fmt.Sprintf("%s -%s-> %s\n", oldPrefix, nodeType, prefix)
	}

	if n.typ == TYPE_INTERIOR {
		for _, child := range n.childs {
			output += child.Mermaid(prefix, fmt.Sprintf("%s_%d", prefix, child.keys[0]))
		}
	}

	return output
}

func printTree(tree BTree) {
	result := tree.Mermaid()
	// write to file
	err := os.WriteFile("tree.mermaid", []byte(result), 0644)
	if err != nil {
		panic(err)
	}

	// execute mermaid cli
	// mmdc -i tree.mermaid -o tree.png

	cmd := exec.Command("mmdc", "-i", "tree.mermaid", "-o", "tree.svg")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	// open tree.png
	cmd = exec.Command("open", "tree.svg")
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func MermaidHtml(tree BTree) string {
	result := tree.Mermaid()

	// wrap the string in "<div class="mermaid">" and "</div>"
	return fmt.Sprintf("<div class=\"mermaid active\">\n%s\n</div>", result)
}

func mermaidToHtml(arr []string) {
	content, err := ioutil.ReadFile("mermaid.html")
	if err != nil {
		panic(err)
	}

	htmlString := string(content)
	startTag := `<section id="list">`
	endTag := `</section>`

	startIdx := strings.Index(htmlString, startTag)
	endIdx := strings.Index(htmlString, endTag)

	if startIdx == -1 || endIdx == -1 {
		panic(fmt.Sprintf("section with ID 'list' not found"))
	}

	startIdx += len(startTag)
	newContent := strings.Join(arr, "\n")

	newHTMLString := htmlString[:startIdx] + newContent + htmlString[endIdx:]

	err = os.WriteFile("mermaid.html", []byte(newHTMLString), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
