package huffman

import (
	"fmt"
	"sort"
	"strings"
)

type Tree struct {
	Nodes []Node
}

type Node struct {
	leaf      bool
	freq      SymbolFreq
	value     Value
	left      int
	right     int
	binPrefix string
}

type (
	Value      rune
	SymbolFreq int
	Queue      []Node
)

func (q *Queue) CreateQueue(nodes map[Value]SymbolFreq) error {
	if len(nodes) == 0 {
		return fmt.Errorf("no nodes to create a queue")
	}
	var sortedNodes []Node
	for val, freq := range nodes {
		sortedNodes = append(sortedNodes, Node{value: val, freq: freq, leaf: true, left: -1, right: -1})
	}
	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].freq < sortedNodes[j].freq
	})
	*q = sortedNodes
	return nil
}

func GetSymbolFreq(data string) (map[Value]SymbolFreq, error) {
	if data == "" {
		return nil, fmt.Errorf("data was empty")
	}
	freqMap := make(map[Value]SymbolFreq)
	for _, char := range data {
		freqMap[Value(char)]++
	}
	return freqMap, nil
}

func (t *Tree) makeTree(pq *Queue) error {
	if len(*pq) < 2 {
		return fmt.Errorf("not enough nodes to make a tree")
	}
	for len(*pq) > 1 {
		leftNode, rightNode := (*pq)[0], (*pq)[1]
		*pq = (*pq)[2:]
		freqSum := leftNode.freq + rightNode.freq
		newNode := Node{
			freq:  freqSum,
			leaf:  false,
			left:  len(t.Nodes),
			right: len(t.Nodes) + 1,
		}
		t.Nodes = append(t.Nodes, leftNode, rightNode)
		*pq = append(*pq, newNode)
		sort.Slice(*pq, func(i, j int) bool {
			return (*pq)[i].freq < (*pq)[j].freq
		})
		fmt.Printf("Added internal node with freq: %d, left index: %d, right index: %d\n", freqSum, len(t.Nodes)-2, len(t.Nodes)-1)
	}
	if len(*pq) == 1 {
		t.Nodes = append(t.Nodes, (*pq)[0])
		fmt.Printf("Added root node with freq: %d\n", (*pq)[0].freq)
	}
	return nil
}

func (t *Tree) AssignBinaryPrefixes(index int, binPrefix string) {
	if index == -1 || index >= len(t.Nodes) {
		return
	}
	node := &t.Nodes[index]
	node.binPrefix = binPrefix
	if node.value != 0 {
		fmt.Printf("Assigning prefix %s to node with freq %d and value of %c \n", binPrefix, node.freq, rune(node.value))
	}
	if node.left != -1 {
		t.AssignBinaryPrefixes(node.left, binPrefix+"0")
	}
	if node.right != -1 {
		t.AssignBinaryPrefixes(node.right, binPrefix+"1")
	}
}

func (t *Tree) PrintTree() error {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print("\nHeres a visual representation of the tree\n\n")
	if len(t.Nodes) == 0 {
		return fmt.Errorf("tree is empty")
	}
	var printNode func(int, string)
	printNode = func(index int, prefix string) {
		if index == -1 || index >= len(t.Nodes) {
			return
		}
		node := t.Nodes[index]
		if node.leaf {
			fmt.Printf("%s└── %s (%c, %d)\n", prefix, node.binPrefix, rune(node.value), node.freq)
		} else {
			fmt.Printf("%s├── %s node (%d)\n", prefix, node.binPrefix, node.freq)
			newPrefix := prefix + "│   "
			printNode(node.left, newPrefix)
			if node.right != -1 {
				printNode(node.right, prefix+"    ")
			}
		}
	}
	printNode(len(t.Nodes)-1, "") // Start printing from the root
	return nil
}

func (t *Tree) EncodedMessage(input string) string {
	var binaryString string
	encodingMap := make(map[rune]string)
	for _, n := range t.Nodes {
		if n.leaf {
			encodingMap[rune(n.value)] = n.binPrefix
		}
	}
	for _, char := range input {
		binaryString += encodingMap[char]
	}
	return fmt.Sprint("Encoded message as binary value: ", binaryString)
}

func main() {
	var input string
	fmt.Println("Enter the text to encode:")
	fmt.Scanln(&input)

	freqMap, err := GetSymbolFreq(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	var pq Queue
	err = pq.CreateQueue(freqMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	var huffTree Tree
	err = huffTree.makeTree(&pq)
	if err != nil {
		fmt.Println(err)
		return
	}

	huffTree.AssignBinaryPrefixes(len(huffTree.Nodes)-1, "")

	err = huffTree.PrintTree()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(huffTree.EncodedMessage(input))
}
