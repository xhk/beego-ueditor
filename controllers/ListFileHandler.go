package controllers

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// 获取目录下所有文件
func getAllFilelist(path string) []string {
	var fileList = make([]string, 0, 100)
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		//println(path)
		fileList = append(fileList, f.Name())
		return nil
	})
	// if err != nil {
	//     fmt.Printf("filepath.Walk() returned %v\n", err)
	// }

	return fileList
}

type ByName []string

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i] < a[j] }

const (
	LfmSuccess       = 0
	LfmInvalidParam  = 1
	LfmAuthorizError = 2
	LfmIOError       = 3
	LfmPathNotFound  = 4
)

type ListFileManager struct {
	Handler
	Start            int
	Size             int
	Total            int
	State            int
	PathToList       string
	FileList         []string
	SearchExtensions []string
}

func NewListFileManager(c *UeditorController, pathToList string, searchExtensions []string) *ListFileManager {
	obj := new(ListFileManager)
	obj.context = c
	obj.PathToList = pathToList
	for i, v := range searchExtensions {
		searchExtensions[i] = strings.ToLower(v)
	}
	obj.SearchExtensions = searchExtensions
	return obj
}

func (this *ListFileManager) Process() {
	this.Start, _ = this.context.GetInt("start")
	this.Size, _ = this.context.GetInt("size")

	var localPath = this.PathToList
	var buildingList = getAllFilelist(localPath)
	sort.Sort(ByName(buildingList))
	this.Total = len(buildingList)
	this.FileList = buildingList[this.Start : this.Start+this.Size]
	this.State = LfmSuccess
	this.WriteResult()
}

func (this *ListFileManager) WriteResult() {
	this.WriteJson(map[string]interface{}{
		"state": this.GetStateString(),
		"list":  this.FileList,
		"start": this.Start,
		"size":  this.Size,
		"total": this.Total,
	})
}

func (this *ListFileManager) GetStateString() string {
	switch this.State {
	case LfmSuccess:
		return "SUCCESS"
	case LfmInvalidParam:
		return "参数不正确"
	case LfmPathNotFound:
		return "路径不存在"
	case LfmAuthorizError:
		return "文件系统权限不足"
	case LfmIOError:
		return "文件系统读取错误"
	}

	return "未知错误"
}
