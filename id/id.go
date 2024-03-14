package id

import (
	"math/rand"

	"github.com/bwmarrin/snowflake"
)

var gGenerator *Generator

type Generator struct {
	*snowflake.Node
}
type ID struct {
	snowflake.ID
}

func init() {
	gGenerator, _ = NewGenerator()
}

func Next() string {
	return gGenerator.Next().String()
}

func NewGenerator() (*Generator, error) {
	node := rand.Int63() % 1024
	n, err := snowflake.NewNode(node)
	if err != nil {
		return nil, err
	}
	return &Generator{Node: n}, nil
}

func (i ID) WithPrefix(prefix string) string {
	return prefix + i.String()
}

func (g *Generator) Next() ID {
	return ID{
		ID: g.Generate(),
	}
}
