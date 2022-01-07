package merkletree

import (
	"crypto/sha256"
	"log"
	"testing"
)

type TestContent struct {
	x string
}

// CalculateHash hashes the values of a TestContent
func (t TestContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(t.x)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t TestContent) Equals(other Content) (bool, error) {
	return t.x == other.(TestContent).x, nil
}

func Test(t *testing.T) {
	//Build list of Content to build tree
	var list []Content
	list = append(list, TestContent{x: "So"})
	list = append(list, TestContent{x: "Much"})
	list = append(list, TestContent{x: "Content"})
	list = append(list, TestContent{x: "Wow"})

	//Create a new Merkle Tree from the list of Content
	tree, err := NewTree(list)
	if err != nil {
		log.Fatal(err)
	}

	//Get the Merkle Root of the tree
	mr := tree.MerkleRoot()
	log.Println(mr)

	// Verify the entire tree (hashes for each node) is valid
	vt, err := tree.VerifyTree()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Verify Tree: ", vt)

	// Verify a specific content in the tree
	vc, err := tree.VerifyContent(list[0])
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Verify Content: ", vc)

	//String representation
	log.Println(t)
}
