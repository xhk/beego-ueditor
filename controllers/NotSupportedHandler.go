package controllers

type NotSupportedHandler struct {
	Handler
}

func NewNotSupportedHandler(c *UeditorController) *NotSupportedHandler {
	this := new(NotSupportedHandler)
	this.context = c

	return this
}

func (this *NotSupportedHandler) Process() {
	this.WriteJson(map[string]string{"state": "action 参数为空或者 action 不被支持。"})
}
