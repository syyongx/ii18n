package ii18n

import (
	"errors"
	"strings"
	"sync"
)

// TMsgs type messages
type TMsgs map[string]string

// Source interface
type Source interface {
	Translate(category string, message string, lang string) (string, error)
	TranslateMsg(category string, message string, lang string) (string, error)
	GetMsgFilePath(category string, lang string) string
	LoadMsgs(category string, lang string) (TMsgs, error)
	LoadFallbackMsgs(category string, fallbackLang string, msgs TMsgs, originalMsgFile string) (TMsgs, error)
}

// MessageSource MessageSource
type MessageSource struct {
	// string the language that the original messages are in
	OriginalLang     string
	ForceTranslation bool
	BasePath         string
	FileMap          map[string]string
	fileSuffix       string
	loadFunc         func(filename string) (TMsgs, error)
	messages         map[string]TMsgs
	mutex            sync.RWMutex
}

// Translate translate
func (ms *MessageSource) Translate(category string, message string, lang string) (string, error) {
	if ms.ForceTranslation || lang != ms.OriginalLang {
		return ms.TranslateMsg(category, message, lang)
	}
	return "", nil
}

// TranslateMsg translate message
func (ms *MessageSource) TranslateMsg(category string, message string, lang string) (string, error) {
	cates := strings.Split(category, ".")
	key := cates[0] + "/" + lang + "/" + cates[1]

	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	if _, ok := ms.messages[key]; !ok {
		val, err := ms.LoadMsgs(category, lang)
		if err != nil {
			return "", err
		}
		ms.messages[key] = val
	}
	if msg, ok := ms.messages[key][message]; ok && msg != "" {
		return msg, nil
	}

	ms.messages[key] = TMsgs{message: ""}
	return "", nil
}

// GetMsgFilePath Get messages file path.
func (ms *MessageSource) GetMsgFilePath(category string, lang string) string {
	suffix := strings.Split(category, ".")[1]
	path := ms.BasePath + "/" + lang + "/"
	if v, ok := ms.FileMap[suffix]; !ok {
		path += v
	} else {
		path += strings.Replace(suffix, "\\", "/", -1)
		if ms.fileSuffix != "" {
			path += "." + ms.fileSuffix
		}
	}
	return path
}

// LoadMsgs Loads the message translation for the specified $language and $category.
// If translation for specific locale code such as `en-US` isn't found it
// tries more generic `en`. When both are present, the `en-US` messages will be merged
// over `en`. See [[loadFallbackTMsgs]] for details.
// If the lang is less specific than [[originalLang]], the method will try to
// load the messages for [[originalLang]]. For example: [[originalLang]] is `en-GB`,
// language is `en`. The method will load the messages for `en` and merge them over `en-GB`.
func (ms *MessageSource) LoadMsgs(category string, lang string) (TMsgs, error) {
	msgFile := ms.GetMsgFilePath(category, lang)
	msgs, err := ms.loadFunc(msgFile)
	if err != nil {
		return nil, err
	}
	fbLang := lang[0:2]
	fbOriginalLang := ms.OriginalLang[0:2]
	if lang != fbLang {
		msgs, err = ms.LoadFallbackMsgs(category, fbLang, msgs, msgFile)
	} else if lang == fbOriginalLang {
		msgs, err = ms.LoadFallbackMsgs(category, ms.OriginalLang, msgs, msgFile)
	} else {
		if msgs == nil {
			return nil, errors.New("the message file for category " + category + " does not exist: " + msgFile)
		}
	}
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// LoadFallbackMsgs Loads the message translation for the specified $language and $category.
// If translation for specific locale code such as `en-US` isn't found it
// tries more generic `en`. When both are present, the `en-US` messages will be merged
func (ms *MessageSource) LoadFallbackMsgs(category string, fallbackLang string, msgs TMsgs, originalMsgFile string) (TMsgs, error) {
	fbMsgFile := ms.GetMsgFilePath(category, fallbackLang)
	fbMsgs, _ := ms.loadFunc(fbMsgFile)
	if msgs == nil && fbMsgs == nil &&
		fallbackLang != ms.OriginalLang &&
		fallbackLang != ms.OriginalLang[0:2] {
		return nil, errors.New("The message file for category " + category + " does not exist: " + originalMsgFile + " Fallback file does not exist as well: " + fbMsgFile)
	} else if msgs == nil {
		return fbMsgs, nil
	} else if fbMsgs != nil {
		ms.mutex.Lock()
		defer ms.mutex.Unlock()

		for key, val := range fbMsgs {
			v, ok := msgs[key]
			if val != "" && (!ok || v == "") {
				msgs[key] = val
			}
		}
	}

	return msgs, nil
}

// LoadMsgsFromFile Get messages file path.
func LoadMsgsFromFile(filename string) (TMsgs, error) {
	return nil, nil
}
