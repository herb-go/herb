package ui

//Translated translated message
type Translated struct {
	translateable Translatable
	lang          string
}

//Label return translated string
func (t *Translated) Label() string {
	return t.translateable.Translate(t.lang)
}

// NewTranslated create new translated message
func NewTranslated(translateable Translatable, lang string) *Translated {
	return &Translated{
		translateable: translateable,
		lang:          lang,
	}
}
