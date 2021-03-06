package render

//Data render data struct.
type Data map[string]interface{}

//NewData create new render data.
func NewData() *Data {
	return &Data{}
}

//Set set data field value by key.
func (d *Data) Set(key string, data interface{}) {
	if *d == nil {
		data := Data(map[string]interface{}{})
		*d = data
	}
	(*d)[key] = data
}

//Del delete data field value by key.
func (d *Data) Del(key string) {
	if *d == nil {
		return
	}
	delete(*d, key)
}

//Get Get data field value by key.
func (d *Data) Get(key string) interface{} {
	if *d == nil {
		return nil
	}
	data, ok := (*d)[key]
	if ok == false {
		return nil
	}
	return data
}

//Merge merger two render data.
func (d *Data) Merge(data *Data) {
	if *d == nil {
		data := Data(map[string]interface{}{})
		*d = data
	}
	if data == nil || *data == nil {
		return
	}
	for k, v := range *data {
		d.Set(k, v)
	}
}
