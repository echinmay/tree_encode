package main

import (
	_ "bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
 *
 * Sample code showing a way to encode a binary tree (not balanced and not complete) to a file and
 * then decode the file contents to get the tree back.
 *
 * For this sample each Node will have a key (of type int) and some data, a string here.
 *
 */

/*
 * The Data struct will carry the data
 */
type Data struct {
	Key int
	Val string
}

/*
 * The Tree data structure
 */
type Node struct {
	D     Data
	Left  *Node
	Right *Node
}

/*
 * We will declare a function type that we can use to walk through the tree and
 * perform any function.
 *
 * One function will be to walk through the tree and encode the data.
 * Another function can be to extract all the data into an array. The usage of
 * functions (closures) will help make the code generic
 *
 */
type parseTreeDataFunc func(root Data) error

/*
 * Generic Algorithm
 *
 * We will do a pre order traversal of the tree and encode each node as we do it into a file
 *
 * Then we will reconstruct the tree with that order.
 *
 * The workhorse of encoding the data in the tree is the encoding/gob function in golang.
 *
 * The data is encoded in a binary format (similar to a TLV format).
 *
 * The encoded file should be portable between different architectures (big endian and little endian).
 * Havent tested that functionality yet
 *
 */

/*
 * One way to do a pre-order traversal is to use recursion. Here we will use a stack so that
 * the call stack does not grow indefintely.
 */

// Define a generic stack data structure.
// Couple of points.
// 1. This is NOT re-entrant or thread safe
// 2. Calling Pop() on an empty stack will cause an exception
type stack []*Node

func (s stack) Push(v *Node) stack {
	s2 := append(s, v)
	return s2
}

// Warning : Will fail if called on an empty stack
func (s stack) Pop() (stack, *Node) {
	l := len(s)
	item := s[l-1]
	newstack := s[0 : l-1]
	return newstack, item
}

// Generic function to create a new Node with the given Data
func newnode(d Data) *Node {
	var n Node
	n.Left = nil
	n.Right = nil
	n.D = d
	return &n
}

/*
 * Simple function to add data to a tree.
 *
 * At each node the keys are compared and data is inserted to the left
 * or right tree.
 *
 * While recursion can be used a stack is used to guard against deep trees
 */
func addtotree(root *Node, d Data) *Node {
	newn := newnode(d)
	if root == nil {
		return newn
	}

	var parent *Node = root

	for {
		if d.Key < parent.D.Key {
			// Add to left tree
			if parent.Left != nil {
				parent = parent.Left
			} else {
				parent.Left = newn
				break
			}
		} else if d.Key > parent.D.Key {
			// Add to left tree
			if parent.Right != nil {
				parent = parent.Right
			} else {
				parent.Right = newn
				break
			}
		} else {
			// Cannot insert duplicate entries
			break
		}
	}
	return root
}

func processTree(f parseTreeDataFunc, root *Node) {
	if root == nil {
		return
	}
	nodestack := make(stack, 0)
	nodestack = nodestack.Push(root)
	var n *Node
	var e error
	for len(nodestack) > 0 {
		// Pop the topmost entry
		nodestack, n = nodestack.Pop()
		// Apply the given function
		e = f(n.D)
		if e != nil {
			// Break if there is a problem reported
			fmt.Println(e.Error())
			log.Fatal("Process data function return failure")
			return
		}
		// Check for right and left subtrees. Push the right and then the left
		if n.Right != nil {
			nodestack = nodestack.Push(n.Right)
		}
		if n.Left != nil {
			nodestack = nodestack.Push(n.Left)
		}
	}
	return

}

func printTreeFunc() parseTreeDataFunc {
	return func(d Data) error {
		fmt.Println("Current node : ", d.Key)
		return nil
	}
}

// Function that will encode the data to the file.
func encodeTreeFunc(f *os.File) parseTreeDataFunc {
	enc := gob.NewEncoder(f)
	return func(d Data) error {
		return (enc.Encode(d))
	}
}

// Function that will return (via the pointer) a slice of Data present in the
// tree
func getPreOrderDataVals(pptr **[]Data) parseTreeDataFunc {
	ds := make([]Data, 0)
	return func(d Data) error {
		ds = append(ds, d)
		*pptr = &ds
		return nil
	}
}

// Function to construct a tree given an array of Data values
func MakeTree(darr []Data) *Node {
	var root *Node = nil
	if darr == nil {
		return root
	}
	for idx := 0; idx < len(darr); idx++ {
		root = addtotree(root, darr[idx])
	}
	return root
}

// Helper function to compare two trees
func CompareTrees(tree1 *Node, tree2 *Node) bool {
	if tree1 == nil {
		if tree2 == nil {
			return true
		} else {
			return false
		}
	} else {
		if tree2 == nil {
			return false
		}
	}
	var data1 *[]Data
	collectDatafn := getPreOrderDataVals(&data1)
	processTree(collectDatafn, tree1)

	var data2 *[]Data
	collectDatafn2 := getPreOrderDataVals(&data2)
	processTree(collectDatafn2, tree2)

	arr1 := *data1
	arr2 := *data2

	if len(arr1) != len(arr2) {
		return false
	}

	for idx := 0; idx < len(arr1); idx++ {
		d1 := arr1[idx]
		d2 := arr2[idx]

		if d1.Key != d2.Key {
			return false
		}

		if strings.Compare(d1.Val, d2.Val) != 0 {
			return false
		}
	}
	return true
}

// Read a file and return an array of Data present in the file
func DecodeFile(f *os.File) []Data {
	dec := gob.NewDecoder(f)
	outputDataSlice := make([]Data, 0)
	for {
		var d Data
		err := dec.Decode(&d)
		if err != nil {
			break
		}
		outputDataSlice = append(outputDataSlice, d)
	}
	return outputDataSlice
}

// Encode a tree into a file
func EncodeIntoFile(f *os.File, root *Node) {
	encFn := encodeTreeFunc(f)
	processTree(encFn, root)
}

// This example shows the basic usage of the package: Create an encoder,
// save the tree to a local file and then decode it and reconstruct the tree
func main() {

	if len(os.Args) != 2 {
		fmt.Println(" Usage : ", os.Args[0], " <Name to file to encode in>")
		return
	}

	// Get a list of keys to insert
	var entries = [10]int{10, 7, 4, 1, 5, 9, 17, 15, 20, 30}

	// Prepare the data
	data := make([]Data, 0)
	for idx := 0; idx < len(entries); idx++ {
		s := strconv.Itoa(entries[idx])
		data = append(data, Data{entries[idx], s})
	}

	// Create the file in which the tree will be encoded
	f, e := os.Create(os.Args[1])
	if e != nil {
		log.Fatal("Could not create a file to store the tree in")
		return
	}

	// Make the tree from the data
	root := MakeTree(data)

	// Encode it into the file
	EncodeIntoFile(f, root)

	// Close the file
	f.Close()

	// Re-open the file
	f2, err2 := os.Open(os.Args[1])
	if err2 != nil {
		log.Fatal("Could not read the file back for reading")
		return
	}

	// Get the data encoded in the file
	outputDataSlice := DecodeFile(f2)
	defer f2.Close()

	// Recreate the tree
	newRoot := MakeTree(outputDataSlice)

	// Compare the original and new tree
	same := CompareTrees(root, newRoot)

	// Print the result
	fmt.Println(same)
}
