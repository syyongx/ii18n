package ii18n

import (
	"regexp"
	"strings"
	"sync"
)

// DefaultOriginalLang default original language
var DefaultOriginalLang = "en-US"

// Translator translator
var Translator *I18N

// T translate.
// 1. T('common', 'hot', [], 'zh-CN') // default app.common
// 2. T('app.common', 'hot', [], 'zh-CN') // result same to 1.
// 3. T('msg.a', 'hello', ['{foo}' => 'bar', '{key}' => 'val'] 'ja-JP')
func T(category string, message string, params map[string]string, lang string) string {
	if strings.Index(category, ".") == -1 {
		category = "app." + category
	}
	return Translator.translate(category, message, params, lang)
}

// Config config
type Config struct {
	SourceNewFunc    func(*Config) Source
	OriginalLang     string
	ForceTranslation bool
	BasePath         string
	FileMap          map[string]string
	source           Source
}

// I18N i18n
type I18N struct {
	Translations map[string]*Config
	formatter    Formatter
	mutex        sync.RWMutex
}

// NewI18N returns an instance of I18N.
func NewI18N(config map[string]Config) *I18N {
	Translator = &I18N{
		Translations: make(map[string]*Config),
	}
	for key, conf := range config {
		if conf.SourceNewFunc == nil {
			panic("Config SourceNewFunc is illegal")
		}
		if conf.OriginalLang == "" {
			conf.OriginalLang = DefaultOriginalLang
		}
		if len(conf.OriginalLang) < 2 {
			panic("Config OriginalLang length cannot be less than 2")
		}
		if conf.BasePath == "" {
			panic("Config BasePath is illegal")
		}
		if conf.FileMap == nil {
			panic("Config FileMap is illegal")
		}
		if _, ok := Translator.Translations[key]; !ok {
			Translator.Translations[key] = &conf
		}
	}
	return Translator
}

// translate
func (i *I18N) translate(category string, message string, params map[string]string, lang string) string {
	s, ol := i.getSource(category)
	translation, err := s.Translate(category, message, lang)
	if err != nil || translation == "" {
		return i.format(message, params, ol)
	}
	return i.format(translation, params, lang)
}

func (i *I18N) format(message string, params map[string]string, lang string) string {
	if params == nil {
		return message
	}
	if ok, _ := regexp.MatchString(`~{\s*[\d\w]+\s*,~u`, message); ok {
		result, err := i.formatter.format(message, params, lang)
		if err != nil {
			return message
		}
		return result
	}
	oldnew := make([]string, len(params)*2)
	for name, val := range params {
		oldnew = append(oldnew, "{"+name+"}", val)
	}
	return strings.NewReplacer(oldnew...).Replace(message)
}

// getFormatter Get the the message formatter.
func (i *I18N) getFormatter(category string) Formatter {
	return i.formatter
}

// getSource Get the message source for the given category.
func (i *I18N) getSource(category string) (Source, string) {
	prefix := strings.Split(category, ".")[0]
	if val, ok := i.Translations[prefix]; ok {
		i.mutex.Lock()
		defer i.mutex.Unlock()
		if val.source == nil {
			i.Translations[prefix].source = i.Translations[prefix].SourceNewFunc(i.Translations[prefix])
		}
		return i.Translations[prefix].source, i.Translations[prefix].OriginalLang
	}
	panic("Unable to locate message source for category " + category + ".")
}
