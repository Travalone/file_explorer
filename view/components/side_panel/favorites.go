package side_panel

import (
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
)

type FavoritesList struct {
	*packed_widgets.Tree

	feContext *store.FeContext
}

func (favoritesList *FavoritesList) RefreshData() {
	feConfig := favoritesList.feContext.GetFeConfig()
	favoritesDirs := make([]*packed_widgets.TreeNode, len(feConfig.Favorites))
	for i, favorites := range feConfig.Favorites {
		favoritesDirs[i] = &packed_widgets.TreeNode{
			NodeId: favorites.Path,
			Label:  favorites.Name,
			Depth:  0,
		}
	}
	favoritesList.Tree.SetData(&packed_widgets.TreeNode{
		NodeId:   "",
		Label:    "\\",
		Children: favoritesDirs,
	})
	favoritesList.Tree.Refresh()
}

func NewFavorites(feContext *store.FeContext) *FavoritesList {
	favoritesList := &FavoritesList{
		feContext: feContext,
		Tree:      packed_widgets.NewTree(nil),
	}

	favoritesList.RefreshData()

	setDirTreeEvents(feContext, favoritesList.Tree)

	feContext.FeConfig.BindListener(func() {
		favoritesList.RefreshData()
	})

	return favoritesList
}
