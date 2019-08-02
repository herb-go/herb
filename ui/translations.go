package ui

//Translations messages collection grouped by lang and module.
type Translations map[string]map[string]*Messages

//SetMessages set collection messages by given lang and module
func (c *Translations) SetMessages(lang string, module string, m *Messages) {
	if (*c)[lang] == nil {
		(*c)[lang] = map[string]*Messages{}
	}
	(*c)[lang][module] = m
}

//Load load translated message by given lang.module and key.
//Return translated message and true if found.
//Return raw message and false if not found
func (c *Translations) Load(lang string, module string, key string) (string, bool) {
	m := c.GetMessages(lang, module)
	if m == nil {
		return key, false
	}
	return m.Load(key)
}

//Get load get message by given lang,module and key.
//Return translated message  if found.
//Return raw message  if not found
func (c *Translations) Get(lang string, module string, key string) string {
	v, _ := c.Load(lang, module, key)
	return v
}

//GetMessages get messages by given lang and module
//Return nil if messages not found
func (c *Translations) GetMessages(lang string, module string) *Messages {
	if lang == "" {
		return nil
	}
	if (*c)[lang] == nil {
		return nil
	}
	if module == "" {
		return nil
	}
	if (*c)[lang][module] == nil {
		return nil
	}
	return (*c)[lang][module]
}

// NewTranslations  create new messages Translations
func NewTranslations() *Translations {
	return &Translations{}
}

//DefaultTranslations default messages Translations
var DefaultTranslations = NewTranslations()

//Lang default langauge used by Translate func.
var Lang = ""

//Get translate messsage  by default langauge ,given module and ky.
func Get(module string, key string) string {
	return GetIn(Lang, module, key)
}

//GetIn translate messsage by given langauge module and ky.
func GetIn(lang string, module string, key string) string {
	return DefaultTranslations.Get(lang, module, key)
}

//GetMessages by given lang and modules
//Default Lang will used if lang is empty
func GetMessages(lang string, module string) *Messages {
	if lang == "" {
		lang = Lang
	}
	return DefaultTranslations.GetMessages(lang, module)
}
