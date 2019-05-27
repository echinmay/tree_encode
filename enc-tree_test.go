package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

// Helper functions
func getTestData1(e *[]int) []Data {
	data := make([]Data, 0)
	entries := *e
	for idx := 0; idx < len(entries); idx++ {
		s := strconv.Itoa(entries[idx])
		data = append(data, Data{entries[idx], s})
	}
	return data
}

func TestWrongRootNodeCmpFn(t *testing.T) {
	var entries = [10]int{10, 7, 4, 1, 5, 9, 17, 15, 20, 30}
	data := make([]Data, 0)
	for idx := 0; idx < len(entries); idx++ {
		s := strconv.Itoa(entries[idx])
		data = append(data, Data{entries[idx], s})
	}
	var entries2 = [10]int{11, 7, 4, 1, 5, 9, 17, 15, 20, 30}
	data2 := make([]Data, 0)
	for idx := 0; idx < len(entries2); idx++ {
		s := strconv.Itoa(entries2[idx])
		data = append(data, Data{entries2[idx], s})
	}
	tree1 := MakeTree(data)
	tree2 := MakeTree(data2)
	cmp := CompareTrees(tree1, tree2)
	if cmp != false {
		t.Errorf("The root node is different in the two trees. Should have returned false")
	}
}

func TestWrongLeafNode(t *testing.T) {
	var entries = [10]int{10, 7, 4, 1, 5, 9, 17, 15, 20, 30}
	data := make([]Data, 0)
	for idx := 0; idx < len(entries); idx++ {
		s := strconv.Itoa(entries[idx])
		data = append(data, Data{entries[idx], s})
	}
	var entries2 = [10]int{10, 7, 4, 2, 5, 9, 17, 15, 20, 30}
	data2 := make([]Data, 0)
	for idx := 0; idx < len(entries2); idx++ {
		s := strconv.Itoa(entries2[idx])
		data2 = append(data2, Data{entries2[idx], s})
	}
	tree1 := MakeTree(data)
	tree2 := MakeTree(data2)
	cmp := CompareTrees(tree1, tree2)
	if cmp != false {
		t.Errorf("The root node is different in the two trees. Should have returned false")
	}
}

func TestEmptyTreeCmp(t *testing.T) {
	//var entries = [10]int{10, 7, 4, 1, 5, 9, 17, 15, 20, 30}
	data := make([]Data, 0)
	data2 := make([]Data, 0)
	tree1 := MakeTree(data)
	tree2 := MakeTree(data2)
	cmp := CompareTrees(tree1, tree2)
	if cmp != true {
		t.Errorf("Empty trees should be compared to true ")
	}
}

func testReadingFromDifferentFile(t *testing.T) {
	var entries = [10]int{10, 7, 4, 1, 5, 9, 17, 15, 20, 30}
	data := make([]Data, 0)
	for idx := 0; idx < len(entries); idx++ {
		s := strconv.Itoa(entries[idx])
		data = append(data, Data{entries[idx], s})
	}
	var entries2 = [10]int{10, 7, 4, 2, 5, 9, 17, 15, 20, 30}
	data2 := make([]Data, 0)
	for idx := 0; idx < len(entries2); idx++ {
		s := strconv.Itoa(entries2[idx])
		data2 = append(data2, Data{entries2[idx], s})
	}
	tree1 := MakeTree(data)
	tree2 := MakeTree(data2)
	f1, err := os.Create("/tmp/testing-1")
	if err != nil {
		fmt.Println("Error creating file")
	}
	EncodeIntoFile(f1, tree1)
	f1.Close()
	f2, err := os.Create("/tmp/testing-2")
	if err != nil {
		fmt.Println("Error creating file")
	}
	EncodeIntoFile(f2, tree2)
	f2.Close()
	f2, err = os.Open("/tmp/testing-2")
	if err != nil {
		fmt.Println("Error creating file")
	}
	testData := DecodeFile(f2)
	f2.Close()
	tree3 := MakeTree(testData)
	cmp1 := CompareTrees(tree3, tree1)
	if cmp1 != false {
		t.Errorf("Expected a comparison fail")
	}
	cmp2 := CompareTrees(tree3, tree2)
	if cmp2 != true {
		t.Errorf("Expected a comparison pass")
	}

}

func TestVeryDeepTrees(t *testing.T) {
	entries := make([]int, 0)
	for idx := 0; idx < 10000; idx++ {
		entries = append(entries, idx)
	}

	data := make([]Data, 0)
	for idx := 0; idx < len(entries); idx++ {
		s := strconv.Itoa(entries[idx])
		data = append(data, Data{entries[idx], s})
	}
	tree := MakeTree(data)
	f, err := os.Create("/tmp/testing-big-file-1")
	if err != nil {
		fmt.Println("Error creating file")
	}
	EncodeIntoFile(f, tree)
	f.Close()
	f1, err1 := os.Open("/tmp/testing-big-file-1")
	if err1 != nil {
		fmt.Println("Error creating file")
	}
	data2 := DecodeFile(f1)
	tree2 := MakeTree(data2)
	f1.Close()
	cmp := CompareTrees(tree, tree2)
	if cmp != true {
		t.Errorf("Even big trees should be correctly encoded")
	}
}

func TestVeryVeryDeepTrees(t *testing.T) {
	entries := make([]int, 0)
	for idx := 0; idx < 100000; idx++ {
		entries = append(entries, idx)
	}

	data := make([]Data, 0)
	for idx := 0; idx < len(entries); idx++ {
		s := strconv.Itoa(entries[idx])
		data = append(data, Data{entries[idx], s})
	}
	tree := MakeTree(data)
	f, err := os.Create("/tmp/testing-verybig-file-1")
	if err != nil {
		fmt.Println("Error creating file")
	}
	EncodeIntoFile(f, tree)
	f.Close()
	f1, err1 := os.Open("/tmp/testing-verybig-file-1")
	if err1 != nil {
		fmt.Println("Error creating file")
	}
	data2 := DecodeFile(f1)
	tree2 := MakeTree(data2)
	f1.Close()
	cmp := CompareTrees(tree, tree2)
	if cmp != true {
		t.Errorf("Even big trees should be correctly encoded")
	}
}
