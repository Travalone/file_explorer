package model

import (
	"sort"
	"strings"
)

type Favorites struct {
	Path string `yaml:"path"`
	Name string `yaml:"name"`
}

type FileExplorerConfig struct {
	Root              string       `yaml:"root"`
	Favorites         []*Favorites `yaml:"favorites,omitempty"`
	FavoritesNotExist byte         `yaml:"favorites_not_exist"`
}

func (config *FileExplorerConfig) AddFavorites(path string, name string) {
	if config.Favorites == nil {
		config.Favorites = make([]*Favorites, 0)
	}
	config.Favorites = append(config.Favorites, &Favorites{
		Path: path,
		Name: name,
	})
	config.sortFavorites()
}

func (config *FileExplorerConfig) sortFavorites() {
	interfaces := make([]interface{}, len(config.Favorites))
	for i, item := range config.Favorites {
		interfaces[i] = item
	}

	sort.Slice(interfaces, func(i, j int) bool {
		a, b := interfaces[i].(*Favorites), interfaces[j].(*Favorites)
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	favorites := make([]*Favorites, len(interfaces))
	for i, item := range interfaces {
		favorites[i] = item.(*Favorites)
	}
	config.Favorites = favorites
}

func (config *FileExplorerConfig) DeleteFavorites(path string) {
	index := config.FindFavorites(path)
	if index < 0 {
		return
	}
	newFavorites := config.Favorites[:index]
	if index+1 < len(config.Favorites) {
		newFavorites = append(newFavorites, config.Favorites[index+1:]...)
	}
	config.Favorites = newFavorites
}

func (config *FileExplorerConfig) FindFavorites(path string) int {
	for index, favorites := range config.Favorites {
		if favorites.Path == path {
			return index
		}
	}
	return -1
}
