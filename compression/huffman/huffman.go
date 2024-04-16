package huffman

import (
	"fmt"
	"sort"
	"strings"
)

type tree struct {
	nodes []node
}

type node struct {
	leaf      bool
	freq      symbolFreq
	value     value
	left      int
	right     int
	binPrefix string
}

type value rune
type symbolFreq int
type queue []node

func (q *queue) createQueue(nodes map[value]symbolFreq) error {
	if len(nodes) == 0 {
		return fmt.Errorf("no nodes to create a queue")
	}
	var sortedNodes []node
	for val, freq := range nodes {
		sortedNodes = append(sortedNodes, node{value: val, freq: freq, leaf: true, left: -1, right: -1})
	}
	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].freq < sortedNodes[j].freq
	})
	*q = sortedNodes
	return nil
}

func getSymbolFreq(data string) (map[value]symbolFreq, error) {
	if data == "" {
		return nil, fmt.Errorf("data was empty")
	}
	freqMap := make(map[value]symbolFreq)
	for _, char := range data {
		freqMap[value(char)]++
	}
	return freqMap, nil
}

func (t *tree) makeTree(pq *queue) error {
	if len(*pq) < 2 {
		return fmt.Errorf("not enough nodes to make a tree")
	}
	for len(*pq) > 1 {
		leftNode, rightNode := (*pq)[0], (*pq)[1]
		*pq = (*pq)[2:]
		freqSum := leftNode.freq + rightNode.freq
		newNode := node{
			freq:  freqSum,
			leaf:  false,
			left:  len(t.nodes),
			right: len(t.nodes) + 1,
		}
		t.nodes = append(t.nodes, leftNode, rightNode)
		*pq = append(*pq, newNode)
		sort.Slice(*pq, func(i, j int) bool {
			return (*pq)[i].freq < (*pq)[j].freq
		})
		fmt.Printf("Added internal node with freq: %d, left index: %d, right index: %d\n", freqSum, len(t.nodes)-2, len(t.nodes)-1)
	}
	if len(*pq) == 1 {
		t.nodes = append(t.nodes, (*pq)[0])
		fmt.Printf("Added root node with freq: %d\n", (*pq)[0].freq)
	}
	return nil
}

func (t *tree) assignBinaryPrefixes(index int, binPrefix string) {
	if index == -1 || index >= len(t.nodes) {
		return
	}
	node := &t.nodes[index]
	node.binPrefix = binPrefix
	if node.value != 0 {
		fmt.Printf("Assigning prefix %s to node with freq %d and value of %c \n", binPrefix, node.freq, rune(node.value))
	}
	if node.left != -1 {
		t.assignBinaryPrefixes(node.left, binPrefix+"0")
	}
	if node.right != -1 {
		t.assignBinaryPrefixes(node.right, binPrefix+"1")
	}
}

func (t *tree) printTree() error {
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print("\nHeres a visual representation of the tree\n\n")
	if len(t.nodes) == 0 {
		return fmt.Errorf("tree is empty")
	}
	var printNode func(int, string)
	printNode = func(index int, prefix string) {
		if index == -1 || index >= len(t.nodes) {
			return
		}
		node := t.nodes[index]
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
	printNode(len(t.nodes)-1, "") // Start printing from the root
	return nil
}

func (t *tree) encodedMessage(input string) string {
	var binaryString string
	encodingMap := make(map[rune]string)
	for _, n := range t.nodes {
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

	freqMap, err := getSymbolFreq(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	var pq queue
	err = pq.createQueue(freqMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	var huffTree tree
	err = huffTree.makeTree(&pq)
	if err != nil {
		fmt.Println(err)
		return
	}

	huffTree.assignBinaryPrefixes(len(huffTree.nodes)-1, "")

	err = huffTree.printTree()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(huffTree.encodedMessage(input))
}
