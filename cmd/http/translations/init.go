package translations

import (
	"encoding/json"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var transJson = `{}`

type Trans map[string]map[string]string

func init() {
	var trans Trans
	err := json.Unmarshal([]byte(transJson), &trans)
	if err != nil {
		panic(err)
	}

	var langTag language.Tag
	for key, items := range trans {
		for lang, msg := range items {
			err = langTag.UnmarshalText([]byte(lang))
			if err != nil {
				panic(err)
			}
			err = message.SetString(langTag, key, msg)
			if err != nil {
				panic(err)
			}
		}
	}
}
