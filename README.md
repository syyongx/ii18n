# ii18n
Go i18n library.

## Download & Install
```shell
go get github.com/syyongx/ii18n
```

## Quick Start
```go
import github.com/syyongx/ii18n

func main() {
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
    message := T("app", "hello", nil, "zh-CN")
}
```

## LICENSE
ii18n source code is licensed under the [MIT](https://github.com/syyongx/ii18n/blob/master/LICENSE) Licence.
