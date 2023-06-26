package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"kkparse/parser"
	"os"
	"regexp"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Parse(shareUrl string) (*parser.VideoParseInfo, error) {
	// 解析
	urlReg := regexp.MustCompile(`http[s]?:\/\/[\w.-]+[\w\/-]*[\w.-]*\??[\w=&:\-\+\%]*[/]*`)
	videoShareUrl := urlReg.FindString(shareUrl)

	parseRes, err := parser.ParseVideoShareUrl(videoShareUrl)
	if err != nil {
		return nil, err
	}
	return parseRes, nil
}

func (a *App) Download(videoUrl string) error {
	// 判断是否存在默认保存的目录
	_, err := os.Stat("./download.txt")
	if err != nil {
		directory, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
			Title: "请先选择默认保存目录",
		})
		if err != nil {
			return err
		}
		ioutil.WriteFile("./download.txt", []byte(directory), 0644)
		return err
	}

	return nil
}
