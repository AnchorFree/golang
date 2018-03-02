package kv

import (
	"errors"
	consul "github.com/hashicorp/consul/api"
	"net/http"
	"strconv"
	"time"
)

type Consul struct {
	kv *consul.KV
}

func (c *Consul) Init(opts []string) error {

	if len(opts) < 1 {
		return errors.New("address required to init consul store")
	}
	var (
		timeout int
		err     error
	)

	config := consul.DefaultConfig()
	config.Address = opts[0]
	if len(opts) > 1 {
		timeout, err = strconv.Atoi(opts[1])
		if err != nil {
			return err
		}
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	config.HttpClient = client
	ins, err := consul.NewClient(config)
	if err != nil {
		return err
	}
	c.kv = ins.KV()
	return nil

}

func (c *Consul) Put(key string, value []byte) error {

	kvp := &consul.KVPair{Key: key, Value: value}
	_, err := c.kv.Put(kvp, nil)
	if err != nil {
		return err
	}
	return nil

}

func (c *Consul) Get(key string) ([]byte, error) {

	kvp, _, err := c.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	if kvp == nil {
		return nil, errors.New("key does not exist")
	}
	return kvp.Value, nil

}

func (c *Consul) Delete(prefix string) error {

	_, err := c.kv.Delete(prefix, nil)
	if err != nil {
		return err
	}
	return nil

}

func (c *Consul) DeleteTree(prefix string) error {

	_, err := c.kv.DeleteTree(prefix, nil)
	if err != nil {
		return err
	}
	return nil

}

func (c *Consul) List(prefix string) ([]string, error) {

	kvpairs, _, err := c.kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, kvp := range kvpairs {
		keys = append(keys, kvp.Key)
	}
	return keys, nil

}

func (c *Consul) Close() error {

	return nil

}
