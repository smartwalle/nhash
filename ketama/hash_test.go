package ketama_test

import (
	"fmt"
	"github.com/smartwalle/nhash/ketama"
	"strconv"
	"testing"
)

func TestHash_Get(t *testing.T) {
	var nHash = ketama.New[string](ketama.WithSpots(200))

	nodes := map[string]int{
		"test1.server.com": 1,
		"test2.server.com": 1,
		"test3.server.com": 2,
		"test4.server.com": 5,
	}

	for k, v := range nodes {
		nHash.Add(k, k, v)
	}

	nHash.Prepare()

	m := make(map[string]int)
	for i := 0; i < 1e6; i++ {
		m[nHash.Get("test value"+strconv.FormatUint(uint64(i), 10))]++
	}

	for k := range nodes {
		fmt.Println(k, m[k])
	}
}
