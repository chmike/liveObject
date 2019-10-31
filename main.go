package main

import (
	"fmt"
	"runtime"
	"time"
)

// T is a live object type with an embedded pointer to the tObj value.
type T struct {
	*tObj
	a string // field not accessed in the goroutine
}

// tObj holds the flag monitored by the goroutine associated to T and all member
// variables of T accessed by the goroutine. It’s required that the goroutine
// does not reference the pointer to T, only a pointer to tOblj.
type tObj struct {
	done bool
	b    string // field accessed in the goroutine
}

// NewT instantiate a live T object.
func NewT() *T {
	t := &T{tObj: &tObj{b: "I’m T"}, a: "hello world"}

	// the finalizer sets done to true when T is garbage collected
	runtime.SetFinalizer(t, func(o interface{}) { o.(*T).tObj.done = true })

	// the goroutine associated to T monitors done and terminate when true. It must not
	// hold a pointer to T, only to tObj. So never use t in the goroutine, only o.
	go func(o *tObj) {
		for !o.done {
			time.Sleep(1 * time.Second)
			fmt.Println(o.b, ": I’m alive")
		}
		fmt.Println(o.b, ": I terminate")
	}(t.tObj)

	// return the pointer to struct with an associated goroutine (live object)
	return t
}

func main() {
	t := NewT()
	// note how you can directly access the fields of tObj due to embedding
	fmt.Println("main : t is instantiated:", t.a, ":", t.b)

	fmt.Println("main : wait 2 seconds")
	time.Sleep(2 * time.Second)

	fmt.Println("main : set t to nil")
	t = nil

	fmt.Println("main : run GC")
	runtime.GC()

	fmt.Println("main : wait 2 seconds")
	time.Sleep(2 * time.Second)

	fmt.Println("main : terminate")
}