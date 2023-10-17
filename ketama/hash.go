package ketama

// Package ketama implements consistent hashing compatible with Algorithm::ConsistentHash::Ketama
// This is a fork of https://github.com/mncaudill/ketama/blob/master/ketama.go

import (
	"github.com/smartwalle/nhash/ketama/internal/sha1"
	"hash"
	"sort"
	"strconv"
)

type node[T any] struct {
	key   string
	value T
	point uint32
}

func (n *node[T]) String() string {
	return n.key
}

type nodeList[T any] []node[T]

func (nl nodeList[T]) Len() int           { return len(nl) }
func (nl nodeList[T]) Less(i, j int) bool { return nl[i].point < nl[j].point }
func (nl nodeList[T]) Swap(i, j int)      { nl[i], nl[j] = nl[j], nl[i] }

type Option func(opts *options)

type options struct {
	spots int
	hash  func() hash.Hash32
}

func WithSpots(spots int) Option {
	return func(opts *options) {
		if spots <= 0 {
			spots = 4
		}
		opts.spots = spots
	}
}

func WithHash(h func() hash.Hash32) Option {
	return func(opts *options) {
		if h == nil {
			return
		}
		opts.hash = h
	}
}

type Hash[T any] struct {
	opts   *options
	nodes  nodeList[T]
	length int
	empty  T
}

func New[T any](opts ...Option) *Hash[T] {
	var nHash = &Hash[T]{}
	nHash.opts = &options{
		spots: 4,
		hash:  sha1.New,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(nHash.opts)
		}
	}

	return nHash
}

func (h *Hash[T]) Add(key string, value T, weight int) {
	var r = h.opts.spots * weight

	var nHash = h.opts.hash()

	for i := 0; i < r; i++ {
		nHash.Write([]byte(key + ":" + strconv.Itoa(i)))
		var nNode = node[T]{
			key:   key,
			value: value,
			point: nHash.Sum32(),
		}
		h.nodes = append(h.nodes, nNode)
		nHash.Reset()
	}
}

func (h *Hash[T]) Prepare() {
	sort.Slice(h.nodes, func(i, j int) bool {
		return h.nodes[i].point < h.nodes[j].point
	})
	h.length = len(h.nodes)
}

func (h *Hash[T]) Get(key string) T {
	if len(h.nodes) == 0 {
		return h.empty
	}

	var nHash = h.opts.hash()
	nHash.Write([]byte(key))
	var value = nHash.Sum32()

	i := sort.Search(h.length, func(i int) bool {
		return h.nodes[i].point >= value
	})

	if i == h.length {
		i = 0
	}
	return h.nodes[i].value
}
