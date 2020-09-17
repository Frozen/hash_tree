package hash_tree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Block struct {
	parentHash uint64
	hash       uint64
}

type Path = []Block

// The input blocks are the blocks in a block tree, the hash of
// the root block is 0. This function finds the longest chain and
// returns it.
func FindLongestChain(blocks []Block) []Block {
	tree := make(map[uint64][]Block)
	for _, block := range blocks {
		tree[block.parentHash] = append(tree[block.parentHash], block)
	}
	return walk(tree, nil, 0)
}

// Recursive walk through tree.
func walk(nodes map[uint64][]Block, path Path, node uint64) Path {
	longestPath := path

	for _, n := range nodes[node] {
		cur := walk(nodes, append(path, n), n.hash)
		if len(cur) > len(longestPath) {
			longestPath = cur
		}
	}
	return longestPath
}

func TestFindLongestChain(t *testing.T) {
	blocks := []Block{
		{0, 1},
		{0, 2},
		{1, 3},
		{2, 4},
		{2, 5},
		{3, 6},
	}
	rs := FindLongestChain(blocks)
	t.Log(rs)
	require.Equal(t, Path{
		Block{0, 1},
		Block{1, 3},
		Block{3, 6},
	}, rs)
}

func FindLongestPathWithoutRecursion(blocks []Block) Path {
	type node = uint64
	type parent = Block
	tree := make(map[node]parent)
	for _, block := range blocks {
		tree[block.hash] = block
	}
	longestPath := Path{}
	for _, block := range blocks {
		cur := block
		path := Path{cur}
		for {
			parent, ok := tree[cur.parentHash]
			if !ok {
				break
			}
			path = append(path, parent)
			cur = parent
			if parent.hash == 0 {
				break
			}
		}
		if len(path) > len(longestPath) {
			longestPath = path
		}
	}
	return revert(longestPath)
}

// Revert path.
func revert(path Path) Path {
	l := len(path) - 1
	for i := 0; i < len(path)/2; i++ {
		tmp := path[i]
		path[i] = path[l-i]
		path[l-i] = tmp
	}
	return path
}

func TestFindLongestPathWithoutRecursion(t *testing.T) {
	blocks := []Block{
		{0, 1},
		{0, 2},
		{1, 3},
		{2, 4},
		{2, 5},
		{3, 6},
	}
	rs := FindLongestPathWithoutRecursion(blocks)
	t.Log(rs)
	require.Equal(t, Path{
		Block{0, 1},
		Block{1, 3},
		Block{3, 6},
	}, rs)
}
