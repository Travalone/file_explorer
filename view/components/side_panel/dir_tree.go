package side_panel

import (
	"file_explorer/common/logger"
	"file_explorer/common/utils"
	"file_explorer/service"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
)

func setDirTreeEvents(feContext *store.FeContext, tree *packed_widgets.Tree) {
	tree.OnBranchOpen = func(node *packed_widgets.TreeNode) {
		// 查询子目录
		subDirs, err := service.QuerySubDirs(node.NodeId)
		if err != nil {
			logger.Error("DirTree QuerySubDirs failed, err=%v", err)
			return
		}
		utils.StringsSort(subDirs)

		// 构造子节点
		tree.Data[node.NodeId].Children = make([]*packed_widgets.TreeNode, len(subDirs))
		for index, subDir := range subDirs {
			childNodeId := utils.PathJoin(node.NodeId, subDir)
			tree.Data[childNodeId] = &packed_widgets.TreeNode{
				NodeId: childNodeId,
				Label:  subDir,
				Depth:  node.Depth + 1,
				Parent: node,
			}
			tree.Data[node.NodeId].Children[index] = tree.Data[childNodeId]
		}
	}

	tree.OnDoubleTapped = func(node *packed_widgets.TreeNode) {
		fileTabContext := store.NewFileTabContext(node.NodeId, feContext)
		feContext.AddTab(fileTabContext)
	}
}

func NewDirTree(feContext *store.FeContext) *packed_widgets.Tree {
	feConfig := feContext.GetFeConfig()

	rootNode := &packed_widgets.TreeNode{
		NodeId: feConfig.Root, // 绝对路径
		Label:  feConfig.Root, // 显示名称
		Depth:  0,
	}

	tree := packed_widgets.NewTree(&packed_widgets.TreeNode{
		NodeId:   "",
		Label:    "\\",
		Children: []*packed_widgets.TreeNode{rootNode},
	})

	setDirTreeEvents(feContext, tree)

	return tree
}
