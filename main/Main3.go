package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strings"
	"uv-web-server/walk/dialog"
)

func main3() {
	var mw *walk.MainWindow

	//d:=walk.NewDialog(path)

	win := MainWindow{
		AssignTo: &mw,
		Title:    "UV轻量级Web服务器",
		Size:     Size{600, 400},
		Layout:   VBox{},
	}
	var path *walk.LineEdit

	var inTE, outTE *walk.TextEdit
	ws := []Widget{
		HSplitter{
			Children: []Widget{
				Label{Text: "资源路径：", StretchFactor: 1},
				LineEdit{AssignTo: &path, StretchFactor: 8},
				PushButton{Text: "选择", StretchFactor: 1,
					OnClicked: func() {
						path.SetText("ok")
						//walk.MsgBox(nil, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
						dialog.OpenFolderDialog(func(dialog *dialog.MyDialog) {
							println(dialog)
						})

					},
				},
			},
		},
		HSplitter{
			Children: []Widget{
				TextEdit{
					AssignTo:      &inTE,
					StretchFactor: 1,
				},
				TextEdit{
					AssignTo:      &outTE,
					ReadOnly:      true,
					StretchFactor: 5,
				},
			},
		},
		PushButton{
			Text: "SCREAM",
			OnClicked: func() {
				outTE.SetText(strings.ToUpper(inTE.Text()))
			},
		},
	}
	win.Children = ws
	win.Run()
}
