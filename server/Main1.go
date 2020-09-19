package server

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"uv-web-server/andlabs/dialog"
)

type NewAreaHandler struct {
	ui.AreaHandler
}

func Main1() {
	err := ui.Main(func() {
		// 生成：窗口（标题，宽度，高度，是否有 菜单 控件）
		window := ui.NewWindow(`UV轻量级Web服务器`, 600, 400, false)
		// 设置：窗口关闭时
		window.OnClosing(func(*ui.Window) bool {
			// 窗体关闭
			ui.Quit()
			return true
		})
		// 窗体显示
		window.Show()
		window.SetMargined(true)
		//area:=ui.NewScrollingArea(NewAreaHandler{},600,400)
		//window.SetChild(area)

		box := ui.NewVerticalBox()
		box.SetPadded(true)

		sf := ui.NewHorizontalBox()
		sf.SetPadded(true)
		sl := ui.NewLabel("资源目录：")
		sf.Append(sl, false)
		path := ui.NewEntry()
		sf.Append(path, false)
		cb := ui.NewButton("选择")
		cb.OnClicked(func(button *ui.Button) {
			fd := dialog.NewFolderDialog()
			fd.Show()
		})
		sf.Append(cb, false)
		box.Append(sf, false)

		name := ui.NewEntry()
		greeting := ui.NewLabel(``)
		button := ui.NewButton(`欢迎`)
		button.OnClicked(func(*ui.Button) {
			greeting.SetText(`你好，` + name.Text() + `！`)
		})

		box.Append(ui.NewCombobox(), false)
		box.Append(ui.NewDatePicker(), false)
		box.Append(ui.NewDateTimePicker(), false)
		box.Append(ui.NewEditableCombobox(), false)
		box.Append(ui.NewFontButton(), false)
		box.Append(ui.NewForm(), false)
		//box.Append(ui.dia("ok"),false)
		box.Append(ui.NewCheckbox("hello"), false)
		box.Append(ui.NewColorButton(), false)
		box.Append(ui.NewSpinbox(0, 100), false)
		box.Append(ui.NewLabel(`请输入你的名字：`), false)
		box.Append(name, false)
		box.Append(greeting, false)
		box.Append(button, false)

		// 窗口容器绑定
		window.SetChild(box)

	})
	if err != nil {
		panic(err)
	}

}
