package core

import (
	"github.com/nikoksr/assert-go"
)

type BNode struct {
	data []byte // can be dumped to the disk
}

const (
	BNODE_NODE = 1 // internal nodes without values
	BNODE_LEAF = 2 // leaf nodes with values
)

type BTree struct {
	// pointer (a nonzero page number)
	root uint64

	// callbacks for managing on-disk pages
	get func(uint64) BNode // dereference a pointer
	new func(BNode) uint64 // allacate a new page
	del func(uint64)       // deallocate a page
}

const (
	HEADER             = 4
	BTREE_PAGE_SIZE    = 4096
	BTREE_MAX_KEY_SIZE = 1000
	BTREE_MAX_VAL_SIZE = 3000
)

func init() {
	node1max := HEADER + 8 + 2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	assert.Assert(node1max <= BTREE_PAGE_SIZE, "node1max is too large", "node1max", node1max)
}
