// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package dialog

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
	"uv-web-server/util"
)

type myDialogUI struct {
}

type MyDialog struct {
	*walk.Dialog
	ui       myDialogUI
	callback func(dialog *MyDialog)
	Path     walk.LineEdit
	ListBox  *walk.ListBox
	Items    []MyListEntry
	sEventId int
	aEventId int
	styler   *MyStyler
}

type MyListModel struct {
	walk.ReflectListModelBase
	items []MyListEntry
}

type MyListEntry struct {
	idx       string
	img       walk.Image
	timestamp time.Time
	message   string
	path      string
}

type MyStyler struct {
	myDialog            MyDialog
	lb                  **walk.ListBox
	canvas              *walk.Canvas
	model               *MyListModel
	font                *walk.Font
	dpi2StampSize       map[int]walk.Size
	widthDPI2WsPerLine  map[widthDPI]int
	textWidthDPI2Height map[textWidthDPI]int // in native pixels
}
type widthDPI struct {
	width int // in native pixels
	dpi   int
}

type textWidthDPI struct {
	text  string
	width int // in native pixels
	dpi   int
}

func (s *MyStyler) ItemHeightDependsOnWidth() bool {
	return true
}

func (s *MyStyler) DefaultItemHeight() int {
	dpi := (*s.lb).DPI()
	marginV := walk.IntFrom96DPI(marginV96dpi, dpi)

	return s.StampSize().Height + marginV*2
}

const (
	marginH96dpi int = 6
	marginV96dpi int = 2
	lineW96dpi   int = 1
)

func (s *MyStyler) ItemHeight(index, width int) int {
	dpi := (*s.lb).DPI()
	marginH := walk.IntFrom96DPI(marginH96dpi, dpi+100)
	marginV := walk.IntFrom96DPI(marginV96dpi, dpi)
	lineW := walk.IntFrom96DPI(lineW96dpi, dpi)

	msg := s.model.items[index].message

	twd := textWidthDPI{msg, width, dpi}

	if height, ok := s.textWidthDPI2Height[twd]; ok {
		return height + marginV*2
	}

	canvas, err := s.Canvas()
	if err != nil {
		return 0
	}

	stampSize := s.StampSize()

	wd := widthDPI{width, dpi}
	wsPerLine, ok := s.widthDPI2WsPerLine[wd]
	if !ok {
		bounds, _, err := canvas.MeasureTextPixels("W", (*s.lb).Font(), walk.Rectangle{Width: 9999999}, walk.TextCalcRect)
		if err != nil {
			return 0
		}
		wsPerLine = (width - marginH*4 - lineW - stampSize.Width) / bounds.Width
		s.widthDPI2WsPerLine[wd] = wsPerLine
	}

	if len(msg) <= wsPerLine {
		s.textWidthDPI2Height[twd] = stampSize.Height
		return stampSize.Height + marginV*2
	}

	bounds, _, err := canvas.MeasureTextPixels(msg, (*s.lb).Font(), walk.Rectangle{Width: width - marginH*4 - lineW - stampSize.Width, Height: 255}, walk.TextEditControl|walk.TextWordbreak|walk.TextEndEllipsis)
	if err != nil {
		return 0
	}

	s.textWidthDPI2Height[twd] = bounds.Height

	return bounds.Height + marginV*2
}

func (s *MyStyler) StyleItem(style *walk.ListItemStyle) {
	if canvas := style.Canvas(); canvas != nil {

		if style.Index()%2 == 1 && style.BackgroundColor == walk.Color(win.GetSysColor(win.COLOR_WINDOW)) {
			style.BackgroundColor = walk.Color(win.GetSysColor(win.COLOR_BTNFACE))
			if err := style.DrawBackground(); err != nil {
				return
			}
		}

		pen, err := walk.NewCosmeticPen(walk.PenSolid, style.LineColor)
		if err != nil {
			return
		}
		defer pen.Dispose()

		dpi := (*s.lb).DPI()
		marginH := walk.IntFrom96DPI(marginH96dpi, dpi)
		marginV := walk.IntFrom96DPI(marginV96dpi, dpi)
		lineW := walk.IntFrom96DPI(lineW96dpi, dpi)

		item := s.model.items[style.Index()]
		b := style.BoundsPixels()

		b.X += marginH
		b.Y += marginV

		style.Canvas().DrawImageStretchedPixels(item.img, walk.Rectangle{b.X, b.Y, 15, 15})

		b.X += 15 + marginH
		canvas.DrawLinePixels(pen, walk.Point{b.X, b.Y - marginV}, walk.Point{b.X, b.Y - marginV + b.Height})

		b.X += marginH
		style.DrawText(item.timestamp.Format("2006-01-02 15:04:05.999"), b, walk.TextEditControl|walk.TextWordbreak)

		stampSize := s.StampSize()

		b.X += marginH + stampSize.Width
		canvas.DrawLinePixels(pen, walk.Point{b.X, b.Y - marginV}, walk.Point{b.X, b.Y - marginV + b.Height})

		b.X += marginH
		style.DrawText(item.message, b, walk.TextEditControl|walk.TextExpandTabs|walk.TextEndEllipsis)

		b.X += marginH
		//b.Width -= stampSize.Width + marginH*4 + lineW

		b.X += stampSize.Width + marginH*2 + lineW

	}
}

