package ketama

// Package ketama implements consistent hashing compatible with Algorithm::ConsistentHash::Ketama
// This is a fork of https://github.com/mncaudill/ketama/blob/master/ketama.go written in a
// more extendable way

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

func (this *node[T]) String() string {
	return this.key
}

type nodeList[T any] []node[T]

func (n nodeList[T]) Len() int           { return len(n) }
func (n nodeList[T]) Less(i, j int) bool { return n[i].point < n[j].point }
func (n nodeList[T]) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

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

func (this *Hash[T]) Add(key string, value T, weight int) {
	var r = this.opts.spots * weight

	var nHash = this.opts.hash()

	for i := 0; i < r; i++ {
		nHash.Write([]byte(key + ":" + strconv.Itoa(i)))
		var nNode = node[T]{
			key:   key,
			value: value,
			point: nHash.Sum32(),
		}
		this.nodes = append(this.nodes, nNode)
		nHash.Reset()
	}
}

func (this *Hash[T]) Prepare() {
	sort.Slice(this.nodes, func(i, j int) bool {
		return this.nodes[i].point < this.nodes[j].point
	})
	this.length = len(this.nodes)
}

func (this *Hash[T]) Get(key string) T {
	if len(this.nodes) == 0 {
		return this.empty
	}

	var nHash = this.opts.hash()
	nHash.Write([]byte(key))
	var hValue = nHash.Sum32()

	i := sort.Search(this.length, func(i int) bool {
		return this.nodes[i].point >= hValue
	})

	if i == this.length {
		i = 0
	}
	return this.nodes[i].value
}
