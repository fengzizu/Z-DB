package core

import (
	"encoding/binary"

	"github.com/nikoksr/assert-go"
)

/*
A node consists of:
1. A fixed-sized header containing the type of the node (leaf node or internal node) and the number of keys.
2. A list of pointers to the child nodes. (Used by internal nodes).
3. A list of offsets pointing to each key-value pair.
4. Packed KV pairs.

| type | nkeys |  pointers  |   offsets  | key-values
|  2B  |  2B   | nkeys * 8B | nkeys * 2B | ...

This is the format of the KV pair. Lengths followed by data.
| klen | vlen | key | val |
|  2B  |  2B  | ... | ... |
*/

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

// header
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

// pointers
func (node BNode) getPtr(index uint16) uint64 {
	assert.Assert(index < node.nkeys(), "index is out of bounds", "index", index, "nkeys", node.nkeys())
	pos := HEADER + 8*index
	return binary.LittleEndian.Uint64(node.data[pos:])
}

func (node BNode) setPtr(index uint16, val uint64) {
	assert.Assert(index < node.nkeys(), "index is out of bounds", "index", index, "nkeys", node.nkeys())
	pos := HEADER + 8*index
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

// offset list
func offsetPos(node BNode, index uint16) uint16 {
	assert.Assert(1 <= index && index <= node.nkeys(), "index is out of bounds", "index", index, "nkeys", node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(index-1)
}

func (node BNode) getOffset(index uint16) uint16 {
	if index == 0 {
		return 0
	}

	return binary.LittleEndian.Uint16(node.data[offsetPos(node, index):])
}

func (node BNode) setOffset(index uint16, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, index):], offset)
}

// key values
func (node BNode) kvPos(index uint16) uint16 {
	assert.Assert(index <= node.nkeys(), "index is out of bounds", "index", index, "nkeys", node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(index)
}

func (node BNode) getKey(index uint16) []byte {
	assert.Assert(index < node.nkeys(), "index is out of bounds", "index", index, "nkeys", node.nkeys())
	pos := node.kvPos(index)
	klen := binary.LittleEndian.Uint16(node.data[pos:])
	return node.data[pos+4:][:klen]
}

func (node BNode) getVal(index uint16) []byte {
	assert.Assert(index < node.nkeys(), "index is out of bounds", "index", index, "nkeys", node.nkeys())
	pos := node.kvPos(index)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	return node.data[pos+4+klen:][:vlen]
}

// node size in bytes
/*
图示（nkeys=3的叶节点）
字节位置:  0    4         28        34   36   38       44       54       64
           │    │          │         │    │    │        │        │        │
           ▼    ▼          ▼         ▼    ▼    ▼        ▼        ▼        ▼
┌──────────┬────┬─────────┬─────────┬────┬────┬────────┬────────┬────────┐
│ 类型=2   │n=3 │ [0][0][0│[0] [10] │klen│vlen│ "a"    │klen    │...     │
│ nkeys=3  │    │ ] (24B) │ [20]    │=1  │=5  │"val_a" │=1...   │        │
│ (4B)     │    │         │ (6B)    │    │    │(10B)   │        │        │
└──────────┴────┴─────────┴─────────┴────┴────┴────────┴────────┴────────┘
           │              │         │
           │              │         └─ KV数据区起始于34（4+24+6）
           │              └─ 偏移量列表起始于28 (4+24)
           └─ 头部结束于4

方法调用:
- kvPos(0) = 4 + 24 + 6 + 0 = 34
- kvPos(1) = 4 + 24 + 6 + 10 = 44
- kvPos(2) = 4 + 24 + 6 + 20 = 54
- kvPos(3) = 4 + 24 + 6 + 30 = 64 ← nbytes()用的就是这个

- getKey(0): pos=34, klen=1, return data[38:39] = "a"
- getVal(0): pos=34, klen=1, vlen=5, return data[39:44] = "val_a"
*/
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}
