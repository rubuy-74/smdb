package b3

import (
	"bytes"
	"fmt"
)

type BTree struct {
	root uint64

	get func(uint64) []byte // reads page
	new func([]byte) uint64 // inserts page
	del func(uint64)        // deletes page
}

func TreeInsert(tree *BTree, node BNode, key []byte, val []byte) (BNode, error) {
	new := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	idx, err := NodeLookupLE(node, key)
	if err != nil {
		return nil, err
	}
	switch node.Btype() {
	case BNODE_LEAF:
		nKey, err := node.GetKey(idx)
		if err != nil {
			return nil, err
		}
		if bytes.Equal(key, nKey) {
			LeafUpdate(new, node, idx, key, val)
		} else {
			LeafInsert(new, node, idx, key, val)
		}
	case BNODE_NODE:
		//
		kptr, err := node.GetPtr(idx)
		if err != nil {
			return nil, err
		}
		knode, err := TreeInsert(tree, tree.get(kptr), key, val)
		if err != nil {
			return nil, err
		}
		nsplit, split, err := NodeSplit3(knode)
		if err != nil {
			return nil, err
		}
		tree.del(kptr)
		NodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
	}
	return new, nil
}

func NodeReplaceKidN(
	tree *BTree, new BNode, old BNode, idx uint16,
	kids ...BNode,
) error {
	inc := uint16(len(kids))
	new.SetHeader(BNODE_NODE, old.Nkeys()+inc-1)
	NodeAppendRange(new, old, 0, 0, idx) // add previous KV pairs
	for i, node := range kids {
		nKey, err := node.GetKey(0)
		if err != nil {
			return err
		}
		NodeAppendKV(new, idx+uint16(i), tree.new(node), nKey, nil) // Add references to children
	}
	NodeAppendRange(new, old, idx+inc, idx+1, old.Nkeys()-(idx+1))
	return nil
}

func (tree *BTree) Insert(key []byte, val []byte) error {
	checkLimit := func() bool {
		return len(key) < BTREE_MAX_KEY_SIZE && len(val) < BTREE_MAX_VAL_SIZE
	}
	if !checkLimit() {
		return fmt.Errorf("failed because of key or val size constraint on Insert")
	}

	if tree.root == 0 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.SetHeader(BNODE_LEAF, 2)
		NodeAppendKV(root, 0, 0, nil, nil)
		NodeAppendKV(root, 1, 0, key, val)
		tree.root = tree.new(root)
		return nil
	}

	node, err := TreeInsert(tree, tree.get(tree.root), key, val)
	if err != nil {
		return err
	}
	nsplit, split, err := NodeSplit3(node)
	if err != nil {
		return err
	}
	if nsplit > 1 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.SetHeader(BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			key, err := knode.GetKey(0)
			if err != nil {
				return err
			}
			ptr := tree.new(knode)
			NodeAppendKV(root, uint16(i), ptr, key, nil)
		}

		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(split[0])
	}

	return nil
}

func (tree *BTree) Delete(key []byte) (bool, error) {
	if tree.root == 0 {
		return false, nil
	}
	node, err := TreeDelete(tree, tree.get(tree.root), key)
	if err != nil {
		return false, err
	}
	tree.root = tree.new(node)
	return true, nil
}

func NodeMerge(new BNode, left BNode, right BNode) error {
	if left.Btype() == BNODE_LEAF && right.Btype() == BNODE_NODE {
		return fmt.Errorf("failed to merge nodes [Leaf-Node]")
	}
	if left.Btype() == BNODE_NODE && right.Btype() == BNODE_LEAF {
		return fmt.Errorf("failed to merge nodes [Node-Leaf]")
	}
	if left.Btype() == BNODE_LEAF {
		new.SetHeader(BNODE_LEAF, left.Nkeys()+right.Nkeys())
	} else {
		new.SetHeader(BNODE_NODE, left.Nkeys()+right.Nkeys())
	}
	NodeAppendRange(new, left, 0, 0, left.Nkeys())
	NodeAppendRange(new, right, 0, 0, left.Nkeys())
	return nil
}

