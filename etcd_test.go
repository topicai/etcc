package etcd

import "flag"

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
