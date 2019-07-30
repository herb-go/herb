package translate

import (
	"strings"
)

//Replace replace string with given replace mapper.
//Words begin with { and end with } will be replace with mappedvalue
//Use "\{" , "\}" , "\\" to keep "{" , "}"  "\" char.
func Replace(str string, m map[string]string) string {
	values := make([]string, 2*len(m)+6)
	values[0] = "\\\\"
	values[1] = "\\"
	values[2] = "\\{"
	values[3] = "{"
	values[4] = "\\}"
	values[5] = "}"
	var i = 6
	for k := range m {
		values[i] = "{" + k + "}"
		values[i+1] = m[k]
		i += 2
	}
	r := strings.NewReplacer(values...)
	return r.Replace(str)
}