func NodeReplace2Kid(new BNode, old BNode, idx uint16, ptr uint64, key []byte) error {
	if new.Btype() == BNODE_LEAF {
		return fmt.Errorf("function NodeReplace2Kid called on a leaf")
	}
	new.SetHeader(BNODE_NODE, old.Nkeys()-1)
	NodeAppendRange(new, old, 0, 0, idx)
	NodeAppendKV(new, idx, ptr, key, nil)
	NodeAppendRange(new, old, idx+2, idx+2, old.Nkeys()-(idx+2))
	return nil
}

func ShouldMerge(
	tree *BTree, node BNode, idx uint16, updated BNode,
) (int, BNode, error) {
	if updated.Nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}, nil
	}
	if idx > 0 {
		nPtr, err := node.GetPtr(idx - 1)
		if err != nil {
			return 0, BNode{}, err
		}
		sibling := BNode(tree.get(nPtr))
		merged_size := sibling.Nbytes() + (updated.Nbytes() - 4)
		if merged_size <= BTREE_PAGE_SIZE {
			return -1, sibling, nil
		}
	}
	if idx+1 < node.Nkeys() {
		nPtr, err := node.GetPtr(idx + 1)
		if err != nil {
			return 0, BNode{}, err
		}
		sibling := BNode(tree.get(nPtr))
		merged_size := sibling.Nbytes() + (updated.Nbytes() - 4)
		if merged_size <= BTREE_PAGE_SIZE {
			return +1, sibling, nil
		}
	}
	return 0, BNode{}, nil
}

func TreeDelete(tree *BTree, node BNode, key []byte) (BNode, error) {

	idx, err := NodeLookupLE(node, key)
	new := BNode(make([]byte, BTREE_PAGE_SIZE))
	if err != nil {
		return BNode{}, err
	}
	switch node.Btype() {
	case BNODE_LEAF:
		LeafDelete(new, node, idx)
	case BNODE_NODE:
		new, err = NodeDelete(tree, node, idx, key)
		if err != nil {
			return BNode{}, err
		}
	}
	return new, nil
}

func NodeDelete(tree *BTree, node BNode, idx uint16, key []byte) (BNode, error) {
	kptr, err := node.GetPtr(idx)
	if err != nil {
		return BNode{}, err
	}
	updated, err := TreeDelete(tree, tree.get(kptr), key)
	if err != nil {
		return BNode{}, err
	}
	tree.del(kptr)

	new := BNode(make([]byte, BTREE_PAGE_SIZE))
	mergeDir, sibling, err := ShouldMerge(tree, node, idx, updated)
	if err != nil {
		return BNode{}, err
	}
	switch {
	case mergeDir < 0:
		merged := BNode(make([]byte, BTREE_PAGE_SIZE))
		NodeMerge(merged, sibling, updated)
		nPtr, err := node.GetPtr(idx - 1)
		if err != nil {
			return BNode{}, err
		}
		tree.del(nPtr)
		mergedKey, err := merged.GetKey(0)
		if err != nil {
			return BNode{}, err
		}
		NodeReplace2Kid(new, node, idx-1, tree.new(merged), mergedKey)
	case mergeDir > 0:
		merged := BNode(make([]byte, BTREE_PAGE_SIZE))
		NodeMerge(merged, sibling, updated)
		nPtr, err := node.GetPtr(idx + 1)
		if err != nil {
			return BNode{}, err
		}
		tree.del(nPtr)
		mergedKey, err := merged.GetKey(0)
		if err != nil {
			return BNode{}, err
		}
		NodeReplace2Kid(new, node, idx, tree.new(merged), mergedKey)
	case mergeDir == 0 && updated.Nkeys() == 0:
		if node.Nkeys() != 1 || idx != 0 {
			return BNode{}, fmt.Errorf("empty child with no sibling")
		}
		new.SetHeader(BNODE_NODE, 0)
	case mergeDir == 0 && updated.Nkeys() > 0:
		NodeReplaceKidN(tree, new, node, idx, updated)
	}
	return new, nil
}
