package controllers

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

type Config struct {
	Items map[string]interface{}
}

func NewConfig() *Config {
	cfg := new(Config)
	cfgFile := "conf/ueditor/config.json"
	content, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return cfg
	}

	// 需要去除其中的注释，才能够正确解释
	re, err := regexp.Compile("/\\*(?:\\s|.)*?\\*\\/")
	content = re.ReplaceAll(content, []byte(""))
	json.Unmarshal(content, &cfg.Items)

	return cfg
}

func (this *Config) GetValue(key string) interface{} {
	return this.Items[key]
}

func (this *Config) GetStringList(key string) []string {
	var ret []string
	val := this.Items[key].([]interface{})
	for _, v := range val {
		ret = append(ret, v.(string))
	}
	return ret
	//return []string(this.Items[key])
}

func (this *Config) GetString(key string) string {
	val := this.Items[key].(string)
	return val
}

func (this *Config) GetInt(key string) int {
	val := this.Items[key].(float64)
	return int(val)
}
