package nanocube

import (
	"fmt"
	"testing"
)

func TestStruct(t *testing.T) {
	list := make([]*SpatNode, 4)
	fmt.Println(list)
	b := Bounds{1, 1, 1, 1}
	sn := SpatNode{Bounds: b, Children: make([]*SpatNode, 4)}
	fmt.Println(sn)
}
func TestNanoCubeSetUp(t *testing.T) {
	nb := SetUpCube(10, Bounds{1, -1, 10, 10}, []string{"Android", "iphone"})
	fmt.Println(nb.MaxLevel)
}

func TestAssignIndexOnBounds(t *testing.T) {
	b := Bounds{1, -1, 3, 3}
	obj := Object{1.5, 1.5, "A", 3}
	fmt.Println(AssignIndexOnBounds(obj, b))
}

func TestAddObj(t *testing.T) {
	nb := SetUpCube(3, Bounds{0, 0, 8, 8}, []string{"Android", "iPhone"})
	nb.AddObject(Object{3, -3, "Android", 50})
	// fmt.Println(nb.Root.Summary)
	// fmt.Println(nb.Root.Children[0])
	// fmt.Println(nb.Root.Children[0].Children[3].Summary)
	if (nb.Root.Summary != nb.Root.Children[0].Summary) || (nb.Root.Children[0].Summary != nb.Root.Children[0].Children[3].Summary) {
		t.Errorf("These three address should be equal")
	}

	nb.AddObject(Object{2.5, -3.5, "iPhone", 50})
	if (nb.Root.Summary != nb.Root.Children[0].Summary) || (nb.Root.Children[0].Summary != nb.Root.Children[0].Children[3].Summary) {
		t.Errorf("These three address should be equal")
	}
	nb.AddObject(Object{5, -5, "iPhone", 50})
	fmt.Println(nb.Root.Summary)
	fmt.Println(nb.Root.Children[0].Summary)
	fmt.Println(nb.Root.Children[3])
	fmt.Println(nb.Root.Children[3].Children[0].Summary)
	if (nb.Root.Summary == nb.Root.Children[0].Summary) || (nb.Root.Children[0].Summary == nb.Root.Children[3].Summary) {
		t.Errorf("These three address should not be equal")
	}
}
