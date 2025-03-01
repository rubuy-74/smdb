package b3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/rubuy-74/smDB/internal/utils"
)

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)
const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

const TYPE_SIZE = 2
const NKEYS_SIZE = 2
const POINTER_SIZE = 8
const OFFSET_SIZE = 2

const KEY_SIZE_SIZE = 2
const VALUE_SIZE_SIZE = 2

func init() {
	node1max := 4 + 1*8 + 1*2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VAL_SIZE
	//assert(node1max <= BTREE_PAGE_SIZE)
	if node1max > BTREE_PAGE_SIZE {
		log.Fatal("Node exceeded max size")
	}
}

// file pointer is not an address in memory but it is a file offset, the page number

// Data needs not to be serialized, as it can be acessed directely
// Node structure
// type | nkeys | pointers | offsets | key-values
//
//	2B      2B      8B*p       2B*o        ?
//
// Key Value Structure
// key_size | val_size | key | value
//
//	2B				 2B				?			?
type BNode []byte // can be dumped into disk

func (node BNode) PrintKV(idx uint16) error {
	key, err := node.GetKey(idx)
	if err != nil {
		return err
	}
	val, err := node.GetVal(idx)
	if err != nil {
		return err
	}

	fmt.Printf("{%s : %s}\n", utils.ToString(key), utils.ToString(val))
	return nil
}

func (node BNode) PrintAllKV() error {
	for i := uint16(0); i < uint16(node.Nkeys()); i++ {
		err := node.PrintKV(i)
		if err != nil {
			return err
		}
	}
	return nil
}

// reads btype
func (node BNode) Btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

// reads nkeys
func (node BNode) Nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

// writes the fixed-size header (btype+nkeys)
func (node BNode) SetHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype) // overrides btype
	binary.LittleEndian.PutUint16(node[2:4], nkeys) // overrides nkeys
}

// reads pointer array
func (node BNode) GetPtr(idx uint16) (uint64, error) {
	if idx >= node.Nkeys() {
		return 0, fmt.Errorf("invalid index %d in getPtr", idx)
	}
	pos := 4 + 8*idx
	return binary.LittleEndian.Uint64(node[pos:]), nil
}

// updates pointer array
func (node BNode) SetPtr(idx uint16, val uint64) error {
	if idx >= node.Nkeys() {
		return fmt.Errorf("invalid index %d in setPtr", idx)
	}
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
	return nil
}

// reads offset to locate nth key in O(1)
func (node BNode) GetOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*node.Nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) SetOffset(idx uint16, val uint16) {
	if idx == 0 {
		return
	}
	pos := 4 + 8*node.Nkeys() + 2*(idx-1)
	binary.LittleEndian.PutUint16(node[pos:], val)
}

// first KV pair is at node[node.kvPos(0):]
//
// returns the position of nth key in the Node
func (node BNode) KVPos(idx uint16) (uint16, error) {
	if idx >= node.Nkeys() {
		return 0, fmt.Errorf("invalid index %d in kvPos", idx)
	}
	return 4 + 8*node.Nkeys() + 8*node.Nkeys() + node.GetOffset(idx), nil
}

// returns nth key of the Node
func (node BNode) GetKey(idx uint16) ([]byte, error) {
	if idx >= node.Nkeys() {
		return []byte{}, fmt.Errorf("invalid index %d in getKey", idx)
	}
	pos, err := node.KVPos(idx)
	if err != nil {
		return []byte{}, err
	}
	klen := binary.LittleEndian.Uint16((node[pos:]))
	return node[pos+4:][:klen], nil
}

// returns nth value of the Node (leaf nodes only)
func (node BNode) GetVal(idx uint16) ([]byte, error) {
	if idx >= node.Nkeys() {
		return []byte{}, fmt.Errorf("invalid index %d in getKey", idx)
	}
	pos, err := node.KVPos(idx)
	if err != nil {
		return []byte{}, err
	}

	klen := binary.LittleEndian.Uint16(node[pos+0:])
	vlen := binary.LittleEndian.Uint16(node[pos+KEY_SIZE_SIZE:])

	return node[4+pos+klen:][:vlen], nil
}

func (node BNode) Nbytes() uint16 {
	bytes, _ := node.KVPos(node.Nkeys())
	return bytes
}

// appends key value pair to Node
//
// assumes that keys are in order
func NodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) error {
	new.SetPtr(idx, ptr)
	pos, err := new.KVPos(idx)
	if err != nil {
		return err
	}

	klen := uint16(len(key))
	vlen := uint16(len(val))

	binary.LittleEndian.PutUint16(new[pos+0:], klen)
	binary.LittleEndian.PutUint16(new[pos+KEY_SIZE_SIZE:], vlen)

	/* 	fmt.Printf("[NodeAppendKV] - key=%d\n", key)
	   	fmt.Printf("[NodeAppendKV] - pos=%d\n", pos+(KEY_SIZE_SIZE+VALUE_SIZE_SIZE)) */

	copy(new[pos+4:], key)
	copy(new[pos+4+klen:], val)

	new.SetOffset(idx+1, new.GetOffset(idx)+4+uint16(klen+vlen))
	return nil
}

