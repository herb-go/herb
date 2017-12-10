package user

type ProfileIndex string

const ProfileIndexName = ProfileIndex("Name")
const ProfileIndexEmail = ProfileIndex("Email")
const ProfileIndexNickname = ProfileIndex("Nickname")
const ProfileIndexAvatar = ProfileIndex("Avatar")
const ProfileIndexProfileURL = ProfileIndex("ProfileURL")
const ProfileIndexAccessToken = ProfileIndex("AccessToken")
const ProfileIndexGender = ProfileIndex("Gender")
const ProfileIndexCompany = ProfileIndex("Company")
const ProfileIndexID = ProfileIndex("ID")
const ProfileIndexLocation = ProfileIndex("Location")
const ProfileIndexWebsite = ProfileIndex("Website")

const ProfileGenderMale = "M"
const ProfileGenderFemale = "F"

type Profile map[ProfileIndex][]string

func (p *Profile) Value(index ProfileIndex) string {
	data, ok := (*p)[index]
	if ok == false || len(data) == 0 {
		return ""
	}
	return data[0]
}

func (p *Profile) Values(index ProfileIndex) []string {
	data, ok := (*p)[index]
	if ok == false {
		return nil
	}
	return data
}

func (p *Profile) SetValue(index ProfileIndex, value string) {
	(*p)[index] = []string{value}
}

func (p *Profile) SetValues(index ProfileIndex, values []string) {
	(*p)[index] = values
}

func (p *Profile) AddValue(index ProfileIndex, value string) {
	data, ok := (*p)[index]
	if ok == false {
		data = []string{}
	}
	data = append(data, value)
	(*p)[index] = data
}
