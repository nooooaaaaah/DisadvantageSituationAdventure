package huffman

import (
	"reflect"
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
				t.Errorf("CreateQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(q, tt.want) {
				t.Errorf("CreateQueue() got = %v, want %v", q, tt.want)
			}
		})
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
			err := tr.MakeTree(&tt.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Optionally check the structure of the tree
			if !tt.wantErr && len(tr.Nodes) < 2 {
				t.Errorf("MakeTree() resulted in an incorrect tree size = %v, want at least 2", len(tr.Nodes))
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
		t.Errorf("AssignBinaryPrefixes() failed to assign correct prefixes, got = %v, %v", tr.Nodes[1].binPrefix, tr.Nodes[2].binPrefix)
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
	expected := "Encoded message as binary value: 01"
	if result := tr.EncodedMessage("ab"); result != expected {
		t.Errorf("EncodedMessage() = %v, want %v", result, expected)
	}
}

func TestGetSymbolFreq(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    map[Value]SymbolFreq
		wantErr bool
	}{
		{
			name:  "string input",
			input: "abba",
			want: map[Value]SymbolFreq{
				'a': 2,
				'b': 2,
			},
			wantErr: false,
		},
		{
			name:  "byte slice input",
			input: []byte("abba"),
			want: map[Value]SymbolFreq{
				byte('a'): 2,
				byte('b'): 2,
			},
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty byte slice",
			input:   []byte{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   123,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSymbolFreq(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSymbolFreq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSymbolFreq() = %v, want %v", got, tt.want)
			}
		})
	}
}
