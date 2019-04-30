package ii18n

import (
	"fmt"
	"testing"
)

func TestTranslate(t *testing.T) {
	config := map[string]Config{
		"app": {
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
