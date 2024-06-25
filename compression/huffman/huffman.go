package huffman

import (
	"bufio"
	"fmt"
	"os"
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
	Value      any
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
		if sortedNodes[i].freq == sortedNodes[j].freq {
			return fmt.Sprintf("%v", sortedNodes[i].value) < fmt.Sprintf("%v", sortedNodes[j].value)
		}
		return sortedNodes[i].freq < sortedNodes[j].freq
	})
	*q = sortedNodes
	return nil
}

func GetSymbolFreq(data any) (map[Value]SymbolFreq, error) {
	freqMap := make(map[Value]SymbolFreq)
	switch v := data.(type) {
	case string:
		if v == "" {
			return nil, fmt.Errorf("data was empty")
		}
		for _, char := range v {
			freqMap[char]++
		}
	case []byte:
		if len(v) == 0 {
			return nil, fmt.Errorf("data was empty")
		}
		for _, b := range v {
			freqMap[b]++
		}
	default:
		return nil, fmt.Errorf("unsupported data type")
	}
	return freqMap, nil
}

func (t *Tree) MakeTree(pq *Queue) error {
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
	if node.leaf {
		var valueStr string
		switch v := node.value.(type) {
		case byte:
			valueStr = fmt.Sprintf("'%c'", v)
		case rune:
			valueStr = fmt.Sprintf("'%c'", v)
		default:
			valueStr = fmt.Sprintf("%v", v)
		}
		fmt.Printf("Assigning prefix %s to node with freq %d and value of %s\n", binPrefix, node.freq, valueStr)
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
	fmt.Print("\nHere's a visual representation of the tree\n\n")
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
			var valueStr string
			switch v := node.value.(type) {
			case byte:
				valueStr = fmt.Sprintf("'%c'", v)
			case rune:
				valueStr = fmt.Sprintf("'%c'", v)
			default:
				valueStr = fmt.Sprintf("%v", v)
			}
			fmt.Printf("%s└── %s (%s, %d)\n", prefix, node.binPrefix, valueStr, node.freq)
		} else {
			fmt.Printf("%s├── %s node (%d)\n", prefix, node.binPrefix, node.freq)
			newPrefix := prefix + "│   "
			printNode(node.left, newPrefix)
			if node.right != -1 {
				printNode(node.right, prefix+"    ")
			}
		}
	}
	printNode(len(t.Nodes)-1, "")
	return nil
}

func (t *Tree) EncodedMessage(input any) string {
	var binaryString string
	encodingMap := make(map[Value]string)
	for _, n := range t.Nodes {
		if n.leaf {
			encodingMap[n.value] = n.binPrefix
		}
	}
	switch v := input.(type) {
	case []byte:
		for _, b := range v {
			binaryString += encodingMap[b]
		}
	case string:
		for _, char := range v {
			binaryString += encodingMap[char]
		}
	default:
		return "Unsupported input type"
	}
	return fmt.Sprintf("Encoded message as binary value: %s", binaryString)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the text to encode:")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	input = strings.TrimSpace(input)

	// String encoding
	fmt.Println("Encoding as string:")
	encodeData(input)

	// Byte slice encoding
	fmt.Println("\nNow encoding the same input as a byte slice:")
	encodeData([]byte(input))
}

func encodeData(data any) {
	freqMap, err := GetSymbolFreq(data)
	if err != nil {
		fmt.Println("Error getting symbol frequencies:", err)
		return
	}

	var pq Queue
	if err := pq.CreateQueue(freqMap); err != nil {
		fmt.Println("Error creating queue:", err)
		return
	}

	var huffTree Tree
	if err := huffTree.MakeTree(&pq); err != nil {
		fmt.Println("Error making tree:", err)
		return
	}

	huffTree.AssignBinaryPrefixes(len(huffTree.Nodes)-1, "")

	if err := huffTree.PrintTree(); err != nil {
		fmt.Println("Error printing tree:", err)
		return
	}

	fmt.Println(huffTree.EncodedMessage(data))
}
