package controllers

//---------------ConfigHandler----------------
type ConfigHandler struct {
	Handler
}

func NewConfigHandler(c *UeditorController) *ConfigHandler {
	var obj = new(ConfigHandler)
	obj.context = c
	return obj
}

func (this *ConfigHandler) Process() {
	cfg := NewConfig()
	this.WriteJson(cfg.Items)
}
