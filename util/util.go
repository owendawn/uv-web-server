package util

import (
	"container/list"
	"os"
	"syscall"
)

func GetLogicalDrives() []string {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	GetLogicalDrives := kernel32.MustFindProc("GetLogicalDrives")
	n, _, _ := GetLogicalDrives.Call()
	print(n)
	s := "12"
	var drives_all = []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:", "N:", "O:", "P：", "Q：", "R：", "S：", "T：", "U：", "V：", "W：", "X：", "Y：", "Z："}
	temp := drives_all[0:len(s)]
	var d []string
	for i, v := range s {
		if v == 49 {
			l := len(s) - i - 1
			d = append(d, temp[l])
		}
	}
	var drives []string
	for i, v := range d {
		drives = append(drives[i:], append([]string{v}, drives[:i]...)...)
	}
	return drives
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetDiskList() *list.List {
	l := list.New()
	var drives_all = []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:", "N:", "O:", "P：", "Q：", "R：", "S：", "T：", "U：", "V：", "W：", "X：", "Y：", "Z："}
	for i := 0; i < len(drives_all); i++ {
		re, _ := PathExists(drives_all[i])
		if re {
			l.PushBack(drives_all[i])
		}
	}
	return l
}
