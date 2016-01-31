package etcd

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"
	"golang.org/x/net/context"
)

type Etcd struct {
	api client.KeysAPI
}

// New trys to connect to etcd server.  endpoints must be addreses
// delimited by comma, like "http://127.0.0.1:4001,http://127.0.0.1:2379".
func New(endpoints string) (*Etcd, error) {
	eps := strings.Split(endpoints, ",")
	for i, ep := range eps {
		u, e := url.Parse(ep)
		if e != nil {
			return nil, fmt.Errorf("url.Parse: %v", e)
		}

		if u.Scheme == "" {
			u.Scheme = "http"
		}
		eps[i] = u.String()
	}

	tr, e := transport.NewTransport(transport.TLSInfo{},
		10*time.Second) // timeout = 10sec
	if e != nil {
		return nil, fmt.Errorf("transport.NewTransport: %v", e)
	}

	c, e := client.New(client.Config{Endpoints: eps, Transport: tr})
	if e != nil {
		return nil, fmt.Errorf("client.New: %v", e)
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	e = c.Sync(ctx)
	cancel()
	if e != nil {
		return nil, fmt.Errorf("(etc)client.Sync: %v", e)
	}

	return &Etcd{client.NewKeysAPI(c)}, nil
}

// Mkdir creates a directory. The directory could be multiple-level,
// like /home/yi/hello. But it must not exist before; otherwise Mkdir
// returns an error.
func (c *Etcd) Mkdir(dir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if _, e := c.api.Set(ctx, dir, "", &client.SetOptions{Dir: true, PrevExist: client.PrevNoExist}); e != nil {
		return fmt.Errorf("Etcd.Mkdir: %v", e)
	}
	return nil
}

func (c *Etcd) Set(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if _, e := c.api.Set(ctx, key, value, &client.SetOptions{}); e != nil {
		return fmt.Errorf("Etcd.Set: %v", e)
	}
	return nil
}

func (c *Etcd) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	r, e := c.api.Get(ctx, key, &client.GetOptions{Sort: true})
	if e != nil {
		return "", fmt.Errorf("Etcd.Get: %v", e)
	}
	return r.Node.Value, nil
}

// Rmdir removes a directory and recursively its all content, like bash command `rm -rf`.
func (c *Etcd) Rmdir(dir string) error {
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()
	if _, e := c.api.Delete(ctx, "/home", &client.DeleteOptions{Dir: true, Recursive: true}); e != nil {
		return fmt.Errorf("Etcd.Rmdir: %v", e)
	}
	return nil
}
