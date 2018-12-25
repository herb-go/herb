package notificationpartmsgpack

import "github.com/herb-go/herb/notification"
import "github.com/vmihailenco/msgpack"

var Marshal = func(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}
var Unmarshal = func(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

type StringPart struct {
	Name string
}

func NewStringPart(name string) *StringPart {
	return &StringPart{
		Name: name,
	}
}

func (p *StringPart) Set(Notification *notification.PartedNotification, data string) error {
	Notification.Parts[p.Name] = []byte(data)
	return nil
}
func (p *StringPart) Get(Notification *notification.PartedNotification) (string, error) {
	bytes := Notification.Parts[p.Name]
	return string(bytes), nil
}

type BinaryPart struct {
	Name string
}

func NewBinaryPart(name string) *BinaryPart {
	return &BinaryPart{
		Name: name,
	}
}

func (p *BinaryPart) Set(Notification *notification.PartedNotification, data []byte) error {
	Notification.Parts[p.Name] = data
	return nil
}
func (p *BinaryPart) Get(Notification *notification.PartedNotification) ([]byte, error) {
	data := Notification.Parts[p.Name]
	return data, nil
}

type StringListPart struct {
	Name string
}

func NewStringListPart(name string) *StringListPart {
	return &StringListPart{
		Name: name,
	}
}

func (p *StringListPart) Set(Notification *notification.PartedNotification, data []string) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	Notification.Parts[p.Name] = bs
	return nil
}
func (p *StringListPart) Get(Notification *notification.PartedNotification) ([]string, error) {
	bytes := Notification.Parts[p.Name]
	data := []string{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type StringMapPart struct {
	Name string
}

func NewStringMapPart(name string) *StringMapPart {
	return &StringMapPart{
		Name: name,
	}
}

func (p *StringMapPart) Set(Notification *notification.PartedNotification, data map[string]string) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	Notification.Parts[p.Name] = bs
	return nil
}
func (p *StringMapPart) Get(Notification *notification.PartedNotification) (map[string]string, error) {
	bytes := Notification.Parts[p.Name]
	data := map[string]string{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type BinaryListPart struct {
	Name string
}

func NewBinaryListPart(name string) *BinaryListPart {
	return &BinaryListPart{
		Name: name,
	}
}

func (p *BinaryListPart) Set(Notification *notification.PartedNotification, data [][]byte) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	Notification.Parts[p.Name] = bs
	return nil
}
func (p *BinaryListPart) Get(Notification *notification.PartedNotification) ([][]byte, error) {
	bytes := Notification.Parts[p.Name]
	data := [][]byte{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type BinaryMapPart struct {
	Name string
}

func NewBinaryMapPart(name string) *BinaryMapPart {
	return &BinaryMapPart{
		Name: name,
	}
}

func (p *BinaryMapPart) Set(Notification *notification.PartedNotification, data map[string][]byte) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	Notification.Parts[p.Name] = bs
	return nil
}
func (p *BinaryMapPart) Get(Notification *notification.PartedNotification) (map[string][]byte, error) {
	bytes := Notification.Parts[p.Name]
	data := map[string][]byte{}
	err := Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type InterfacePart struct {
	Name string
}

func NewInterfacePart(name string) *InterfacePart {
	return &InterfacePart{
		Name: name,
	}
}

func (p *InterfacePart) Marshal(Notification *notification.PartedNotification, data interface{}) error {
	bs, err := Marshal(data)
	if err != nil {
		return err
	}
	Notification.Parts[p.Name] = bs
	return nil
}
func (p *InterfacePart) Unmarshal(Notification *notification.PartedNotification, data interface{}) error {
	bytes := Notification.Parts[p.Name]
	return Unmarshal(bytes, &data)
}