func (s *MyStyler) StampSize() walk.Size {
	dpi := (*s.lb).DPI()

	stampSize, ok := s.dpi2StampSize[dpi]
	if !ok {
		canvas, err := s.Canvas()
		if err != nil {
			return walk.Size{}
		}

		bounds, _, err := canvas.MeasureTextPixels("2006-01-02 15:04:05.999", (*s.lb).Font(), walk.Rectangle{Width: 9999999}, walk.TextCalcRect)
		if err != nil {
			return walk.Size{}
		}

		stampSize = bounds.Size()
		s.dpi2StampSize[dpi] = stampSize
	}

	return stampSize
}

func (s *MyStyler) Canvas() (*walk.Canvas, error) {
	if s.canvas != nil {
		return s.canvas, nil
	}

	canvas, err := (*s.lb).CreateCanvas()
	if err != nil {
		return nil, err
	}
	s.canvas = canvas
	(*s.lb).AddDisposable(canvas)

	return canvas, nil
}

func InitEvent(s *MyStyler) {
	if s.myDialog.aEventId >= 0 {
		(*s.lb).ItemActivated().Detach(s.myDialog.aEventId)
		s.myDialog.aEventId = -1
	}
	if s.myDialog.sEventId >= 0 {
		(*s.lb).SelectedIndexesChanged().Detach(s.myDialog.sEventId)
		s.myDialog.sEventId = -1
	}

	s.myDialog.aEventId = (*s.lb).ItemActivated().Attach(func() {
		i := (*s.lb).CurrentIndex()
		it := s.myDialog.Items[i]
		p := s.myDialog.Path
		println("p=" + it.path)
		if it.message == ".." {
			p.SetText(it.path)
			RefreshListBox(&s.myDialog)
		} else if util.IsDir(it.path) {
			p.SetText(it.path)
			RefreshListBox(&s.myDialog)
		}
	})
	s.myDialog.sEventId = (*s.lb).SelectedIndexesChanged().Attach(func() {
		i := (*s.lb).CurrentIndex()

		println("111-" + s.myDialog.Items[i].message)
		it := s.myDialog.Items[i]
		p := s.myDialog.Path
		f := p.Text() + "/" + it.message
		if util.IsDir(f) && it.message != ".." {
			//p.SetText(f)
		}
	})
}

func (m *MyListModel) Items() interface{} {
	return m.items
}

func (dlg *MyDialog) setState(state walk.PIState) {
	if err := dlg.ProgressIndicator().SetState(state); err != nil {
		log.Print(err)
	}
}

