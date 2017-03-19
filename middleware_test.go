package main

import (
	"testing"
	"time"
	"sync"
	"reflect"
)

// TestNewContext test creation of new context object.
func TestNewContext(t *testing.T) {
	freshContext := NewContext(false)
	if freshContext == nil {
		t.Fatal("Creation of context failed.")
	}
}

// TestFlushContext test for flushing the context to hard disk.
func TestFlushContext(t *testing.T) {
	server = &API{started_at: time.Now(),}
	freshContext := NewContext(false)
	err := freshContext.flush()
	if err != nil {
		t.Fatal("Flushing is broken.")
	}
}


// TestLoadContext test for sourcing the context object from hard disk.
func TestLoadContext(t *testing.T) {
	server = &API{started_at: time.Now(),}
	freshContext := NewContext(false)
	freshContext.flush()

	oldContext := NewContext(true)
	if oldContext == nil {
		t.Fatal("Sourcing old context failed.")
	}

	if !reflect.DeepEqual(freshContext, oldContext) {
		t.Fatal("Context sourcing is broken.")
	}

}

// TestContextStartBeat test concept behind heartbeat which takes care of resetting atomic rpm.
func TestContextStartBeat(t *testing.T) {
	var wg sync.WaitGroup

	c := &Context{
		rpm: 1,
		cache: SingleCache(),
	}

	var delay int = 1
	wg.Add(1)

	go func (){
		resetInterval := 2

		if delay < 60 {
			resetInterval = resetInterval - delay
		}

		<-time.After(time.Duration(resetInterval) * time.Second)
		c.rpm = 0
		wg.Done()

	}()

	wg.Wait()
	if c.rpm != 0 {
		t.Fatal("Heartbeat is broken.")
	}

}

// TestNewRoutes test creation of routes object.
func TestNewRoutes(t *testing.T) {
	r := NewRoutes()
	if r == nil {
		t.Fatal("Creation of routes object is broken.")
	}
}

// TestAPIRouter test creation of api router object with given context.
func TestAPIRouter(t *testing.T) {
	ctx := NewContext(false)
	router := APIRouter(ctx)

	if ctx == nil {
		t.Fatal("Creation of context is broken.")
	}

	if router == nil {
		t.Fatal("Creation of router is broken.")
	}

}
