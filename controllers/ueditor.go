package controllers

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type UeditorController struct {
	beego.Controller
}

type IHandler interface {
	WriteJson(data interface{})
	Process()
}

func PathFormat(originFileName string, pathFormat string) string {
	if pathFormat == "" {
		pathFormat = "{filename}{rand:6}"
	}

	// 去掉路径之不合法的字符
	invalidPattern, _ := regexp.Compile("[\\\\\\/\\:\\*\\?\\042\\<\\>\\|]")
	originFileName = invalidPattern.ReplaceAllString(originFileName, "")
	var extension = path.Ext(originFileName)
	var filename = originFileName[:strings.LastIndex(originFileName, ".")]
	pathFormat = strings.Replace(pathFormat, "{filename}", filename, -1)

	re, _ := regexp.Compile("\\{rand(\\:?)(\\d+)\\}")
	pathFormat = re.ReplaceAllStringFunc(pathFormat, func(data string) string {
		m := re.FindStringSubmatch(data)
		var digit = 6
		if len(m) > 2 {
			//fmt.Println(m[2])
			digit, _ = strconv.Atoi(m[2])
		}
		var num int64 = 1
		for i := 0; i < digit; i++ {
			num *= 10
		}

		// 用时间来代替随机数，反正随机数一般也是拿时间戳做随机种子的
		return strconv.Itoa(int(time.Now().Unix() % num))
	})

	now := time.Now()
	pathFormat = strings.Replace(pathFormat, "{time}", strconv.Itoa(int(now.Unix())), -1)
	pathFormat = strings.Replace(pathFormat, "{yyyy}", strconv.Itoa(now.Year()), -1)
	pathFormat = strings.Replace(pathFormat, "{yy}", fmt.Sprintf("%02d", now.Year()%100), -1)
	pathFormat = strings.Replace(pathFormat, "{mm}", fmt.Sprintf("%02d", now.Month()%100), -1)
	pathFormat = strings.Replace(pathFormat, "{dd}", fmt.Sprintf("%02d", now.Day()%100), -1)
	pathFormat = strings.Replace(pathFormat, "{hh}", fmt.Sprintf("%02d", now.Hour()%100), -1)
	pathFormat = strings.Replace(pathFormat, "{ii}", fmt.Sprintf("%02d", now.Minute()%100), -1)
	pathFormat = strings.Replace(pathFormat, "{ss}", fmt.Sprintf("%02d", now.Second()%100), -1)

	return pathFormat + extension
}

// 判断文件夹是否存在
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

func (this *UeditorController) Get() {
	var action IHandler
	context := this
	Config := NewConfig()
	switch this.GetString("action") {
	default:
		action = NewNotSupportedHandler(this)
	case "config":
		action = (NewConfigHandler(this))
	case "uploadimage":
		uploadCfg := UploadConfig{}

		uploadCfg.AllowExtensions = Config.GetStringList("imageAllowFiles")
		uploadCfg.PathFormat = Config.GetString("imagePathFormat")
		uploadCfg.SizeLimit = Config.GetInt("imageMaxSize")
		uploadCfg.UploadFieldName = Config.GetString("imageFieldName")
		action = NewUploadHandler(this, uploadCfg)
	case "uploadscrawl":
		uploadCfg := UploadConfig{}
		uploadCfg.AllowExtensions = []string{".png"}
		uploadCfg.PathFormat = Config.GetString("scrawlPathFormat")
		uploadCfg.SizeLimit = Config.GetInt("scrawlMaxSize")
		uploadCfg.UploadFieldName = Config.GetString("scrawlFieldName")
		uploadCfg.Base64 = true
		uploadCfg.Base64Filename = "scrawl.png"
		action = NewUploadHandler(context, uploadCfg)

	case "uploadvideo":
		uploadCfg := UploadConfig{}
		uploadCfg.AllowExtensions = Config.GetStringList("videoAllowFiles")
		uploadCfg.PathFormat = Config.GetString("videoPathFormat")
		uploadCfg.SizeLimit = Config.GetInt("videoMaxSize")
		uploadCfg.UploadFieldName = Config.GetString("videoFieldName")
		action = NewUploadHandler(context, uploadCfg)
		break
	case "uploadfile":
		uploadCfg := UploadConfig{}
		uploadCfg.AllowExtensions = Config.GetStringList("fileAllowFiles")
		uploadCfg.PathFormat = Config.GetString("filePathFormat")
		uploadCfg.SizeLimit = Config.GetInt("fileMaxSize")
		uploadCfg.UploadFieldName = Config.GetString("fileFieldName")
		action = NewUploadHandler(context, uploadCfg)
	case "listimage":
		action = NewListFileManager(context, Config.GetString("imageManagerListPath"), Config.GetStringList("imageManagerAllowFiles"))
	case "listfile":
		action = NewListFileManager(context, Config.GetString("fileManagerListPath"), Config.GetStringList("fileManagerAllowFiles"))
	case "catchimage":
		action = NewCrawlerHandler(context)
	}

	action.Process()
}

func (this *UeditorController) Post() {
	this.Get()
}
