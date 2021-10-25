package idgenerator

import (
	"testing"
)

func TestAll(t *testing.T) {
	NewIdGenerator()
	id := Generate()
	t.Log(id)
	id = Generate()
	t.Log(id)
	id = Generate()
	t.Log(id)
}
