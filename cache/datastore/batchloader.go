package datastore

// BatchLoader batch data loader interface
type BatchLoader interface {
	// NewDataElement create empty data element.
	// Must return pointer of data.
	NewDataElement() interface{}
	//BatchLoadData   load data by given keys to map of data pointers.
	// Return map of data pointers and any error if raised.
	BatchLoadData(...string) (map[string]interface{}, error)
}
