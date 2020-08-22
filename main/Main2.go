package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"os"
	"uv-web-server/dialog"

	"fyne.io/fyne/widget"
)

func main2() {
	dir, _ := os.Getwd()
	fmt.Println("当前路径：", dir)
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	res, e := fyne.LoadResourceFromPath(dir + "/asserts/uv.png")
	print(e)
	a.SetIcon(res)
	w := a.NewWindow("uv-web-server")
	w.Resize(fyne.NewSize(600, 500))

	pathLabel := widget.NewLabelWithStyle("WEB ROOT:", fyne.TextAlignLeading, fyne.TextStyle{
		Bold:   true, // Should text be bold
		Italic: false,
	})
	//pathLabel.TextStyle.Bold
	pathLabel.Resize(fyne.NewSize(100, 100))

	pathInput := widget.NewEntry()
	//pathInput.Resize(fyne.NewSize(400,100))
	//pathInput.Move(fyne.NewPos(300,100))

	fo := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if file != nil {
			pathInput.SetText(file.URI().String())
		}
	}, w)
	pathBtn := widget.NewButton("Choose", func() {
		fo.Show()
	})
	//pathBtn.Resize(fyne.NewSize(100,10))

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), pathLabel, pathInput, pathBtn)

	box := widget.NewVBox()
	container.Resize(fyne.NewSize(600, 500))
	box.Append(container)
	w.SetContent(box)

	w.ShowAndRun()
}
