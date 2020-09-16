package main

import (
	"context"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	myhttp "uv-web-server/custom"
	"uv-web-server/util"
	"uv-web-server/walk/dialog"
)

func main3() {
	const serverName = "UV轻量级Web服务器"
	const configPath = "./config.json"
	var mw *walk.MainWindow
	var server *http.Server
	var logger = log.New(os.Stdout, "http: ", log.LstdFlags)
	//d:=walk.NewDialog(path)

	win := MainWindow{
		AssignTo: &mw,
		Title:    serverName,
		MinSize:  Size{Width: 600, Height: 200},
		Size:     Size{Width: 600, Height: 200},
		MaxSize:  Size{Width: 600},
		Layout:   VBox{},
		Icon:     util.NewSystemIcon(),
	}
	var path *walk.LineEdit
	var port *walk.LineEdit
	var runBtn *walk.PushButton
	var stopBtn *walk.PushButton

	//var inTE, outTE *walk.TextEdit
	ws := []Widget{
		HSplitter{
			MaxSize: Size{
				Height: 25,
			},
			Children: []Widget{
				Label{Text: "资源路径：", StretchFactor: 1},
				LineEdit{AssignTo: &path, StretchFactor: 8},
				PushButton{Text: "选择", StretchFactor: 1,
					OnClicked: func() {
						path.SetText("ok")
						//walk.MsgBox(nil, "Open", "Pretend to open a file...", walk.MsgBoxIconInformation)
						dialog.OpenFolderDialog(func(dialog *dialog.MyDialog) {
							println(dialog.Path.Text())
							path.SetText(strings.Replace(dialog.Path.Text(), "\\", "/", -1))
						})
					},
				},
			},
		},
		HSplitter{
			MaxSize: Size{
				Height: 25,
			},
			Children: []Widget{
				Label{Text: "端口：", StretchFactor: 1},
				LineEdit{AssignTo: &port, StretchFactor: 9},
			},
		},
		PushButton{
			Text:     "启动",
			AssignTo: &runBtn,
			Background: SolidColorBrush{
				Color: walk.RGB(0, 0xff, 0),
			},
			OnClicked: func() {
				InitConfig(path, port)
				if path.Text() == "" {
					ShowNotifyDialog("请选择资源路径")
				} else if port.Text() == "" {
					ShowNotifyDialog("请填写端口号")
				} else {
					SaveConfig(configPath, Config{
						Root: path.Text(),
						Port: port.Text(),
					})
					stopBtn.SetVisible(true)
					runBtn.SetVisible(false)
					server = newWebServer(path.Text(), port.Text(), logger)
					go server.ListenAndServe()
					println("启动服务")
				}
			},
		},
		PushButton{
			Text:     "关闭",
			Visible:  false,
			AssignTo: &stopBtn,
			OnClicked: func() {
				stopBtn.SetVisible(false)
				runBtn.SetVisible(true)
				runBtn.SetEnabled(false)
				go shutdown(server, logger, runBtn)
				println("关闭服务")
			},
		},
	}
	win.Children = ws

	go (func() {
		time.Sleep(1 * time.Second)
		InitConfig(path, port)
		//托盘图标文件
		ni, err := walk.NewNotifyIcon(mw)
		if err != nil {
			log.Fatal(err)
		}
		if err := ni.SetIcon(util.NewSystemIcon()); err != nil {
			log.Fatal(err)
		}
		if err := ni.SetToolTip(serverName); err != nil {
			log.Fatal(err)
		}
		ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
			if button != walk.LeftButton {
				return
			}
			mw.SetFocus()
		})
		exitAction := walk.NewAction()
		if err := exitAction.SetText("右键icon的菜单按钮"); err != nil {
			log.Fatal(err)
		}
		//Exit 实现的功能
		exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
		if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
			log.Fatal(err)
		}
		if err := ni.SetVisible(true); err != nil {
			log.Fatal(err)
		}
		if err := ni.ShowInfo("Welcome", "Welcome to use my tool -- "+serverName); err != nil {
			log.Fatal(err)
		}
	})()

	win.Run()

}

const configPath = "./config.json"

func InitConfig(path *walk.LineEdit, port *walk.LineEdit) {
	cfg := LoadConfig(configPath)
	if cfg != nil {
		if path.Text() == "" {
			path.SetText(cfg.Root)
		}
		if port.Text() == "" {
			port.SetText(cfg.Port)
		}
	}
}

//初始化 server
func newWebServer(path string, port string, logger *log.Logger) *http.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//w.WriteHeader(http.StatusOK)
		//http.FileServer(http.Dir(path)).ServeHTTP(w,r)
		myhttp.MyFileServer(http.Dir(path)).ServeHTTP(w, r)
	})
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	return server
}

//关闭 server
//quit: 接收关闭信号
//done: 发出已经关闭信号
func shutdown(server *http.Server, logger *log.Logger, btn *walk.PushButton) {
	btn.SetText("服务关闭中……")
	//等待接收到退出信号：
	logger.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Fatalf("Could not gracefully shutdown the server: %v \n", err)
	}

	//do Something ：
	fmt.Println("do something start ..... ", time.Now())
	time.Sleep(5 * time.Second)
	fmt.Println("do something end ..... ", time.Now())
	btn.SetText("启动")
	btn.SetEnabled(true)
}

func ShowNotifyDialog(str string) {
	dlg, _ := walk.NewDialog(nil)
	dlg.SetName("Dialog")
	dlg.SetTitle("提示")
	dlg.SetLayout(walk.NewVBoxLayout())
	dlg.SetClientSize(walk.Size{700, 300})
	dlg.SetSize(walk.Size{700, 500})
	dlg.SetIcon(util.NewSystemIcon())
	label, _ := walk.NewTextLabel(dlg)
	label.SetText(str)
	dlg.Form().Children().Add(label)
	closeBtn, _ := walk.NewPushButton(dlg)
	closeBtn.SetText("确定")
	closeBtn.Clicked().Attach(func() {
		dlg.Close(0)
	})
	dlg.Form().Children().Add(closeBtn)
	dlg.Run()
}