func NodeAppendRange(
	new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16,
) error {
	for i := uint16(0); i < n; i++ {
		dst, src := dstNew+i, srcOld+i
		oldPtr, err := old.GetPtr(src)
		if err != nil {
			return err
		}
		oldKey, err := old.GetKey(src)
		if err != nil {
			return err
		}
		oldVal, err := old.GetVal(src)
		if err != nil {
			return err
		}
		NodeAppendKV(
			new, dst, oldPtr, oldKey, oldVal,
		)
	}
	return nil
}

// inserts key-value pair to Node (maintaining the remaining)
func LeafInsert(
	new BNode, old BNode, idx uint16, key []byte, val []byte,
) {
	new.SetHeader(BNODE_LEAF, old.Nkeys()+1)
	NodeAppendRange(new, old, 0, 0, idx)
	NodeAppendKV(new, idx, 0, key, val)
	NodeAppendRange(new, old, idx+1, idx, old.Nkeys()-idx)
}

func LeafUpdate(
	new BNode, old BNode, idx uint16, key []byte, val []byte,
) {
	new.SetHeader(BNODE_LEAF, old.Nkeys())
	NodeAppendRange(new, old, 0, 0, idx)
	NodeAppendKV(new, idx, 0, key, val)
	NodeAppendRange(new, old, idx+1, idx+1, old.Nkeys()-(idx+1))
}

// TODO: Use binary search
func NodeLookupLE(node BNode, key []byte) (uint16, error) {
	nkeys := node.Nkeys()
	var i uint16
	for i = range nkeys {
		ikey, err := node.GetKey(i)
		if err != nil {
			return 0, err
		}
		cmp := bytes.Compare(ikey, key)
		if cmp == 0 {
			return i, nil
		}
		if cmp > 0 {
			return i - 1, nil
		}
	}
	return i - 1, nil
}

func ChangeKVPair(old BNode, key []byte, val []byte) (BNode, error) {
	new := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	idx, err := NodeLookupLE(old, key)
	if err != nil {
		return nil, err
	}
	nKey, err := old.GetKey(idx)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(key, nKey) {
		LeafUpdate(new, old, idx, key, val)
	} else {
		LeafInsert(new, old, idx, key, val)
	}
	return new, nil
}

func NodeSplit2(left BNode, right BNode, old BNode) error {
	if old.Nkeys() < 2 {
		return fmt.Errorf("node has less than 2 node in NodeSplit2")
	}

	nleft := old.Nkeys() / 2
	left_bytes := func(nleft uint16) uint16 {
		return 4 + 8*nleft + 2*nleft + old.GetOffset(nleft)
	}
	for left_bytes(nleft) > BTREE_PAGE_SIZE {
		nleft--
	}
	if nleft < 1 {
		return fmt.Errorf("nleft is too small (<1) in NodeSplit2")
	}

	right_bytes := func(nleft uint16) uint16 {
		return old.Nbytes() - (left_bytes(nleft) - 4)
	}
	for right_bytes(nleft) > BTREE_PAGE_SIZE {
		nleft++
	}
	if nleft >= old.Nkeys() {
		return fmt.Errorf("nleft is too bit (>= old.Nkeys()) in NodeSplit2")
	}

	// newNodes
	left.SetHeader(old.Btype(), nleft)
	right.SetHeader(old.Btype(), old.Nkeys()-nleft)
	NodeAppendRange(left, old, 0, 0, nleft)
	NodeAppendRange(right, old, 0, nleft, old.Nkeys()-nleft)

	if right.Nbytes() > BTREE_PAGE_SIZE {
		return fmt.Errorf("right node is too big in NodeSplit2")
	}
	return nil
}

func NodeSplit3(old BNode) (uint16, [3]BNode, error) {
	if old.Nbytes() <= BTREE_PAGE_SIZE {
		old = old[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}, nil
	}
	left := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	right := BNode(make([]byte, BTREE_PAGE_SIZE))
	err := NodeSplit2(left, right, old)
	if err != nil {
		return 0, [3]BNode{}, err
	}
	if left.Nbytes() <= BTREE_PAGE_SIZE {
		left = left[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}, nil
	}
	leftleft := BNode(make([]byte, BTREE_PAGE_SIZE))
	middle := BNode(make([]byte, BTREE_PAGE_SIZE))
	NodeSplit2(leftleft, middle, left)
	if leftleft.Nbytes() > BTREE_PAGE_SIZE {
		return 0, [3]BNode{}, fmt.Errorf("most left node is too big in NodeSplit3")
	}
	return 3, [3]BNode{leftleft, middle, right}, nil
}
