// Bit improved from the gotour version.  Handles arbitrary depth trees.
package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	WalkRecursive(t, ch)
	close(ch)
}

func WalkRecursive(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		WalkRecursive(t.Left, ch)
	}

	ch <- t.Value

	if t.Right != nil {
		WalkRecursive(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		a, ok1 := <-ch1
		b, ok2 := <-ch2

		if ok1 == ok2 && !ok1 {
			return true
		}
		if ok1 != ok2 {
			return false
		}

		fmt.Printf("%v vs %v\n", a, b)
		if a != b {
			return false
		}
	}
}

func Print(ch chan int) {

	for {
		a := <-ch
		fmt.Printf("%v", a)
	}
}

func main() {
	fmt.Printf("Were same == %v\n", Same(tree.New(1), tree.New(1)))
	fmt.Printf("Were same == %v\n", Same(tree.New(1), tree.New(2)))
}
