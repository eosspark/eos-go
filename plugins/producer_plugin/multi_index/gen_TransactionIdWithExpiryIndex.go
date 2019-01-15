// Code generated by gotemplate. DO NOT EDIT.

package multi_index

import (
	"github.com/eosspark/eos-go/common/container/multiindex"
	"github.com/eosspark/eos-go/log"
)

// template type MultiIndex(SuperIndex,SuperNode,Value)

type TransactionIdWithExpiryIndex struct {
	super *ById
	count int
}

func NewTransactionIdWithExpiryIndex() *TransactionIdWithExpiryIndex {
	m := &TransactionIdWithExpiryIndex{}
	m.super = &ById{}
	m.super.init(m)
	return m
}

/*generic class*/

type TransactionIdWithExpiryIndexNode struct {
	super *ByIdNode
}

/*generic class*/

//method for MultiIndex
func (m *TransactionIdWithExpiryIndex) GetSuperIndex() interface{} { return m.super }
func (m *TransactionIdWithExpiryIndex) GetFinalIndex() interface{} { return nil }

func (m *TransactionIdWithExpiryIndex) GetIndex() interface{} {
	return nil
}

func (m *TransactionIdWithExpiryIndex) Size() int {
	return m.count
}

func (m *TransactionIdWithExpiryIndex) Clear() {
	m.super.clear()
	m.count = 0
}

func (m *TransactionIdWithExpiryIndex) Insert(v TransactionIdWithExpiry) bool {
	_, res := m.insert(v)
	return res
}

func (m *TransactionIdWithExpiryIndex) insert(v TransactionIdWithExpiry) (*TransactionIdWithExpiryIndexNode, bool) {
	fn := &TransactionIdWithExpiryIndexNode{}
	n, res := m.super.insert(v, fn)
	if res {
		fn.super = n
		m.count++
		return fn, true
	}
	return nil, false
}

func (m *TransactionIdWithExpiryIndex) Erase(iter multiindex.IteratorType) {
	m.super.erase_(iter)
}

func (m *TransactionIdWithExpiryIndex) erase(n *TransactionIdWithExpiryIndexNode) {
	m.super.erase(n.super)
	m.count--
}

func (m *TransactionIdWithExpiryIndex) Modify(iter multiindex.IteratorType, mod func(*TransactionIdWithExpiry)) bool {
	return m.super.modify_(iter, mod)
}

func (m *TransactionIdWithExpiryIndex) modify(mod func(*TransactionIdWithExpiry), n *TransactionIdWithExpiryIndexNode) (*TransactionIdWithExpiryIndexNode, bool) {
	defer func() {
		if e := recover(); e != nil {
			log.Error("#multi modify failed: %v", e)
			m.erase(n)
			m.count--
			panic(e)
		}
	}()
	mod(n.value())
	if sn, res := m.super.modify(n.super); !res {
		m.count--
		return nil, false
	} else {
		n.super = sn
		return n, true
	}
}

func (n *TransactionIdWithExpiryIndexNode) GetSuperNode() interface{} { return n.super }
func (n *TransactionIdWithExpiryIndexNode) GetFinalNode() interface{} { return nil }

func (n *TransactionIdWithExpiryIndexNode) value() *TransactionIdWithExpiry {
	return n.super.value()
}

/// IndexBase
type TransactionIdWithExpiryIndexBase struct {
	final *TransactionIdWithExpiryIndex
}

type TransactionIdWithExpiryIndexBaseNode struct {
	final *TransactionIdWithExpiryIndexNode
	pv    *TransactionIdWithExpiry
}

func (i *TransactionIdWithExpiryIndexBase) init(final *TransactionIdWithExpiryIndex) {
	i.final = final
}

func (i *TransactionIdWithExpiryIndexBase) clear() {}

func (i *TransactionIdWithExpiryIndexBase) GetSuperIndex() interface{} { return nil }

func (i *TransactionIdWithExpiryIndexBase) GetFinalIndex() interface{} { return i.final }

func (i *TransactionIdWithExpiryIndexBase) insert(v TransactionIdWithExpiry, fn *TransactionIdWithExpiryIndexNode) (*TransactionIdWithExpiryIndexBaseNode, bool) {
	return &TransactionIdWithExpiryIndexBaseNode{fn, &v}, true
}

func (i *TransactionIdWithExpiryIndexBase) erase(n *TransactionIdWithExpiryIndexBaseNode) {
	n.pv = nil
}

func (i *TransactionIdWithExpiryIndexBase) erase_(iter multiindex.IteratorType) {
	log.Warn("erase iterator doesn't match all index")
}

func (i *TransactionIdWithExpiryIndexBase) modify(n *TransactionIdWithExpiryIndexBaseNode) (*TransactionIdWithExpiryIndexBaseNode, bool) {
	return n, true
}

func (i *TransactionIdWithExpiryIndexBase) modify_(iter multiindex.IteratorType, mod func(*TransactionIdWithExpiry)) bool {
	log.Warn("modify iterator doesn't match all index")
	return false
}

func (n *TransactionIdWithExpiryIndexBaseNode) value() *TransactionIdWithExpiry {
	return n.pv
}
