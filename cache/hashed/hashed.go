package hashed

//Data cache data struct
type Data struct {
	Key     string
	Expired int64
	Data    []byte
}

//NewData create new data with given key,expiired time stamp and data
func NewData(key string, expired int64, data []byte) *Data {
	return &Data{
		Key:     key,
		Expired: expired,
		Data:    data,
	}
}

//Status hash data status
type Status struct {
	FirstExpired int64
	LastExpired  int64
	Size         int
	Delta        int
	Changed      bool
}

func (s *Status) calc(data *Data, current int64) bool {
	if data.Expired < current {
		s.Changed = true
		return false
	}
	if s.LastExpired <= 0 || s.FirstExpired > data.Expired {
		s.FirstExpired = data.Expired
	}
	if s.LastExpired < data.Expired {
		s.LastExpired = data.Expired
	}
	s.Size = s.Size + len(data.Data)
	return true
}

//NewStatus create new status
func NewStatus() *Status {
	return &Status{
		FirstExpired: 0,
		LastExpired:  0,
		Size:         0,
		Changed:      false,
	}
}

//Hashed cache hashed data.
type Hashed []*Data

//New create new cache hashed data.
func New() *Hashed {
	h := Hashed([]*Data{})
	return &h
}
func (h *Hashed) isEmpty() bool {
	return len(*h) == 0
}
func (h *Hashed) set(data *Data, current int64) *Status {
	result := make(Hashed, 0, len(*h)+1)
	status := NewStatus()
	status.Changed = true
	for k := range *h {
		if (*h)[k].Key != data.Key && status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
		} else {
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	if status.calc(data, current) {
		result = append(result, data)
		status.Delta = status.Delta + len(data.Data)
	}
	*h = result
	return status
}
func (h *Hashed) update(data *Data, current int64) *Status {
	result := make(Hashed, 0, len(*h))
	status := NewStatus()

	for k := range *h {
		var delta int
		d := (*h)[k]
		if d.Key == data.Key {
			status.Changed = true
			delta = len(data.Data) - len(d.Data)
			d.Expired = data.Expired
			d.Data = data.Data
		}
		if status.calc(d, current) {
			result = append(result, d)
			status.Delta = status.Delta + delta
		} else {
			status.Delta = status.Delta - len(d.Data)
		}
	}

	*h = result
	return status
}
func (h *Hashed) expired(key string, expired int64, current int64) *Status {
	result := make(Hashed, 0, len(*h)+1)
	status := NewStatus()
	status.Changed = true
	for k := range *h {
		d := (*h)[k]
		if d.Key == key {
			d.Expired = expired
		}
		if status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
		} else {
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	*h = result
	return status
}
func (h *Hashed) delete(key string, current int64) *Status {
	result := make(Hashed, 0, len(*h))
	status := NewStatus()
	for k := range *h {
		if (*h)[k].Key != key && status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
		} else {
			status.Changed = true
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	*h = result
	return status
}
func (h *Hashed) selfcheck(current int64) *Status {
	result := make(Hashed, 0, len(*h))
	status := NewStatus()
	for k := range *h {
		if status.calc((*h)[k], current) {
			result = append(result, (*h)[k])
		} else {
			status.Delta = status.Delta - len((*h)[k].Data)
		}
	}
	*h = result
	return status
}

func (h *Hashed) get(key string, current int64) *Data {
	for k := range *h {
		if (*h)[k].Key == key {
			if (*h)[k].Expired >= current {
				return (*h)[k]
			}
			return nil
		}
	}
	return nil
}
