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
		nodes   map[value]symbolFreq
		want    queue
		wantErr bool
	}{
		{
			name: "single node",
			nodes: map[value]symbolFreq{
				'a': 1,
			},
			want: queue{
				{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
			},
			wantErr: false,
		},
		{
			name:    "empty nodes",
			nodes:   map[value]symbolFreq{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "multiple nodes",
			nodes: map[value]symbolFreq{
				'a': 3,
				'b': 2,
				'c': 1,
			},
			want: queue{
				{leaf: true, freq: 1, value: 'c', left: -1, right: -1},
				{leaf: true, freq: 2, value: 'b', left: -1, right: -1},
				{leaf: true, freq: 3, value: 'a', left: -1, right: -1},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var q queue
			err := q.createQueue(tt.nodes)
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
		want    map[value]symbolFreq
		wantErr bool
	}{
		{"aaaabbc", map[value]symbolFreq{'a': 4, 'b': 2, 'c': 1}, false},
		{"", nil, true},
		{"cccccc", map[value]symbolFreq{'c': 6}, false},
		{"ab", map[value]symbolFreq{'a': 1, 'b': 1}, false},
	}

	for _, tt := range tests {
		got, err := getSymbolFreq(tt.data)
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
		nodes   queue
		wantErr bool
	}{
		{
			name: "simple tree",
			nodes: queue{
				{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
				{leaf: true, freq: 2, value: 'b', left: -1, right: -1},
			},
			wantErr: false,
		},
		{
			name: "insufficient nodes",
			nodes: queue{
				{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tr tree
			err := tr.makeTree(&tt.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeTree() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Optionally check the structure of the tree
			if !tt.wantErr && len(tr.nodes) < 2 {
				t.Errorf("makeTree() resulted in an incorrect tree size = %v, want at least 2", len(tr.nodes))
			}
		})
	}
}

func TestAssignBinaryPrefixes(t *testing.T) {
	tr := tree{
		nodes: []node{
			{leaf: false, freq: 3, left: 1, right: 2},
			{leaf: true, freq: 1, value: 'a', left: -1, right: -1},
			{leaf: true, freq: 2, value: 'b', left: -1, right: -1},
		},
	}
	tr.assignBinaryPrefixes(0, "")
	if tr.nodes[1].binPrefix != "0" || tr.nodes[2].binPrefix != "1" {
		t.Errorf("assignBinaryPrefixes() failed to assign correct prefixes, got = %v, %v", tr.nodes[1].binPrefix, tr.nodes[2].binPrefix)
	}
}

func TestEncodedMessage(t *testing.T) {
	tr := tree{
		nodes: []node{
			{leaf: false, freq: 3, binPrefix: "", left: 1, right: 2},
			{leaf: true, freq: 1, value: 'a', binPrefix: "0", left: -1, right: -1},
			{leaf: true, freq: 2, value: 'b', binPrefix: "1", left: -1, right: -1},
		},
	}
	var output strings.Builder
	fmt.Fprintf(&output, "Encoded message as binary value: %s", "01")
	expected := output.String()

	if result := tr.encodedMessage("ab"); result != expected {
		t.Errorf("encodedMessage() = %v, want %v", result, expected)
	}
}
