package kv

type Store interface {
	Init(opts []string) error
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Delete(key string) error
	List(prefix string) ([]string, error)
	DeleteTree(prefix string) error
	Close() error
}
