package idgenerator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	var err error
	randint, err := rand.Int(rand.Reader, big.NewInt(1024))
	if err != nil {
		panic("snowflake randint failed")
	}
	randomNub := randint.Int64()
	node, err = snowflake.NewNode(randomNub)
	if err != nil {
		panic("snowflake init failed")
	}

	fmt.Printf("!!!warning, please check random number in every machine\n!!!this machine's random nuber is %d.\n", randomNub)
	if err = os.WriteFile("warning_check_randomnub", []byte(strconv.FormatInt(randomNub, 10)), 0666); err != nil {
		panic(err)
	}
}

func Generate() int64 {
	id := int64(node.Generate())
	return id
}
