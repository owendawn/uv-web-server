package dialog

import "C"
import (
	"github.com/andlabs/ui"
	"uv-web-server/util"
)

type FolderDialog struct {
	ui.ControlBase
	win *ui.Window
}

func initUI(f *FolderDialog) {
	f.win = ui.NewWindow("选择文件夹", 400, 400, false)
	f.win.OnClosing(func(*ui.Window) bool {
		return true
	})
	hbox := ui.NewHorizontalBox()
	lbox := ui.NewVerticalBox()
	rbox := ui.NewVerticalBox()

	click := func(b *ui.Button) {
		println(b.Text())
	}

	l := util.GetDiskList()
	for i := l.Front(); i != nil; i = i.Next() {
		b := ui.NewButton(i.Value.(string))
		b.OnClicked(click)
		lbox.Append(b, false)
	}

	hbox.Append(lbox, true)
	hbox.Append(rbox, false)
	f.win.SetChild(hbox)
	util.GetLogicalDrives()

}

func (c *FolderDialog) Show() {
	c.win.Show()
}

func NewFolderDialog() *FolderDialog {
	f := FolderDialog{}
	if f.win == nil {
		initUI(&f)
	}
	return &f
}
