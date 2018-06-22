package ii18n

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type JsonSource struct {
	MessageSource
}

func NewJsonSource(conf *Config) Source {
	s := &JsonSource{}
	s.OriginalLang = conf.OriginalLang
	s.BasePath = conf.BasePath
	s.ForceTranslation = conf.ForceTranslation
	s.FileMap = conf.FileMap
	s.messages = make(map[string]TMsgs)
	s.fileSuffix = "json"
	s.loadFunc = loadMsgsFromJsonFile

	return s
}

// Get messages file path.
func (js *JsonSource) GetMsgFilePath(category string, lang string) string {
	suffix := strings.Split(category, ".")[1]
	path := js.BasePath + "/" + lang + "/"
	if v, ok := js.FileMap[suffix]; !ok {
		path += v
	} else {
		path += strings.Replace(suffix, "\\", "/", -1)
	}
	return path
}

// Get messages file path.
func loadMsgsFromJsonFile(filename string) (TMsgs, error) {
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
