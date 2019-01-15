// Code generated by gotemplate. DO NOT EDIT.

package forkdb_multi_index

import (
	"fmt"

	"github.com/eosspark/eos-go/common/container/multiindex"
	"github.com/eosspark/eos-go/log"
)

// template type OrderedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Comparator,Multiply)

// OrderedIndex holds elements of the red-black tree
type ByLibBlockNum struct {
	super *MultiIndexBase // index on the OrderedIndex, IndexBase is the last super index
	final *MultiIndex     // index under the OrderedIndex, MultiIndex is the final index

	Root *ByLibBlockNumNode
	size int
}

func (tree *ByLibBlockNum) init(final *MultiIndex) {
	tree.final = final
	tree.super = &MultiIndexBase{}
	tree.super.init(final)
}

func (tree *ByLibBlockNum) clear() {
	tree.Clear()
	tree.super.clear()
}

/*generic class*/

/*generic class*/

// OrderedIndexNode is a single element within the tree
type ByLibBlockNumNode struct {
	Key    ByLibBlockNumComposite
	super  *MultiIndexBaseNode
	final  *MultiIndexNode
	color  colorByLibBlockNum
	Left   *ByLibBlockNumNode
	Right  *ByLibBlockNumNode
	Parent *ByLibBlockNumNode
}

/*generic class*/

/*generic class*/

func (node *ByLibBlockNumNode) value() *BlockStatePtr {
	return node.super.value()
}

type colorByLibBlockNum bool

const (
	blackByLibBlockNum, redByLibBlockNum colorByLibBlockNum = true, false
)

func (tree *ByLibBlockNum) Insert(v BlockStatePtr) (IteratorByLibBlockNum, bool) {
	fn, res := tree.final.insert(v)
	if res {
		return tree.makeIterator(fn), true
	}
	return tree.End(), false
}

func (tree *ByLibBlockNum) insert(v BlockStatePtr, fn *MultiIndexNode) (*ByLibBlockNumNode, bool) {
	key := ByLibBlockNumFunc(v)

	node, res := tree.put(key)
	if !res {
		log.Warn("#ordered index insert failed")
		return nil, false
	}
	sn, res := tree.super.insert(v, fn)
	if res {
		node.super = sn
		node.final = fn
		return node, true
	}
	tree.remove(node)
	return nil, false
}

func (tree *ByLibBlockNum) Erase(iter IteratorByLibBlockNum) (itr IteratorByLibBlockNum) {
	itr = iter
	itr.Next()
	tree.final.erase(iter.node.final)
	return
}

func (tree *ByLibBlockNum) Erases(first, last IteratorByLibBlockNum) {
	for first != last {
		first = tree.Erase(first)
	}
}

func (tree *ByLibBlockNum) erase(n *ByLibBlockNumNode) {
	tree.remove(n)
	tree.super.erase(n.super)
	n.super = nil
	n.final = nil
}

func (tree *ByLibBlockNum) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(IteratorByLibBlockNum); ok {
		tree.Erase(itr)
	} else {
		tree.super.erase_(iter)
	}
}

