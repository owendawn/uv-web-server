package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"

	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	a.SetIcon(fyne.NewStaticResource("icon", []byte("asserts/icon.png")))
	w := a.NewWindow("uv-web-server")
	w.Resize(fyne.NewSize(600, 500))

	pathLabel := widget.NewLabel("Web Resource:")
	//pathLabel.Resize(fyne.NewSize(100,100))

	pathInput := widget.NewEntry()
	//pathInput.Resize(fyne.NewSize(400,100))
	//pathInput.Move(fyne.NewPos(300,100))

	fo := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		pathInput.SetText(file.URI().String())
	}, w)
	pathBtn := widget.NewButton("choose", func() {
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