func RunMyDialog(owner walk.Form, callback func(dialog *MyDialog)) (int, error) {
	dlg := new(MyDialog)
	dlg.Dialog, _ = walk.NewDialog(owner)
	dlg.aEventId = -1
	dlg.sEventId = -1

	succeeded := false
	defer func() {
		if !succeeded {
			dlg.Dispose()
		}
	}()

	dlg.callback = callback
	dlg.SetName("Dialog")
	dlg.SetTitle("选择文件夹")
	dlg.SetLayout(walk.NewVBoxLayout())
	dlg.SetClientSize(walk.Size{700, 300})
	dlg.SetSize(walk.Size{700, 500})
	dlg.SetIcon(util.NewSystemIcon())

	hs1, _ := walk.NewHSplitter(dlg)
	dlg.Children().Add(hs1)
	l1, _ := walk.NewLabel(hs1)
	l1.SetText("文件目录：")
	nl, _ := walk.NewLineEdit(hs1)
	nl.SetText(util.NewUserHomePath())
	dlg.Path = *nl
	dlg.Path.SetWidth(500)
	dlg.Path.SetWidthPixels(500)
	btn1, _ := walk.NewPushButton(hs1)
	btn1.SetText("确定")
	btn1.Clicked().Attach(func() {
		dlg.Dialog.Hide()
		callback(dlg)
	})
	hs1.Children().Add(l1)
	hs1.Children().Add(nl)
	hs1.Children().Add(btn1)

	hs, _ := walk.NewHSplitter(dlg)
	dlg.Children().Add(hs)
	lpart, _ := walk.NewVSplitter(hs)
	lpart.SetMinMaxSize(walk.Size{50, 100}, walk.Size{100, 400})
	lpart.SetSize(walk.Size{50, 400})
	cspace, _ := walk.NewHSeparator(hs)
	rpart, _ := walk.NewVSplitter(hs)
	hs.Children().Add(lpart)
	hs.Children().Add(cspace)
	hs.Children().Add(rpart)

	dBtn, _ := walk.NewPushButton(dlg)
	dBtn.SetName("Desktop")
	//dBtn.SetBounds(walk.Rectangle{40, 120, 100, 23})
	dBtn.SetText(`Desktop`)
	dBtn.SetMinMaxSize(walk.Size{0, 0}, walk.Size{50, 50})
	dBtn.Clicked().Attach(func() {
		//fmt.Println("SetState normal")
		//dlg.setState(walk.PINormal)
		dlg.Path.SetText(util.NewUserHomePath() + "\\Desktop")
		RefreshListBox(&dlg.styler.myDialog)
	})
	lpart.Children().Add(dBtn)
	hBtn, _ := walk.NewPushButton(dlg)
	hBtn.SetName("Home")
	//hBtn.SetBounds(walk.Rectangle{40, 120, 100, 23})
	hBtn.SetText(`Home`)
	hBtn.SetMinMaxSize(walk.Size{0, 0}, walk.Size{50, 50})
	hBtn.Clicked().Attach(func() {
		//fmt.Println("SetState error")
		//dlg.setState(walk.PIError)
		dlg.Path.SetText(util.NewUserHomePath())
		RefreshListBox(&dlg.styler.myDialog)
	})
	lpart.Children().Add(hBtn)
	choose := func(p string) {
		dlg.Path.SetText(p)
		RefreshListBox(&dlg.styler.myDialog)
	}
	devices := util.GetDiskList()
	for device := devices.Front(); nil != device; device = device.Next() {
		dbtn, _ := walk.NewPushButton(lpart)
		dbtn.SetText(device.Value.(string))
		dbtn.SetName(device.Value.(string))
		//dbtn.SetBounds(walk.Rectangle{40, 150, 161, 23})
		dbtn.SetMinMaxSize(walk.Size{0, 0}, walk.Size{50, 50})
		dbtn.Clicked().Attach(func() {
			choose(dbtn.Text())
		})
		lpart.Children().Add(dbtn)
	}

	cps, _ := walk.NewComposite(dlg)
	cps.SetLayout(walk.NewVBoxLayout())
	rpart.Children().Add(cps)
	rpart.SetLayout(walk.NewVBoxLayout())
	var lb *walk.ListBox
	dlg.ListBox = lb
	model := &MyListModel{items: dlg.Items}
	dlg.styler = &MyStyler{
		myDialog:            *dlg,
		lb:                  &dlg.ListBox,
		model:               model,
		dpi2StampSize:       make(map[int]walk.Size),
		widthDPI2WsPerLine:  make(map[widthDPI]int),
		textWidthDPI2Height: make(map[textWidthDPI]int),
	}
	declarative.ListBox{
		AssignTo:       &dlg.ListBox,
		MultiSelection: true,
		Model:          model,
		ItemStyler:     dlg.styler,
	}.Create(declarative.NewBuilder(cps))
	dlg.ListBox.SetMinMaxSize(walk.Size{400, 300}, walk.Size{500, 500})

	cps.Children().Add(dlg.ListBox)
	RefreshListBox(dlg)

	dlg.callback(dlg)

	succeeded = true
	return dlg.Run(), nil
}

func RefreshListBox(dlg *MyDialog) {
	ii, _ := walk.NewImageFromFileForDPI(util.NewSystemPath()+"/asserts/fileicon.jpg", 10)
	ii2, _ := walk.NewImageFromFileForDPI(util.NewSystemPath()+"/asserts/foldericon.jpg", 10)
	fs, _ := ioutil.ReadDir(dlg.Path.Text())
	dir := util.ReadDir(dlg.Path.Text())
	var items []MyListEntry
	thepath := strings.Replace(dlg.Path.Text(), "\\", "/", -1)
	if !strings.HasSuffix(thepath, ":") {
		idx := -1
		thepath2 := []rune(thepath)
		for i := utf8.RuneCountInString(thepath) - 1; i > 0; i-- {
			if string(thepath2[i:i+1]) == "/" {
				idx = i
				break
			}
		}
		thepath = util.Substr(thepath, 0, idx)
		items = append(items, MyListEntry{strconv.Itoa(1), ii2, dir.ModTime(), "..", thepath})
	}
	for i := 0; i < len(fs); i++ {
		f := fs[i]
		if f.Name() == "System Volume Information" {
			continue
		}
		if strings.HasPrefix(f.Name(), "$") {
			continue
		}
		if f.IsDir() {
			items = append(items, MyListEntry{strconv.Itoa(1 + len(items)), ii2, f.ModTime(), f.Name(), dlg.Path.Text() + "/" + f.Name()})
		} else {
			items = append(items, MyListEntry{strconv.Itoa(1 + len(items)), ii, time.Now(), f.Name(), dlg.Path.Text() + "/" + f.Name()})
		}
	}
	println(len(items))
	dlg.Items = items
	dlg.styler.myDialog = *dlg
	dlg.styler.model.items = dlg.Items
	dlg.ListBox.SetItemStyler(dlg.styler)
	dlg.ListBox.SetModel(dlg.styler.model)
	InitEvent(dlg.styler)
}

func OpenFolderDialog(callback func(dialog *MyDialog)) {
	RunMyDialog(nil, callback)
}
