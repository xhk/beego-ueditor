package controllers

type Handler struct {
	context *UeditorController
}

func (this *Handler) WriteJson(data interface{}) {
	var jsonpCallback = this.context.GetString("callback")
	if jsonpCallback == "" {
		this.context.Data["json"] = data
		this.context.ServeJSON()
	} else {
		this.context.Data["jsonp"] = data
		this.context.ServeJSON()
	}
}

func (this *Handler) Process() {

}
