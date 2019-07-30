package translate

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

//NewMessage create new message with given module and text.
func NewMessage(module string, message string) *Message {
	return &Message{
		Module: module,
		Text:   message,
	}
}

type MessageWithTokens struct {
	message *Message
	tokens  map[string]string
}

//Translate translate message with default translations.
//if lang if empty,Lang will be used.
func (m *MessageWithTokens) Translate(lang string) string {
	if lang == "" {
		lang = Lang
	}
	return m.TranslateWith(DefaultTranslations, lang)
}

//TranslateWith translate message with given translations and language
func (m *MessageWithTokens) TranslateWith(t *Translations, lang string) string {
	str := m.TranslateWith(t, lang)
	return Replace(str, m.tokens)
}

//NewMessageWithTokens create new message with given module and text.
func NewMessageWithTokens(module string, message string, tokens map[string]string) *MessageWithTokens {
	return &MessageWithTokens{
		message: NewMessage(module, message),
		tokens:  tokens,
	}
}
