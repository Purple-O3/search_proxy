package router

import (
	"hash"
	"hash/fnv"
	"math/rand"
	"search_proxy/internal/util/tools"
)

type Router interface {
	LoadBalance(key string) int
}

func RouterFactory(routerType string, machineNub int) Router {
	switch routerType {
	case "poll":
		return newPoller(machineNub)
	case "rand":
		return newRander(machineNub)
	case "hash":
		return newHasher(machineNub)
	default:
		return newPoller(machineNub)
	}
}

type poller struct {
	machineNub int
	index      int
}

func newPoller(machineNub int) *poller {
	p := new(poller)
	p.machineNub = machineNub
	p.index = 0
	return p
}

func (p *poller) LoadBalance(key string) int {
	p.index = (p.index + 1) % p.machineNub
	return p.index
}

type rander struct {
	machineNub int
}

func newRander(machineNub int) *rander {
	r := new(rander)
	r.machineNub = machineNub
	return r
}

func (r *rander) LoadBalance(key string) int {
	return rand.Intn(r.machineNub)
}

type hasher struct {
	machineNub int
	hashFn     hash.Hash64
}

func newHasher(machineNub int) *hasher {
	h := new(hasher)
	h.machineNub = machineNub
	h.hashFn = fnv.New64()
	return h
}

func (h *hasher) LoadBalance(hashKey string) int {
	keyByte := tools.Str2Bytes(hashKey)
	h.hashFn.Reset()
	h.hashFn.Write(keyByte)
	value := h.hashFn.Sum64()
	return int(value) % h.machineNub
}
