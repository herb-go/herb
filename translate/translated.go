package translate

//Translated translated message
type Translated struct {
	translateable Translateable
	lang          string
}

//String return translated string
func (t *Translated) String() string {
	return t.translateable.Translate(t.lang)
}

// NewTranslated create new translated message
func NewTranslated(translateable Translateable, lang string) *Translated {
	return &Translated{
		translateable: translateable,
		lang:          lang,
	}
}
