package resource

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

//go:embed font/方圆新黑体.ttf
var FontTTF []byte

type FeTheme struct {
}

func (t *FeTheme) Font(s fyne.TextStyle) fyne.Resource {
	return fyne.NewStaticResource("font", FontTTF)
}
func (t *FeTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}
func (t *FeTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}
func (t *FeTheme) Size(s fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(s)
}
