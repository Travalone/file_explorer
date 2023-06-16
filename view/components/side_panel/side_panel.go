package side_panel

import (
	"file_explorer/common"
	"file_explorer/common/utils"
	"file_explorer/view/packed_widgets"
	"file_explorer/view/store"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type SidePanel struct {
	*container.Scroll
	SearchBox *widget.Entry
	DirTree   *packed_widgets.Tree
	//OpList    fyne.CanvasObject
	Favorites *FavoritesList

	feContext *store.FeContext
}

func NewSidePanel(feContext *store.FeContext) *SidePanel {
	sidePanel := &SidePanel{
		SearchBox: NewSearchBox(feContext),
		DirTree:   NewDirTree(feContext),
		Favorites: NewFavorites(feContext),
		//OpList:    NewOpList(),
		feContext: feContext,
	}

	dirTreeTab := container.NewTabItem("目录树", sidePanel.DirTree.Tree) //不加.Tree时点击事件响应慢
	favoritesTab := container.NewTabItem("收藏夹", sidePanel.Favorites.Tree.Tree)
	//opListTab := container.NewTabItem("操作列表", sidePanel.OpList)
	sidePanelTabs := container.NewAppTabs(
		dirTreeTab,
		favoritesTab,
		//opListTab,
	)

	// 检查收藏夹目录是否都存在
	sidePanelTabs.OnSelected = func(item *container.TabItem) {
		if feContext.GetFeConfig().FavoritesNotExist == common.FavoritesNotExistIgnore ||
			item != favoritesTab || len(feContext.GetFeConfig().Favorites) == 0 {
			return
		}

		invalidFavorites := make([]string, 0)
		for _, favorites := range feContext.GetFeConfig().Favorites {
			if !utils.PathExists(favorites.Path) {
				invalidFavorites = append(invalidFavorites, favorites.Name)
			}
		}

		if len(invalidFavorites) == 0 {
			return
		}

		if feContext.GetFeConfig().FavoritesNotExist == common.FavoritesNotExistDelete {
			feContext.ClearInvalidFavorites()
			return
		}

		info := "以下收藏目录已失效: \n" + utils.Conv2Str(invalidFavorites)
		dialog.ShowConfirm("确认删除", info,
			func(confirm bool) {
				if confirm {
					feContext.ClearInvalidFavorites()
				}
			},
			packed_widgets.GetAppWindow())

	}

	sidePanel.Scroll = container.NewHScroll(
		container.NewBorder(sidePanel.SearchBox, nil, nil, nil, sidePanelTabs))

	return sidePanel
}
