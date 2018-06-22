package ii18n

import (
	"testing"
	"fmt"
)

func TestTranslate(t *testing.T) {
	config := map[string]Config{
		"app": Config{
			SourceNewFunc: NewJsonSource,
			OriginalLang:  "en-US",
			BasePath:      "./test",
			FileMap: map[string]string{
				"app":   "app.json",
				"error": "error.json",
			},
		},
	}
	NewI18N(config)
	res := T("app", "hello", nil, "zh-CN")
	fmt.Println(res)
}