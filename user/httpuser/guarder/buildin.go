package guarder

import "net/http"

type FixedIDDriver struct {
	ID string
}

func (g *FixedIDDriver) MustCreateGuarder(f Field) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		f.StoreID(r, g.ID)
		next(w, r)
	}
}

func FixedIDFactory(loader func(v interface{}) error) (Driver, error) {
	c := &FixedIDDriver{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func EmptyFactory(loader func(v interface{}) error) (Driver, error) {
	return &FixedIDDriver{
		ID: "",
	}, nil
}
func RegisterBuildinDrivers() {
	Register("fixedix", FixedIDFactory)
}
