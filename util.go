package merkletree

import "fmt"

func buildIntermediate(nl []*Node, tree *MerkleTree) (*Node, error) {
	var nodes []*Node

	for i := 0; i < len(nl); i += 2 {
		h := tree.hashStrategy()

		var left, right = i, i + 1
		if i+1 == len(nl) {
			right = i
		}

		calculatedHash := append(nl[left].Hash, nl[right].Hash...)
		if _, err := h.Write(calculatedHash); err != nil {
			return nil, err
		}

		n := &Node{
			Tree:  tree,
			Left:  nl[left],
			Right: nl[right],
			Hash:  h.Sum(nil),
		}

		nodes = append(nodes, n)

		nl[left].Parent = n
		nl[right].Parent = n

		if len(nl) == 2 {
			return n, nil
		}
	}

	return buildIntermediate(nodes, tree)
}

func buildWithContent(contents []Content, tree *MerkleTree) (*Node, []*Node, error) {
	if len(contents) == 0 {
		return nil, nil, fmt.Errorf("cannot construct the tree with empty contents")
	}

	var leafs []*Node

	for _, c := range contents {
		hash, err := c.CalculateHash()
		if err != nil {
			return nil, nil, err
		}

		leafs = append(leafs, &Node{
			Tree: tree,
			Hash: hash,
			C:    c,
			leaf: true,
		})
	}

	if len(leafs)%2 == 1 {
		dup := &Node{
			Tree: tree,
			Hash: leafs[len(leafs)-1].Hash,
			C:    leafs[len(leafs)-1].C,
			leaf: true,
			dup:  true,
		}

		leafs = append(leafs, dup)
	}

	root, err := buildIntermediate(leafs, tree)
	if err != nil {
		return nil, nil, err
	}

	return root, leafs, nil
}
