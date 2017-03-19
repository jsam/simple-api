package main

import (
	"fmt"
	"sync"
	"testing"
	"reflect"
)

// TestCacheSet will test if inserting new record in cache works correctly.
func TestCacheSet(t *testing.T) {
	c := NewCache()
	c.Set("test", 123)

	if c.state["test"] != 123 {
		t.Fatal("There is something terrible happenind with cache.")
	}

}


// TestCacheConcurrentSet will insert 5000 records concurrently, wait for their insertions and then check if values where inserted properly.
func TestCacheConcurrentSet(t *testing.T) {
	var wg sync.WaitGroup
	c := NewCache()

	for r := 0; r < 5000; r++ {
		key := fmt.Sprintf("test_%d", r)
		wg.Add(1)
		val := r
		go func() {
			c.Set(key, val)
			defer wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < 5000; i++ {
		key := fmt.Sprintf("test_%d", i)

		if c.state[key] != i {
			t.Fatalf("Concurent cache state is broken. %s: %d expected %d ", key, c.state[key], i)
		}
	}
}

// TestCacheGet will test for retrieval of key.
func TestCacheGet(t *testing.T) {
	c := NewCache()

	c.state["test"] = 1
	c.state["test_X"] = 32

	val, err := c.Get("test")
	if val != 1 || err != nil {
		t.Fatal("There is something terrible happening with cache.")
	}
}

// TestCacheConcurrentGet will test if reading from cache concurrently is correct.
func TestCacheConcurrentGet(t *testing.T) {
	var wg sync.WaitGroup
	c := NewCache()

	for r := 0; r < 5000; r++ {
		key := fmt.Sprintf("test_%d", r)
		c.state[key] = r
	}

	for i := 0; i < 5000; i++ {
		key := fmt.Sprintf("test_%d", i)
		wg.Add(1)
		expected := i
		go func() {
			val, err := c.Get(key)
			if val != expected || err != nil{
				t.Fatal("Something terrible happend with cache.")
			}
		}()
	}
}

// TestNewCache will test if creation of new object is successful.
func TestNewCache(t *testing.T) {
	c := NewCache()
	if c == nil {
		t.Fatal("Cache was not allocated.")
	}
}


// TestSingleCache will check if singleton pattern works on top of Cache object.
func TestSingleCache(t *testing.T) {
	s := SingleCache()
	if s == nil {
		t.Fatal("Creation of cache is broken.")
	}

	c := SingleCache()
	if c == nil {
		t.Fatal("Refrence fetch is broken.")
	}

	if !reflect.DeepEqual(s, c) {
		t.Fatal("Objects are not the same.")
	}

}