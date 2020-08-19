package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(600, 500))

	hello := widget.NewLabel("Hello Fyne!")
	files := widget.NewLabel("file!")
	path := widget.NewLabel("Hello Fyne!")

	box := widget.NewVBox()
	//box.Resize(fyne.NewSize(600,500))

	box.Append(hello)
	box.Append(widget.NewButton("Hi!", func() {
		hello.SetText("Welcome :)")
	}))
	box.Append(path)
	fo := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		files.SetText(file.Name())
	}, w)
	box.Append(widget.NewButton("file!", func() {
		fo.Show()
	}))
	folder := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		files.SetText(file.Name())
	}, w)
	box.Append(widget.NewButton("folder!", func() {
		folder.Show()
	}))
	w.SetContent(box)

	w.ShowAndRun()
}
