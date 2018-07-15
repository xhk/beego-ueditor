package routers

import (
	"beego-UEditor/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ueditor_controller", &controllers.UeditorController{})
}
