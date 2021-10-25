package router

import "testing"

func TestAll(t *testing.T) {
	router := RouterFactory("poll", 10)
	ret := router.LoadBalance("127.0.0.1")
	t.Log(ret)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)

	router = RouterFactory("rand", 10)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)

	router = RouterFactory("hash", 10)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)
	ret = router.LoadBalance("127.0.0.1")
	t.Log(ret)
}
