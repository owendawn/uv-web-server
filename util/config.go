package util

import (
	"github.com/lxn/walk"
	"os"
)

func NewSystemIcon() walk.Image {
	dir, _ := os.Getwd()
	icon, _ := walk.NewImageFromFileForDPI(dir+"/asserts/uv.png", 500)
	return icon
}
func NewSystemPath() string {
	dir, _ := os.Getwd()
	return dir
}

func NewUserHomePath() string {
	p, _ := os.UserHomeDir()
	return p
}
