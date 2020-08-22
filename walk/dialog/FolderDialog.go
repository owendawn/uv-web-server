// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package dialog

import (
	"fmt"
	"github.com/lxn/walk"
	"log"
)

type myDialogUI struct {
	hSpitter  *walk.Splitter
	normalBtn *walk.PushButton
	errBtn    *walk.PushButton
}

type MyDialog struct {
	*walk.Dialog
	ui       myDialogUI
	callback func(dialog *MyDialog)
}

func (dlg *MyDialog) setState(state walk.PIState) {
	if err := dlg.ProgressIndicator().SetState(state); err != nil {
		log.Print(err)
	}
}

func RunMyDialog(owner walk.Form, callback func(dialog *MyDialog)) (int, error) {
	dlg := new(MyDialog)
	dlg.Dialog, _ = walk.NewDialog(owner)

	succeeded := false
	defer func() {
		if !succeeded {
			dlg.Dispose()
		}
	}()

	dlg.callback = callback
	dlg.SetLayout(walk.NewVBoxLayout())
	dlg.SetName("Dialog")
	dlg.SetClientSize(walk.Size{598, 300})
	dlg.SetTitle("hello1")

	// normalBtn
	normalBtn, _ := walk.NewPushButton(dlg)
	normalBtn.SetName("normalBtn")
	normalBtn.SetBounds(walk.Rectangle{40, 120, 161, 23})
	normalBtn.SetText(`Normal`)
	normalBtn.SetMinMaxSize(walk.Size{0, 0}, walk.Size{161, 16777215})
	normalBtn.Clicked().Attach(func() {
		fmt.Println("SetState normal")
		dlg.setState(walk.PINormal)
	})
	dlg.Children().Add(normalBtn)
	// errBtn
	errBtn, _ := walk.NewPushButton(dlg)
	errBtn.SetName("errBtn")
	errBtn.SetBounds(walk.Rectangle{40, 150, 161, 23})
	errBtn.SetText(`Error`)
	errBtn.Clicked().Attach(func() {
		fmt.Println("SetState error")
		dlg.setState(walk.PIError)
	})
	dlg.Children().Add(errBtn)

	lable, _ := walk.NewLabel(dlg)
	lable.SetText("hhhhhh")
	dlg.Children().Add(lable)

	dlg.callback(dlg)

	succeeded = true
	return dlg.Run(), nil
}

func OpenFolderDialog(callback func(dialog *MyDialog)) {
	RunMyDialog(nil, callback)
}
