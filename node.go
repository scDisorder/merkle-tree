package merkletree

import "fmt"

type Content interface {
	CalculateHash() ([]byte, error)
	Equals(other Content) (bool, error)
}

type Node struct {
	Tree   *MerkleTree
	Parent *Node
	Left   *Node
	Right  *Node
	Hash   []byte
	C      Content
	leaf   bool
	dup    bool
}

// verify goes down the tree and calculates hash at each level until finds a leaf, which returns resulting hash
func (n *Node) verify() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}

	right, err := n.Right.verify()
	if err != nil {
		return nil, err
	}

	left, err := n.Left.verify()
	if err != nil {
		return nil, err
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(left, right...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// hash calculates node hash
func (n *Node) hash() ([]byte, error) {
	if n.leaf {
		return n.C.CalculateHash()
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// returns a string representation of the node.
func (n *Node) String() string {
	return fmt.Sprintf("%t %t %v %s", n.leaf, n.dup, n.Hash, n.C)
}
