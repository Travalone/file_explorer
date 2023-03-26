package packed_widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TreeNode struct {
	NodeId   string
	Label    string
	Depth    int
	Parent   *TreeNode
	Children []*TreeNode
}

type Tree struct {
	*widget.Tree
	Data           map[string]*TreeNode
	OnBranchOpen   func(*TreeNode)
	OnBranchClose  func(*TreeNode)
	OnDoubleTapped func(*TreeNode)
}

func (tree *Tree) SetData(rootNode *TreeNode) {
	if rootNode == nil {
		return
	}
	tree.Data = map[string]*TreeNode{rootNode.NodeId: rootNode}
	for _, rootChild := range tree.Data[rootNode.NodeId].Children {
		tree.Data[rootChild.NodeId] = rootChild
	}
}

func NewTree(rootNode *TreeNode) *Tree {
	tree := &Tree{}
	tree.SetData(rootNode)

	tree.Tree = widget.NewTree(
		func(nodeId widget.TreeNodeID) []widget.TreeNodeID {
			childIds := make([]string, len(tree.Data[nodeId].Children))
			for index, child := range tree.Data[nodeId].Children {
				childIds[index] = child.NodeId
			}
			return childIds
		},
		func(nodeId widget.TreeNodeID) bool {
			_, ok := tree.Data[nodeId]
			return ok
		},
		func(b bool) fyne.CanvasObject {
			return NewLabel("")
		},
		func(nodeId widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
			node.(*Label).SetText(tree.Data[nodeId].Label)
			node.(*Label).OnTapped = func() {
				if tree.IsBranchOpen(nodeId) {
					tree.CloseBranch(nodeId)
				} else {
					tree.OpenBranch(nodeId)
				}
			}
			if tree.OnDoubleTapped != nil {
				node.(*Label).OnDoubleTapped = func() {
					tree.OnDoubleTapped(tree.Data[nodeId])
				}
			}
		},
	)

	tree.OnBranchOpened = func(nodeId widget.TreeNodeID) {
		if tree.OnBranchOpen != nil {
			tree.OnBranchOpen(tree.Data[nodeId])
		}
	}
	tree.OnBranchClosed = func(nodeId widget.TreeNodeID) {
		if tree.OnBranchClose != nil {
			tree.OnBranchClose(tree.Data[nodeId])
		}
	}

	return tree
}
