package ii18n

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// JSONSource JSONSource
type JSONSource struct {
	MessageSource
}

// NewJSONSource New JSONSource
func NewJSONSource(conf *Config) Source {
	s := &JSONSource{}
	s.OriginalLang = conf.OriginalLang
	s.BasePath = conf.BasePath
	s.ForceTranslation = conf.ForceTranslation
	s.FileMap = conf.FileMap
	s.messages = make(map[string]TMsgs)
	s.fileSuffix = "json"
	s.loadFunc = loadMsgsFromJSONFile

	return s
}

// GetMsgFilePath Get messages file path.
func (js *JSONSource) GetMsgFilePath(category string, lang string) string {
	suffix := strings.Split(category, ".")[1]
	path := js.BasePath + "/" + lang + "/"
	if v, ok := js.FileMap[suffix]; !ok {
		path += v
	} else {
		path += strings.Replace(suffix, "\\", "/", -1)
	}
	return path
}

// loadMsgsFromJSONFile Get messages file path.
func loadMsgsFromJSONFile(filename string) (TMsgs, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var msgs TMsgs
	e := json.Unmarshal(data, &msgs)
	if e != nil {
		return nil, e
	}

	return msgs, nil
}
