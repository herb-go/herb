package render

type Data map[string]interface{}

func (d *Data) Set(key string, data interface{}) {
	if d == nil {
		data := Data(map[string]interface{}{})
		d = &(data)
	}
	(*d)[key] = data
}

func (d *Data) Del(key string) {
	if d == nil {
		return
	}
	delete(*d, key)
}

func (d *Data) Get(key string) interface{} {
	if d == nil {
		return nil
	}
	data, ok := (*d)[key]
	if ok == false {
		return nil
	}
	return data
}
func (d *Data) Merge(data *Data) {
	if data == nil || *data == nil {
		return
	}
	for k, v := range *data {
		d.Set(k, v)
	}
}
