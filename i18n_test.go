package ii18n

import (
	"testing"
	"fmt"
)

func TestTranslate(t *testing.T) {
	config := map[string]Config{
		"app": Config{
			SourceLang: "en-Us",
			BasePath:   "./test",
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
