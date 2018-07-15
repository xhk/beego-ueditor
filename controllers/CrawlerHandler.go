package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CrawlerHandler struct {
	Handler
	Sources  []string
	Crawlers []Crawler
}

func NewCrawlerHandler(c *UeditorController) *CrawlerHandler {
	var this = new(CrawlerHandler)
	this.context = c

	return this
}

func (this *CrawlerHandler) Process() {
	this.Sources = this.context.GetStrings("source[]")
	if len(this.Sources) == 0 {
		this.WriteJson(map[string]string{"state": "参数错误：没有指定抓取源"})
		return
	}

	cfg := NewConfig()
	pathFormat := cfg.GetString("catcherPathFormat")
	for _, source := range this.Sources {
		crawler := Crawler{source, "", "", pathFormat}
		crawler.Fetch()
		this.Crawlers = append(this.Crawlers, crawler)
	}
}

type Crawler struct {
	SourceUrl  string
	ServerUrl  string
	State      string
	PathFormat string
}

func (this *Crawler) Fetch() {
	rsp, err := http.Get(this.SourceUrl)
	if err != nil {
		this.State = "Url return " + err.Error()
		return
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		this.State = "Url return " + strconv.Itoa(rsp.StatusCode)
		return
	}

	contentType := rsp.Header.Get("Content-Type")
	if !(contentType != "" && strings.ContainsAny(contentType, "image")) {
		this.State = "Url is not an image"
		return
	}

	_, fileName := filepath.Split(this.SourceUrl)
	this.ServerUrl = PathFormat(fileName, this.PathFormat)
	savePath := "static/" + this.ServerUrl
	this.ServerUrl = "static" + this.ServerUrl
	localPath := savePath
	localDir := localPath[:strings.LastIndex(localPath, "\\")]
	exist, err := PathExists(localDir)
	if !exist {
		err = os.Mkdir(localDir, os.ModePerm)
		if err != nil {
			this.State = "Save file failed"
			return
		}
	}

	imgBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		this.State = "Read image content failed"
		return
	}

	if err = ioutil.WriteFile(savePath, imgBytes, 064); err != nil {
		this.State = "Save file failed"
		return
	}
}
