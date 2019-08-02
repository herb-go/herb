package translate

//TranslationLanguage translateable interface
type TranslationLanguage interface {
	Lang() string
	SetLang(string)
}

//Language translation language struct
type Language struct {
	lang string
}

//Lang get  language
func (l *Language) Lang() string {
	return l.lang
}

//SetLang set language
func (l *Language) SetLang(lang string) {
	l.lang = lang
}
