package main

import (
	"context"
	"fmt"
	"kkparse/parser"
	"regexp"
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
