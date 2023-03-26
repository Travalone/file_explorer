package side_panel

import (
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SidePanel struct {
	*fyne.Container
	SearchBox *widget.Entry
	DirTree   *packed_widgets.Tree
	Favorites *FavoritesList

	feContext *store.FeContext
}

func NewSidePanel(feContext *store.FeContext) *SidePanel {
	sidePanel := &SidePanel{
		SearchBox: widget.NewEntry(),
		DirTree:   NewDirTree(feContext),
		Favorites: NewFavorites(feContext),
		feContext: feContext,
	}

	sidePanelTabs := container.NewAppTabs(
		container.NewTabItem("目录树", sidePanel.DirTree.Tree), //不加.Tree时点击事件响应慢
		container.NewTabItem("收藏夹", sidePanel.Favorites.Tree.Tree),
	)

	sidePanel.Container = container.NewBorder(sidePanel.SearchBox, nil, nil, nil, sidePanelTabs)

	return sidePanel
}
