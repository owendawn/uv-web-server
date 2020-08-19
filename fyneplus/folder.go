package fyneplus

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	//"go/ast"
)
import (
	//"io/ioutil"
	"os"
	"path/filepath"
	//"strings"

	//"fyne.io/fyne"
	//"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	//"fyne.io/fyne/widget"
)

type textWidget interface {
	fyne.Widget
	SetText(string)
}
type folderDialog struct {
	win        *widget.PopUp
	folder     *FolderDialog
	dismiss    *widget.Button
	folderName textWidget
	breadcrumb *widget.Box
	fileScroll *widget.ScrollContainer
	open       *widget.Button
	dir        string
	selected   *fileDialogItem
}

type FolderDialog struct {
	dialog           *folderDialog
	callback         interface{}
	parent           fyne.Window
	save             bool
	onClosedCallback func(bool)
	dismissText      string
}

func (f *folderDialog) makeUI() fyne.CanvasObject {
	if f.folder.save {
		saveName := widget.NewEntry()
		saveName.OnChanged = func(s string) {
			if s == "" {
				f.open.Disable()
			} else {
				f.open.Enable()
			}
		}
		f.folderName = saveName
	} else {
		f.folderName = widget.NewLabel("")
	}

	label := "Open"
	if f.folder.save {
		label = "Save"
	}
	f.open = widget.NewButton(label, func() {
		if f.folder.callback == nil {
			f.win.Hide()
			if f.folder.onClosedCallback != nil {
				f.folder.onClosedCallback(false)
			}
			return
		}

		if f.folder.save {
			callback := f.folder.callback.(func(fyne.URIWriteCloser, error))
			name := f.folderName.(*widget.Entry).Text
			path := filepath.Join(f.dir, name)

			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				f.win.Hide()
				if f.folder.onClosedCallback != nil {
					f.folder.onClosedCallback(true)
				}
				callback(storage.SaveFileToURI(storage.NewURI("file://" + path)))
				return
			} else if info.IsDir() {
				//ShowInformation("Cannot overwrite",
				//	"Files cannot replace a directory,\ncheck the file name and try again", f.file.parent)
				return
			}

			//ShowConfirm("Overwrite?", "Are you sure you want to overwrite the file\n"+name+"?",
			//	func(ok bool) {
			//		f.win.Hide()
			//		if !ok {
			//			callback(nil, nil)
			//			return
			//		}
			//
			//		callback(storage.SaveFileToURI(storage.NewURI("file://" + path)))
			//		if f.folder.onClosedCallback != nil {
			//			f.folder.onClosedCallback(true)
			//		}
			//	}, f.folder.parent)
		} else if f.selected != nil {
			callback := f.folder.callback.(func(fyne.URIReadCloser, error))
			f.win.Hide()
			if f.folder.onClosedCallback != nil {
				f.folder.onClosedCallback(true)
			}
			callback(storage.OpenFileFromURI(storage.NewURI("file://" + f.selected.path)))
		}
	})
	f.open.Style = widget.PrimaryButton
	f.open.Disable()
	dismissLabel := "Cancel"
	if f.folder.dismissText != "" {
		dismissLabel = f.folder.dismissText
	}
	f.dismiss = widget.NewButton(dismissLabel, func() {
		f.win.Hide()
		if f.folder.onClosedCallback != nil {
			f.folder.onClosedCallback(false)
		}
		if f.folder.callback != nil {
			if f.folder.save {
				f.folder.callback.(func(fyne.URIWriteCloser, error))(nil, nil)
			} else {
				f.folder.callback.(func(fyne.URIReadCloser, error))(nil, nil)
			}
		}
	})
	buttons := widget.NewHBox(f.dismiss, f.open)
	footer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, buttons),
		buttons, widget.NewHScrollContainer(f.folderName))

	//f.files = fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(fileIconCellWidth,
	//	fileIconSize+theme.Padding()+fileTextSize)),
	//)
	//f.fileScroll = widget.NewScrollContainer(f.files)
	//verticalExtra := int(float64(fileIconSize) * 0.25)
	//f.fileScroll.SetMinSize(fyne.NewSize(fileIconCellWidth*2+theme.Padding(),
	//	(fileIconSize+fileTextSize)+theme.Padding()*2+verticalExtra))
	//
	f.breadcrumb = widget.NewHBox()
	scrollBread := widget.NewScrollContainer(f.breadcrumb)
	body := fyne.NewContainerWithLayout(layout.NewBorderLayout(scrollBread, nil, nil, nil),
		scrollBread, f.fileScroll)
	header := widget.NewLabelWithStyle(label+" File", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	favorites := widget.NewGroup("Favorites", f.loadFavorites()...)
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(header, footer, favorites, nil),
		favorites, header, footer, body)
}

func (f *folderDialog) loadFavorites() []fyne.CanvasObject {
	//home, _ := os.UserHomeDir()
	places := []fyne.CanvasObject{
		//makeFavoriteButton("Home", theme.HomeIcon(), func() {
		//	f.setDirectory(home)
		//}),
		//makeFavoriteButton("Documents", theme.DocumentIcon(), func() {
		//	f.setDirectory(filepath.Join(home, "Documents"))
		//}),
		//makeFavoriteButton("Downloads", theme.DownloadIcon(), func() {
		//	f.setDirectory(filepath.Join(home, "Downloads"))
		//}),
	}

	places = append(places, f.loadPlaces()...)
	return places
}

func (f *folderDialog) loadPlaces() []fyne.CanvasObject {
	var places []fyne.CanvasObject

	//for _, drive := range listDrives() {
	//	driveRoot := drive + string(os.PathSeparator) // capture loop var
	//	places = append(places, makeFavoriteButton(drive, theme.StorageIcon(), func() {
	//		f.setDirectory(driveRoot)
	//	}))
	//}
	return places
}
func showFile(folder *FolderDialog) *folderDialog {
	d := &folderDialog{folder: folder}
	ui := d.makeUI()

	//d.setDirectory(folder.effectiveStartingDir())

	size := ui.MinSize().Add(fyne.NewSize(fileIconCellWidth*2+theme.Padding()*4,
		(fileIconSize+fileTextSize)+theme.Padding()*4))

	d.win = widget.NewModalPopUp(ui, folder.parent.Canvas())
	d.win.Resize(size)

	d.win.Show()
	return d
}

// Show shows the file dialog.
func (f *FolderDialog) Show() {
	if f.dialog != nil {
		f.dialog.win.Show()
		return
	}
	f.dialog = showFile(f)
}

func newFolderDialog(callback func(fyne.URIReadCloser, error), parent fyne.Window) *FolderDialog {
	return &FolderDialog{}
}
