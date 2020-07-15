package token

import (
	"crypto/rand"
	"math/big"
)

type Generator interface {
	Generate() (Secret, error)
}

type BytesGenerator int

func (g BytesGenerator) Generate() (Secret, error) {
	buffer := make([]byte, int(g))
	_, err := rand.Read(buffer)
	if err != nil {
		return nil, err
	}
	return Secret(buffer), nil
}

type ListGenerator struct {
	List []byte
	Min  int
	Max  int
}

func (g *ListGenerator) Generate() (Secret, error) {
	var length int
	if g.Max <= g.Min {
		length = int(g.Min)
	} else {
		max, err := rand.Int(rand.Reader, big.NewInt(int64(g.Max-g.Min)))
		if err != nil {
			return nil, err
		}
		length = g.Min + int(max.Int64())
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(g.List))))
		if err != nil {
			return nil, err
		}
		result[i] = g.List[int(index.Int64())]
	}
	return Secret(result), nil
}
