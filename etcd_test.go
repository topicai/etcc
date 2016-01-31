package etcd

import (
	"flag"
	"log"
	"os/exec"
	"path"
	"testing"
	"time"

	. "github.com/topicai/candy"
)

func ExampleNew() {
	flagEtcd := flag.String("etcd", "http://127.0.0.1:4001,http://127.0.0.1:2379", "Etcd peers addresses")
	flag.Parse()
	c, _ := New(*flagEtcd)
	c.Mkdir("/home/yi")
	c.Set("/home/yi/a", "Apple")
	c.Set("/home/yi/b", "Banana")
	c.Get("/home/yi/a")
	c.Get("home/yi/b")
	c.Rmdir("/home")
}

func init() {
	// Build and run etcd.
	buf, e := exec.Command("go", "get", "github.com/coreos/etcd/...").CombinedOutput()
	if e != nil {
		log.Panicf("Failed go get github.com/coreos/etcd/... : %v\n%s", e, buf)
	}

	go Must(exec.Command(path.Join(GoPath(), "bin/etcd")).Start())
}

func TestEtcd(t *testing.T) {
	time.Sleep(2 * time.Second) // Wait 2 seconds for etcd to get ready to serve.

	if c, e := New("http://127.0.0.1:4001,http://127.0.0.1:2379"); e != nil {
		t.Error(e)
	} else {
		c.Rmdir("/home/yi")
		if e := c.Mkdir("/home/yi"); e != nil {
			t.Error(e)
		}
		if e := c.Set("/home/yi/a", "Apple"); e != nil {
			t.Error(e)
		}
		if e := c.Set("/home/yi/b", "Banana"); e != nil {
			t.Error(e)
		}
		if r, e := c.Get("/home/yi/a"); e != nil || r != "Apple" {
			t.Error(e)
		}
		if r, e := c.Get("home/yi/b"); e != nil || r != "Banana" {
			t.Error(e)
		}

		if e := c.Set("/home/yi/a", "Aloha"); e != nil {
			t.Error(e)
		}
		if r, e := c.Get("/home/yi/a"); e != nil || r != "Aloha" {
			t.Error(e)
		}

		if e := c.Rmdir("/home"); e != nil {
			t.Error(e)
		}
	}
}
