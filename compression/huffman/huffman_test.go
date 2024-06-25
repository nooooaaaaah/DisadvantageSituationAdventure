package huffman

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestCreateQueue(t *testing.T) {
	tests := []struct {
		name    string
		nodes   map[Value]SymbolFreq
		want    Queue
		wantErr bool
	}{
		{
			name: "single node",
			nodes: map[Value]SymbolFreq{
				'a': 1,
			},
			want: Queue{
				{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
			},
			wantErr: false,
		},
		{
			name:    "empty nodes",
			nodes:   map[Value]SymbolFreq{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "multiple nodes",
			nodes: map[Value]SymbolFreq{
				'a': 3,
				'b': 2,
				'c': 1,
			},
			want: Queue{
				{leaf: true, freq: 1, value: 'c', left: -1, right: -1},
				{leaf: true, freq: 2, value: 'b', left: -1, right: -1},
				{leaf: true, freq: 3, value: 'a', left: -1, right: -1},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q Queue
			err := q.CreateQueue(tt.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("createQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(q, tt.want) {
				t.Errorf("createQueue() got = %v, want %v", q, tt.want)
			}
		})
	}
}

func TestGetSymbolFreq(t *testing.T) {
	tests := []struct {
		data    string
		want    map[Value]SymbolFreq
		wantErr bool
	}{
		{"aaaabbc", map[Value]SymbolFreq{'a': 4, 'b': 2, 'c': 1}, false},
		{"", nil, true},
		{"cccccc", map[Value]SymbolFreq{'c': 6}, false},
		{"ab", map[Value]SymbolFreq{'a': 1, 'b': 1}, false},
	}

	for _, tt := range tests {
		got, err := GetSymbolFreq(tt.data)
		if (err != nil) != tt.wantErr {
			t.Errorf("getSymbolFreq(%q) error = %v, wantErr %v", tt.data, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
			t.Errorf("getSymbolFreq(%q) = %v, want %v", tt.data, got, tt.want)
		}
	}
}

func TestMakeTree(t *testing.T) {
	tests := []struct {
		name    string
		nodes   Queue
		wantErr bool
	}{
		{
			name: "simple tree",
			nodes: Queue{
				{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
				{leaf: true, freq: 2, value: 'b', left: -1, right: -1},
			},
			wantErr: false,
		},
		{
			name: "insufficient nodes",
			nodes: Queue{
				{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tr Tree
			err := tr.makeTree(&tt.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeTree() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Optionally check the structure of the tree
			if !tt.wantErr && len(tr.Nodes) < 2 {
				t.Errorf("makeTree() resulted in an incorrect tree size = %v, want at least 2", len(tr.Nodes))
			}
		})
	}
}

func TestAssignBinaryPrefixes(t *testing.T) {
	tr := Tree{
		Nodes: []Node{
			{leaf: false, freq: 3, left: 1, right: 2},
			{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
			{leaf: true, freq: 2, value: 'b', left: -1, right: -1},
		},
	}
	tr.AssignBinaryPrefixes(0, "")
	if tr.Nodes[1].binPrefix != "0" || tr.Nodes[2].binPrefix != "1" {
		t.Errorf("assignBinaryPrefixes() failed to assign correct prefixes, got = %v, %v", tr.Nodes[1].binPrefix, tr.Nodes[2].binPrefix)
	}
}

func TestEncodedMessage(t *testing.T) {
	tr := Tree{
		Nodes: []Node{
			{leaf: false, freq: 3, binPrefix: "", left: 1, right: 2},
			{leaf: true, freq: 1, value: 'a', binPrefix: "0", left: -1, right: -1},
			{leaf: true, freq: 2, value: 'b', binPrefix: "1", left: -1, right: -1},
		},
	}
	var output strings.Builder
	fmt.Fprintf(&output, "Encoded message as binary value: %s", "01")
	expected := output.String()

	if result := tr.EncodedMessage("ab"); result != expected {
		t.Errorf("encodedMessage() = %v, want %v", result, expected)
	}
}
