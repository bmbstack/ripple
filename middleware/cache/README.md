# cache


## import
```
	"github.com/weisd/cache"
	_ "github.com/weisd/cache/redis"

```

## Documentation
```
package cache

import (
	"testing"
)

func Test_TagCache(t *testing.T) {

	c, err := New(Options{Adapter: "memory"})
	if err != nil {
		t.Fatal(err)
	}

	// base use
	err = c.Put("da", "weisd", 300)
	if err != nil {
		t.Fatal(err)
	}

	res := c.Get("da")

	if res != "weisd" {
		t.Fatal("base put faield")
	}

	t.Log("ok")

	// use tags/namespace
	err = c.Tags([]string{"dd"}).Put("da", "weisd", 300)
	if err != nil {
		t.Fatal(err)
	}
	res = c.Tags([]string{"dd"}).Get("da")

	if res != "weisd" {
		t.Fatal("tags put faield")
	}

	t.Log("ok")

	err = c.Tags([]string{"aa"}).Put("aa", "aaa", 300)
	if err != nil {
		t.Fatal(err)
	}

	res = c.Tags([]string{"aa"}).Get("aa")

	if res != "aaa" {
		t.Fatal("not aaa")
	}

	t.Log("ok")

	// flush namespace
	err = c.Tags([]string{"aa"}).Flush()
	if err != nil {
		t.Fatal(err)
	}

	res = c.Tags([]string{"aa"}).Get("aa")
	if res != "" {
		t.Fatal("flush faield")
	}

	res = c.Tags([]string{"aa"}).Get("bb")
	if res != "" {
		t.Fatal("flush faield")
	}

	// still store in
	res = c.Tags([]string{"dd"}).Get("da")
	if res != "weisd" {
		t.Fatal("where ")
	}

	t.Log("ok")

}
```


## echo Middleware
```
e := echo.New()
e.Use(cache.EchoCacher(cache.Options{Adapter: "redis", AdapterConfig: `{"Addr":":6379"}`, Section: "test", Interval: 5}))

e.Get("/test/cache/put", func(c *echo.Context) error {
	err := cache.Store(c).Put("name", "weisd", 10)
	if err != nil {
		return err
	}

	return c.String(200, "store ok")
})

e.Get("/test/cache/get", func(c *echo.Context) error {
	name := cache.Store(c).Get("name")

	return c.String(200, "get name %s", name)
})

e.Run(":1323")

```

##