package main

import "testing"


func TestFmtId(t *testing.T) {
	// errc := make(chan error, 1)
	// list := List{}
	// t.Errorf("%T\n", item.FmtId)
	// err := Enclose(CombineStrings, "Im", "Chad")
	// t.Errorf("Enclose %v\n", err())
	// if err != nil {
	// 	t.Errorf("Failed because %v\n", err())
	// }

}

// TESTING FOR InSynExecute(errc chan error, ...args)
// //////////////////PROBLEM///////////////////
	// Arguments being passed Enclose() are not changing as each closure is executed.
// Each method needs to be made into a closure and then executed before the next method starts the same process
// Might be able to use a function that uses a channel to send each closure to a select stmt
