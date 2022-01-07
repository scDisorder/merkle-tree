package merkletree

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
)

type MerkleTree struct {
	Root         *Node
	merkleRoot   []byte
	Leafs        []*Node
	hashStrategy func() hash.Hash
}

// MerkleRoot returns hash of the root node
func (m *MerkleTree) MerkleRoot() []byte {
	return m.merkleRoot
}

// MerklePath returns merkle path and indexes (left or right leaf)
func (m *MerkleTree) MerklePath(content Content) ([][]byte, []int64, error) {
	for _, l := range m.Leafs {
		ok, err := l.C.Equals(content)
		if err != nil {
			return nil, nil, err
		}

		if ok {
			leafParent := l.Parent

			var merklePath [][]byte
			var indices []int64

			for leafParent != nil {
				if bytes.Equal(leafParent.Left.Hash, leafParent.Hash) {
					merklePath = append(merklePath, leafParent.Right.Hash)
					indices = append(indices, 1)
				} else {
					merklePath = append(merklePath, leafParent.Left.Hash)
					indices = append(indices, 0)
				}

				l = leafParent
				leafParent = leafParent.Parent
			}

			return merklePath, indices, nil
		}
	}

	return nil, nil, nil
}

// Rebuilt is a helper func which rebuilds the tree with existing contents
func (m *MerkleTree) Rebuilt() error {
	var contents []Content
	for _, c := range m.Leafs {
		contents = append(contents, c.C)
	}

	root, leafs, err := buildWithContent(contents, m)
	if err != nil {
		return err
	}

	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash

	return nil
}

// RebuiltWith is a helper func which rebuilds the tree by replacing content passed as argument
func (m *MerkleTree) RebuiltWith(contents []Content) error {
	root, leafs, err := buildWithContent(contents, m)
	if err != nil {
		return err
	}

	m.Root = root
	m.Leafs = leafs
	m.merkleRoot = root.Hash

	return nil
}

// VerifyTree validates hash on each level
func (m *MerkleTree) VerifyTree() (bool, error) {
	merkleRoot, err := m.Root.verify()
	if err != nil {
		return false, err
	}

	if bytes.Compare(merkleRoot, m.MerkleRoot()) == 0 {
		return true, nil
	}

	return false, nil
}

// VerifyContent indicates if given content is a part of the tree
func (m *MerkleTree) VerifyContent(c Content) (bool, error) {
	for _, l := range m.Leafs {
		ok, err := l.C.Equals(c)
		if err != nil {
			return false, err
		}

		if ok {
			leafParent := l.Parent

			for leafParent != nil {
				h := m.hashStrategy()

				right, err := leafParent.Right.verify()
				if err != nil {
					return false, err
				}

				left, err := leafParent.Left.verify()
				if err != nil {
					return false, err
				}

				if _, err := h.Write(append(left, right...)); err != nil {
					return false, err
				}

				if bytes.Compare(h.Sum(nil), leafParent.Hash) != 0 {
					return false, nil
				}

				leafParent = leafParent.Parent
			}

			return true, nil
		}
	}

	return false, nil
}

// returns merkle tree string representation
func (m *MerkleTree) String() string {
	s := ""
	for _, l := range m.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	return s
}

func NewTree(contents []Content) (*MerkleTree, error) {
	defaultHashStrategy := sha256.New

	tree := &MerkleTree{
		hashStrategy: defaultHashStrategy,
	}

	root, leafs, err := buildWithContent(contents, tree)
	if err != nil {
		return nil, err
	}

	tree.Root = root
	tree.Leafs = leafs
	tree.merkleRoot = root.Hash

	return tree, nil
}

func NewTreeWithHashStrategy(contents []Content, hashStrategy func() hash.Hash) (*MerkleTree, error) {
	tree := &MerkleTree{
		hashStrategy: hashStrategy,
	}

	root, leafs, err := buildWithContent(contents, tree)
	if err != nil {
		return nil, err
	}

	tree.Root = root
	tree.Leafs = leafs
	tree.merkleRoot = root.Hash

	return tree, nil
}
