package translate

//Translateable interface which can be translated
type Translateable interface {
	Translate(lang string) string
}

//Message message which should be translated
type Message struct {
	Module string
	Text   string
}

//Translate translate message with default translations.
//if lang if empty,Lang will be used.
func (m *Message) Translate(lang string) string {
	if lang == "" {
		lang = Lang
	}
	return m.TranslateWith(DefaultTranslations, lang)
}

//TranslateWith translate message with given translations and language
func (m *Message) TranslateWith(t *Translations, lang string) string {
	return t.Get(lang, m.Module, m.Text)
}

//Translated create translated message by given language
func (m *Message) Translated(lang string) *Translated {
	return NewTranslated(m, lang)
}

//NewMessage create new message with given module and text.
func NewMessage(module string, message string) *Message {
	return &Message{
		Module: module,
		Text:   message,
	}
}

//TemplateMessage message with template
type TemplateMessage struct {
	message *Message
	tokens  map[string]string
}

//Translate translate message with default translations.
//if lang if empty,Lang will be used.
func (m *TemplateMessage) Translate(lang string) string {
	if lang == "" {
		lang = Lang
	}
	return m.TranslateWith(DefaultTranslations, lang)
}

//TranslateWith translate message with given translations and language
func (m *TemplateMessage) TranslateWith(t *Translations, lang string) string {
	str := m.message.TranslateWith(t, lang)
	return Replace(str, m.tokens)
}

//Translated create translated message by given language
func (m *TemplateMessage) Translated(lang string) *Translated {
	return NewTranslated(m, lang)
}

//NewTemplateMessage create new message with given module and text.
func NewTemplateMessage(module string, message string, tokens map[string]string) *TemplateMessage {
	return &TemplateMessage{
		message: NewMessage(module, message),
		tokens:  tokens,
	}
}
