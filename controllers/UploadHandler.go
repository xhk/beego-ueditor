package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//----------------UploadHandler---------------
type UploadConfig struct {
	/// 文件命名规则
	PathFormat string

	/// 上传表单域名称
	UploadFieldName string

	/// 上传大小限制
	SizeLimit int

	/// 上传允许的文件格式
	AllowExtensions []string

	/// 文件是否以 Base64 的形式上传
	Base64 bool

	/// Base64 字符串所表示的文件名
	Base64Filename string
}

// UploadState
const (
	Success         = 0
	SizeLimitExceed = -1
	TypeNotAllow    = -2
	FileAccessError = -3
	NetworkError    = -4
	Unknown         = 1
)

type UploadResult struct {
	State          int
	Url            string
	OriginFileName string
	ErrorMessage   string
}

type UploadHandler struct {
	Handler
	Config UploadConfig
	Result UploadResult
}

func NewUploadHandler(c *UeditorController, config UploadConfig) *UploadHandler {
	var obj = new(UploadHandler)
	obj.context = c
	obj.Config = config
	//obj.Result = result

	return obj
}

func (this *UploadHandler) Process() {
	var uploadFileName = ""
	var uploadFileBytes []byte
	if this.Config.Base64 {
		uploadFileName = this.Config.Base64Filename
		uploadFileBytes, _ = base64.StdEncoding.DecodeString(this.context.GetString(this.Config.UploadFieldName))
	} else {
		f, file, err := this.context.GetFile(this.Config.UploadFieldName)
		if err != nil {
			this.Result.State = NetworkError
			this.WriteResult()
		}

		defer f.Close()
		uploadFileName = file.Filename
		if !this.CheckFileType(uploadFileName) {
			this.Result.State = TypeNotAllow
			this.WriteResult()
			return
		}

		if !this.CheckFileSize(int(file.Size)) {
			this.Result.State = SizeLimitExceed
			this.WriteResult()
			return
		}

		uploadFileBytes, err = ioutil.ReadAll(f)
		if err != nil {
			this.Result.State = FileAccessError
			this.WriteResult()
			return
		}
	}

	this.Result.OriginFileName = uploadFileName
	var savePath = PathFormat(uploadFileName, this.Config.PathFormat)
	var localPath = "static/" + savePath
	fmt.Printf("savePath:%s localPath:%s\n", savePath, localPath)
	localDir := localPath[:strings.LastIndex(localPath, "/")]
	exist, err := PathExists(localDir)
	if !exist {
		err = os.Mkdir(localDir, os.ModePerm)
		if err != nil {
			this.Result.State = FileAccessError
			this.WriteResult()
			return
		}
	}

	err = ioutil.WriteFile(localPath, uploadFileBytes, 0644)
	if err != nil {
		this.Result.State = FileAccessError
		this.WriteResult()
		return
	}

	this.Result.Url = savePath
	this.Result.State = Success
	this.WriteResult()
}

func (this *UploadHandler) WriteResult() {
	result := make(map[string]string)
	result["state"] = this.GetStateMessage(this.Result.State)
	result["url"] = this.Result.Url
	result["title"] = this.Result.OriginFileName
	result["original"] = this.Result.OriginFileName
	result["error"] = this.Result.ErrorMessage

	this.WriteJson(result)
}

func (this *UploadHandler) GetStateMessage(state int) string {
	switch state {
	case Success:
		return "SUCCESS"
	case FileAccessError:
		return "文件访问出错，请检查写入权限"
	case SizeLimitExceed:
		return "文件大小超出服务器限制"
	case TypeNotAllow:
		return "不允许的文件格式"
	case NetworkError:
		return "网络错误"
	}
	return "未知错误"
}

func (this *UploadHandler) CheckFileType(fileName string) bool {
	//fmt.Printf("filename:%s \n", fileName)
	var fileExternsion = strings.ToLower(path.Ext(fileName))
	//fileExternsion = strings.Replace(fileExternsion, ".", "", 1)
	//fmt.Printf("fileExternsion:%s \n", fileExternsion)
	//fmt.Printf("AllowExtensions:%s \n", this.Config.AllowExtensions)

	var exist = false
	for _, v := range this.Config.AllowExtensions {
		if strings.ToLower(v) == fileExternsion {
			exist = true
			break
		}
	}
	return exist
}

func (this *UploadHandler) CheckFileSize(size int) bool {
	fmt.Printf("file size:%d, limit size:%d\n", size, this.Config.SizeLimit)
	return size < this.Config.SizeLimit
}