func (tree *ByLibBlockNum) Modify(iter IteratorByLibBlockNum, mod func(*BlockStatePtr)) bool {
	if _, b := tree.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (tree *ByLibBlockNum) modify(n *ByLibBlockNumNode) (*ByLibBlockNumNode, bool) {
	n.Key = ByLibBlockNumFunc(*n.value())

	if !tree.inPlace(n) {
		tree.remove(n)
		node, res := tree.put(n.Key)
		if !res {
			log.Warn("#ordered index modify failed")
			tree.super.erase(n.super)
			return nil, false
		}

		//n.Left = node.Left
		//if n.Left != nil {
		//	n.Left.Parent = n
		//}
		//n.Right = node.Right
		//if n.Right != nil {
		//	n.Right.Parent = n
		//}
		//n.Parent = node.Parent
		//if n.Parent != nil {
		//	if n.Parent.Left == node {
		//		n.Parent.Left = n
		//	} else {
		//		n.Parent.Right = n
		//	}
		//} else {
		//	tree.Root = n
		//}
		node.super = n.super
		node.final = n.final
		n = node
	}

	if sn, res := tree.super.modify(n.super); !res {
		tree.remove(n)
		return nil, false
	} else {
		n.super = sn
	}

	return n, true
}

func (tree *ByLibBlockNum) modify_(iter multiindex.IteratorType, mod func(*BlockStatePtr)) bool {
	if itr, ok := iter.(IteratorByLibBlockNum); ok {
		return tree.Modify(itr, mod)
	} else {
		return tree.super.modify_(iter, mod)
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByLibBlockNum) Find(key ByLibBlockNumComposite) IteratorByLibBlockNum {
	if true {
		lower := tree.LowerBound(key)
		if !lower.IsEnd() && ByLibBlockNumCompare(key, lower.Key()) == 0 {
			return lower
		}
		return tree.End()
	} else {
		if node := tree.lookup(key); node != nil {
			return IteratorByLibBlockNum{tree, node, betweenByLibBlockNum}
		}
		return tree.End()
	}
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (tree *ByLibBlockNum) LowerBound(key ByLibBlockNumComposite) IteratorByLibBlockNum {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByLibBlockNumCompare(key, node.Key) > 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenByLibBlockNum
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is greater than the given key.
// Complexity: O(log N).
func (tree *ByLibBlockNum) UpperBound(key ByLibBlockNumComposite) IteratorByLibBlockNum {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByLibBlockNumCompare(key, node.Key) >= 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenByLibBlockNum
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByLibBlockNum) Remove(key ByLibBlockNumComposite) {
	if true {
		for lower := tree.LowerBound(key); lower.position != endByLibBlockNum; {
			if ByLibBlockNumCompare(lower.Key(), key) == 0 {
				node := lower.node
				lower.Next()
				tree.remove(node)
			} else {
				break
			}
		}
	} else {
		node := tree.lookup(key)
		tree.remove(node)
	}
}

func (tree *ByLibBlockNum) put(key ByLibBlockNumComposite) (*ByLibBlockNumNode, bool) {
	var insertedNode *ByLibBlockNumNode
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		ByLibBlockNumCompare(key, key)
		tree.Root = &ByLibBlockNumNode{Key: key, color: redByLibBlockNum}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		if true {
			for loop {
				compare := ByLibBlockNumCompare(key, node.Key)
				switch {
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByLibBlockNumNode{Key: key, color: redByLibBlockNum}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare >= 0:
					if node.Right == nil {
						node.Right = &ByLibBlockNumNode{Key: key, color: redByLibBlockNum}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := ByLibBlockNumCompare(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					return node, false
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByLibBlockNumNode{Key: key, color: redByLibBlockNum}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare > 0:
					if node.Right == nil {
						node.Right = &ByLibBlockNumNode{Key: key, color: redByLibBlockNum}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++

	return insertedNode, true
}

func (tree *ByLibBlockNum) swapNode(node *ByLibBlockNumNode, pred *ByLibBlockNumNode) {
	if node == pred {
		return
	}

	tmp := ByLibBlockNumNode{color: pred.color, Left: pred.Left, Right: pred.Right, Parent: pred.Parent}

	pred.color = node.color
	node.color = tmp.color

	pred.Right = node.Right
	if pred.Right != nil {
		pred.Right.Parent = pred
	}
	node.Right = tmp.Right
	if node.Right != nil {
		node.Right.Parent = node
	}

	if pred.Parent == node {
		pred.Left = node
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}
		node.Parent = pred

	} else {
		pred.Left = node.Left
		if pred.Left != nil {
			pred.Left.Parent = pred
		}
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}

		node.Parent = tmp.Parent
		if node.Parent != nil {
			if node.Parent.Left == pred {
				node.Parent.Left = node
			} else {
				node.Parent.Right = node
			}
		} else {
			tree.Root = node
		}
	}
}

func (tree *ByLibBlockNum) remove(node *ByLibBlockNumNode) {
	var child *ByLibBlockNumNode
	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		tree.swapNode(node, pred)
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == blackByLibBlockNum {
			node.color = nodeColorByLibBlockNum(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = blackByLibBlockNum
		}
	}
	tree.size--
}

func (tree *ByLibBlockNum) lookup(key ByLibBlockNumComposite) *ByLibBlockNumNode {
	node := tree.Root
	for node != nil {
		compare := ByLibBlockNumCompare(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

// Empty returns true if tree does not contain any nodes
func (tree *ByLibBlockNum) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *ByLibBlockNum) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *ByLibBlockNum) Keys() []ByLibBlockNumComposite {
	keys := make([]ByLibBlockNumComposite, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *ByLibBlockNum) Values() []BlockStatePtr {
	values := make([]BlockStatePtr, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *ByLibBlockNum) Left() *ByLibBlockNumNode {
	var parent *ByLibBlockNumNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *ByLibBlockNum) Right() *ByLibBlockNumNode {
	var parent *ByLibBlockNumNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Clear removes all nodes from the tree.
func (tree *ByLibBlockNum) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *ByLibBlockNum) String() string {
	str := "OrderedIndex\n"
	if !tree.Empty() {
		outputByLibBlockNum(tree.Root, "", true, &str)
	}
	return str
}

func (node *ByLibBlockNumNode) String() string {
	if !node.color {
		return fmt.Sprintf("(%v,%v)", node.Key, "red")
	}
	return fmt.Sprintf("(%v)", node.Key)
}

func outputByLibBlockNum(node *ByLibBlockNumNode, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		outputByLibBlockNum(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		outputByLibBlockNum(node.Left, newPrefix, true, str)
	}
}

func (node *ByLibBlockNumNode) grandparent() *ByLibBlockNumNode {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *ByLibBlockNumNode) uncle() *ByLibBlockNumNode {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *ByLibBlockNumNode) sibling() *ByLibBlockNumNode {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (node *ByLibBlockNumNode) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *ByLibBlockNum) rotateLeft(node *ByLibBlockNumNode) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *ByLibBlockNum) rotateRight(node *ByLibBlockNumNode) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *ByLibBlockNum) replaceNode(old *ByLibBlockNumNode, new *ByLibBlockNumNode) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *ByLibBlockNum) insertCase1(node *ByLibBlockNumNode) {
	if node.Parent == nil {
		node.color = blackByLibBlockNum
	} else {
		tree.insertCase2(node)
	}
}

func (tree *ByLibBlockNum) insertCase2(node *ByLibBlockNumNode) {
	if nodeColorByLibBlockNum(node.Parent) == blackByLibBlockNum {
		return
	}
	tree.insertCase3(node)
}

func (tree *ByLibBlockNum) insertCase3(node *ByLibBlockNumNode) {
	uncle := node.uncle()
	if nodeColorByLibBlockNum(uncle) == redByLibBlockNum {
		node.Parent.color = blackByLibBlockNum
		uncle.color = blackByLibBlockNum
		node.grandparent().color = redByLibBlockNum
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *ByLibBlockNum) insertCase4(node *ByLibBlockNumNode) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *ByLibBlockNum) insertCase5(node *ByLibBlockNumNode) {
	node.Parent.color = blackByLibBlockNum
	grandparent := node.grandparent()
	grandparent.color = redByLibBlockNum
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *ByLibBlockNumNode) maximumNode() *ByLibBlockNumNode {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *ByLibBlockNum) deleteCase1(node *ByLibBlockNumNode) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *ByLibBlockNum) deleteCase2(node *ByLibBlockNumNode) {
	sibling := node.sibling()
	if nodeColorByLibBlockNum(sibling) == redByLibBlockNum {
		node.Parent.color = redByLibBlockNum
		sibling.color = blackByLibBlockNum
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *ByLibBlockNum) deleteCase3(node *ByLibBlockNumNode) {
	sibling := node.sibling()
	if nodeColorByLibBlockNum(node.Parent) == blackByLibBlockNum &&
		nodeColorByLibBlockNum(sibling) == blackByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Left) == blackByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Right) == blackByLibBlockNum {
		sibling.color = redByLibBlockNum
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *ByLibBlockNum) deleteCase4(node *ByLibBlockNumNode) {
	sibling := node.sibling()
	if nodeColorByLibBlockNum(node.Parent) == redByLibBlockNum &&
		nodeColorByLibBlockNum(sibling) == blackByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Left) == blackByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Right) == blackByLibBlockNum {
		sibling.color = redByLibBlockNum
		node.Parent.color = blackByLibBlockNum
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *ByLibBlockNum) deleteCase5(node *ByLibBlockNumNode) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColorByLibBlockNum(sibling) == blackByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Left) == redByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Right) == blackByLibBlockNum {
		sibling.color = redByLibBlockNum
		sibling.Left.color = blackByLibBlockNum
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColorByLibBlockNum(sibling) == blackByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Right) == redByLibBlockNum &&
		nodeColorByLibBlockNum(sibling.Left) == blackByLibBlockNum {
		sibling.color = redByLibBlockNum
		sibling.Right.color = blackByLibBlockNum
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *ByLibBlockNum) deleteCase6(node *ByLibBlockNumNode) {
	sibling := node.sibling()
	sibling.color = nodeColorByLibBlockNum(node.Parent)
	node.Parent.color = blackByLibBlockNum
	if node == node.Parent.Left && nodeColorByLibBlockNum(sibling.Right) == redByLibBlockNum {
		sibling.Right.color = blackByLibBlockNum
		tree.rotateLeft(node.Parent)
	} else if nodeColorByLibBlockNum(sibling.Left) == redByLibBlockNum {
		sibling.Left.color = blackByLibBlockNum
		tree.rotateRight(node.Parent)
	}
}

func nodeColorByLibBlockNum(node *ByLibBlockNumNode) colorByLibBlockNum {
	if node == nil {
		return blackByLibBlockNum
	}
	return node.color
}

//////////////iterator////////////////

func (tree *ByLibBlockNum) makeIterator(fn *MultiIndexNode) IteratorByLibBlockNum {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*ByLibBlockNumNode); ok {
			return IteratorByLibBlockNum{tree: tree, node: n, position: betweenByLibBlockNum}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

// Iterator holding the iterator's state
type IteratorByLibBlockNum struct {
	tree     *ByLibBlockNum
	node     *ByLibBlockNumNode
	position positionByLibBlockNum
}

type positionByLibBlockNum byte

const (
	beginByLibBlockNum, betweenByLibBlockNum, endByLibBlockNum positionByLibBlockNum = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *ByLibBlockNum) Iterator() IteratorByLibBlockNum {
	return IteratorByLibBlockNum{tree: tree, node: nil, position: beginByLibBlockNum}
}

func (tree *ByLibBlockNum) Begin() IteratorByLibBlockNum {
	itr := IteratorByLibBlockNum{tree: tree, node: nil, position: beginByLibBlockNum}
	itr.Next()
	return itr
}

func (tree *ByLibBlockNum) End() IteratorByLibBlockNum {
	return IteratorByLibBlockNum{tree: tree, node: nil, position: endByLibBlockNum}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *IteratorByLibBlockNum) Next() bool {
	if iterator.position == endByLibBlockNum {
		goto end
	}
	if iterator.position == beginByLibBlockNum {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Left {
				goto between
			}
			node = iterator.node
		}
	}

end:
	iterator.node = nil
	iterator.position = endByLibBlockNum
	return false

between:
	iterator.position = betweenByLibBlockNum
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *IteratorByLibBlockNum) Prev() bool {
	if iterator.position == beginByLibBlockNum {
		goto begin
	}
	if iterator.position == endByLibBlockNum {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left
		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Right {
				goto between
			}
			node = iterator.node
			//if iterator.tree.Comparator(node.Key, iterator.node.Key) >= 0 {
			//	goto between
			//}
		}
	}

begin:
	iterator.node = nil
	iterator.position = beginByLibBlockNum
	return false

between:
	iterator.position = betweenByLibBlockNum
	return true
}

func (iterator IteratorByLibBlockNum) HasNext() bool {
	return iterator.position != endByLibBlockNum
}

func (iterator *IteratorByLibBlockNum) HasPrev() bool {
	return iterator.position != beginByLibBlockNum
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorByLibBlockNum) Value() BlockStatePtr {
	return *iterator.node.value()
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator IteratorByLibBlockNum) Key() ByLibBlockNumComposite {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *IteratorByLibBlockNum) Begin() {
	iterator.node = nil
	iterator.position = beginByLibBlockNum
}

func (iterator IteratorByLibBlockNum) IsBegin() bool {
	return iterator.position == beginByLibBlockNum
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *IteratorByLibBlockNum) End() {
	iterator.node = nil
	iterator.position = endByLibBlockNum
}

func (iterator IteratorByLibBlockNum) IsEnd() bool {
	return iterator.position == endByLibBlockNum
}

// Delete remove the node which pointed by the iterator
// Modifies the state of the iterator.
func (iterator *IteratorByLibBlockNum) Delete() {
	node := iterator.node
	//iterator.Prev()
	iterator.tree.remove(node)
}

func (tree *ByLibBlockNum) inPlace(n *ByLibBlockNumNode) bool {
	prev := IteratorByLibBlockNum{tree, n, betweenByLibBlockNum}
	next := IteratorByLibBlockNum{tree, n, betweenByLibBlockNum}
	prev.Prev()
	next.Next()

	var (
		prevResult int
		nextResult int
	)

	if prev.IsBegin() {
		prevResult = 1
	} else {
		prevResult = ByLibBlockNumCompare(n.Key, prev.Key())
	}

	if next.IsEnd() {
		nextResult = -1
	} else {
		nextResult = ByLibBlockNumCompare(n.Key, next.Key())
	}

	return (true && prevResult >= 0 && nextResult <= 0) ||
		(!true && prevResult > 0 && nextResult < 0)
}
