package etcc

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/etcdmain"
	"github.com/stretchr/testify/assert"
	. "github.com/topicai/candy"
)

func ExampleNew() {
	flagEtcd := flag.String("etcd", "http://127.0.0.1:4001,http://127.0.0.1:2379", "Etcd peers")
	flag.Parse()
	c, _ := New(*flagEtcd)
	c.Mkdir("/home/yi")
	c.Set("/home/yi/a", "Apple")
	c.Set("/home/yi/b", "Banana")
	c.Get("/home/yi/a")
	c.Get("home/yi/b")
	c.Rmdir("/home")
}

var (
	c *Etcd
)

func init() {
	os.Args = os.Args[0:1] // Make etcdmain.Main parse and get all-default settings.
	go etcdmain.Main()

	time.Sleep(3 * time.Second) // NOTE: Give etcd 3 seconds to start before connecting to it.

	var e error
	c, e = New("http://localhost:4001,http://localhost:2379")
	Must(e)
}

func TestEtcdBasicSession(t *testing.T) {
	assert := assert.New(t)

	c.Rmdir("/home/yi")

	assert.Nil(c.Mkdir("/home/yi"))
	assert.NotNil(c.Mkdir("/home/yi"))

	assert.Nil(c.Set("/home/yi/a", "Apple"))
	assert.Nil(c.Set("/home/yi/b", "Banana"))

	r, e := c.Get("/home/yi/a")
	assert.Nil(e)
	assert.Equal("Apple", r)

	r, e = c.Get("home/yi/b")
	assert.Nil(e)
	assert.Equal("Banana", r)

	assert.Nil(c.Set("/home/yi/a", "Aloha"))
	r, e = c.Get("/home/yi/a")
	assert.Nil(e)
	assert.Equal("Aloha", r)

	assert.Nil(c.Rmdir("/home"))
}

func TestSetWithTTL(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(c.SetWithTTL("key-ttl", "value-ttl", 1000*time.Millisecond))

	r, e := c.Get("key-ttl")
	assert.Nil(e)
	assert.Equal("value-ttl", r)

	// NOTE: it seems that etcd doesn't keep TTL exactly, if we
	// wait for exactly the TTL time (1000 * timeMillisecond in
	// this case), we might find the key-value pair still there.
	time.Sleep(1500 * time.Millisecond)
	r, e = c.Get("key-ttl")
	assert.NotNil(e)
	assert.Equal("", r)
}
