package idgenerator

import (
	"crypto/rand"
	"math/big"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	randint, err := rand.Int(rand.Reader, big.NewInt(1024))
	if err != nil {
		panic("snowflake randint failed")
	}
	node, err = snowflake.NewNode(randint.Int64())
	if err != nil {
		panic("snowflake init failed")
	}
}

func Generate() int64 {
	id := int64(node.Generate())
	return id
}
