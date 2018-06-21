package ii18n

import (
	"strings"
	"regexp"
)

var DefaultSourceLang = "en-US"
var Translator *I18N

// translate.
// 1) T('common', 'hot', [], 'zh_cn'); // default app.common
// 2) T('app.common', 'hot', [], 'zh_cn'); // result same to 1)
// 3) T('msg.a', 'hello', ['{foo}' => 'bar', '{key}' => 'val'], 'ja_jp');
func T(category string, message string, params map[string]string, lang string) string {
	if strings.Index(category, ".") == -1 {
		category = "app." + category
	}
	return Translator.translate(category, message, params, lang)
}

type Config struct {
	SourceLang       string
	BasePath         string
	FileMap          map[string]string
	ForceTranslation bool
	// message source
	source Source
}

// I18N
type I18N struct {
	Translations map[string]*Config
	formatter    Formatter
}

// New I18N
func NewI18N(config map[string]Config) *I18N {
	Translator = &I18N{
		Translations: make(map[string]*Config),
	}
	for key, conf := range config {
		if conf.SourceLang == "" {
			conf.SourceLang = DefaultSourceLang
		}
		if conf.BasePath == "" {
			panic("BasePath config invalid")
		}
		if conf.FileMap == nil {
			panic("FileMap config invalid")
		}
		if _, ok := Translator.Translations[key]; !ok {
			Translator.Translations[key] = &conf
		}
	}
	return Translator
}

// translate
func (i *I18N) translate(category string, message string, params map[string]string, lang string) string {
	s := i.getSource(category)
	translation, err := s.Translate(category, message, lang)
	if err != nil || translation == "" {
		return i.format(message, params, s.(*JsonSource).SourceLang)
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
func (i *I18N) getSource(category string) Source {
	prefix := strings.Split(category, ".")[0]
	if val, ok := i.Translations[prefix]; ok {
		if val.source == nil {
			i.Translations[prefix].source = NewJsonSource(i.Translations[prefix])
		}
		return i.Translations[prefix].source
	}
	panic("Unable to locate message source for category " + category + ".")
}
