package etcd

import (
	"flag"
	"fmt"
	"testing"
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

var testEtcd = flag.String("testEtcd", "", "If non-empty, run tests.  An example is \"http://127.0.0.1:4001,http://127.0.0.1:2379\"")

func TestEtcd(t *testing.T) {
	if *testEtcd == "" {
		fmt.Println("Test doesn't run without specifying testEtcd")
		return
	}

	fmt.Println("Testing using etcd listening on ", *testEtcd)

	if c, e := New(*testEtcd); e != nil {
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
