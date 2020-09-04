package persist

type Factory func(loader func(v interface{}) error) (Store, error)
