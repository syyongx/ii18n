package ii18n

import (
	"testing"
	"fmt"
)

func TestTranslate(t *testing.T) {
	config := map[string]Config{
		"app": Config{
			SourceNewFunc: NewJSONSource,
			OriginalLang:  "en-US",
			BasePath:      "./testdata",
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
