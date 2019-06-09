package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

func main() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cache.New(5*time.Minute, 10*time.Minute)

	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("foo", "bar", cache.DefaultExpiration)

	// Set the value of the key "baz" to 42, with no expiration time
	// (the item won't be removed until it is re-set, or removed using
	// c.Delete("baz")
	var num = 1000
	c.Set("baz", num, cache.NoExpiration)

	// Get the string associated with the key "foo" from the cache
	fooV1, found := c.Get("foo")
	if found {
		fmt.Printf("fooV1: %s \n", fooV1)
	}

	// Since Go is statically typed, and cache values can be anything, type
	// assertion is needed when values are being passed to functions that don't
	// take arbitrary types, (i.e. interface{}). The simplest way to do this for
	// values which will only be used once--e.g. for passing to another
	// function--is:
	fooV2, found := c.Get("foo")
	if found {
		fmt.Printf("fooV2: %s \n", fooV2)
		//MyFunction(foo.(string))
	}

	// This gets tedious if the value is used several times in the same function.
	// You might do either of the following instead:
	if x, found := c.Get("foo"); found {
		fooV3 := x.(string)
		fmt.Printf("fooV3: %s \n", fooV3)
	}
	// or
	var fooV4 int
	if x, found := c.Get("baz"); found {
		fooV4 = x.(int)
		fmt.Printf("fooV4: %d \n", fooV4)
	}
	// ...
	// foo can then be passed around freely as a string

	// Want performance? Store pointers!
	//c.Set("foo", &MyStruct, cache.DefaultExpiration)
	//if x, found := c.Get("foo"); found {
	//	foo := x.(*MyStruct)
	//	// ...
	//}
}
