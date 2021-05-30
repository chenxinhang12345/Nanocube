package nanocube

import "fmt"

/*
Nanocube ...
*/
type Nanocube struct {
	Root     *SpatNode
	MaxLevel int //maximum level allowed for spatial attribute
	Types    []string
	Index    map[string]int //the map stores categorical index
}

//SpatNode for encoding spatial attribute
type SpatNode struct {
	Bounds   Bounds
	Children []*SpatNode
	Summary  *Summary
	CatRoot  *CatNode
	Level    int //current level
}

//CatNode for encoding categorical attribute
type CatNode struct {
	Children []*Summary
	Summary  *Summary
	// Type     string //the category
}

/*Bounds encode spatial information for each node
|---------Lng ++
|  0    1
|
|  2    3
Lat --
*/
type Bounds struct {
	Lng    float64
	Lat    float64
	Width  float64
	Height float64
}

//Object represent event
type Object struct {
	Lng       float64
	Lat       float64
	Type      string
	TimeStamp int
}

//Summary a summary for a bunch of nodes
type Summary struct {
	Count             int64
	TimeStampedCounts []int64
}

//SetUpCube Initialize the cube
func SetUpCube(MaxLevel int, MaxBounds Bounds, Types []string) *Nanocube {
	nc := &Nanocube{Root: &SpatNode{Bounds: MaxBounds, Children: make([]*SpatNode, 4), Level: 1}, MaxLevel: MaxLevel, Types: Types}
	m := make(map[string]int)
	for i := 0; i < len(Types); i++ {
		m[Types[i]] = i
	}
	nc.Index = m
	return nc
}

//AssignIndexOnBounds helper function for assigning index on specific bounds for an object
func AssignIndexOnBounds(obj Object, b Bounds) (int, Bounds) {
	HalfWidth := b.Width / 2
	HalfHeight := b.Height / 2
	MidLng := b.Lng + HalfWidth
	MidLat := b.Lat - HalfHeight
	// fmt.Println("func AssignIndexBounds ", obj, " ", b, "MidLng ", MidLng, "MidLat ", MidLat)
	if obj.Lng <= MidLng && obj.Lat >= MidLat {
		return 0, Bounds{b.Lng, b.Lat, HalfWidth, HalfHeight}
	} else if obj.Lng > MidLng && obj.Lat >= MidLat {
		return 1, Bounds{MidLng, b.Lat, HalfWidth, HalfHeight}
	} else if obj.Lng <= MidLng && obj.Lat < MidLat {
		return 2, Bounds{b.Lng, MidLat, HalfWidth, HalfHeight}
	} else if obj.Lng > MidLng && obj.Lat < MidLat {
		return 3, Bounds{MidLng, MidLat, HalfWidth, HalfHeight}
	} else {
		return 0, Bounds{}
	}
}

func (nc *Nanocube) getIndex(t string) int {
	return nc.Index[t]
}

//HasOnlyOneChild check if the SpatNode has only one child
func (s *SpatNode) HasOnlyOneChild() (bool, *SpatNode) {
	fmt.Println("debug hasonlyonechild:", s.Children)
	counter := 0
	var retptr *SpatNode = nil
	for i := 0; i < 4; i++ {
		if s.Children[i] != nil {
			retptr = s.Children[i]
			counter++
		}
	}
	fmt.Println(counter)
	return (counter == 1), retptr
}

//Copy return a deep copy of a summary
func (s *Summary) Copy() *Summary {
	return &Summary{Count: s.Count}
}

//UpdateSummary update current summary when adding an object to current SpatNode
// func (s *SpatNode) UpdateSummary(obj Object, maxLevel int) {
// 	hasOnlyOneChild, child := s.HasOnlyOneChild()
// 	if s.Level < maxLevel {
// 		if s.Summary == nil { //if it doesn't have summary
// 			s.Summary = child.Summary
// 		} else {
// 			if hasOnlyOneChild {
// 				s.Summary = child.Summary
// 			} else {
// 				s.Summary = s.Summary.Copy()
// 				s.Summary.Count++
// 			}
// 		}
// 	} else {
// 		if s.Summary == nil {
// 			s.Summary = &Summary{Count: 1}
// 		} else {
// 			s.Summary.Count++
// 		}
// 	}
// 	// fmt.Println(s)
// }

//HasOnlyOneChild check if the cat node has only one child
func (c *CatNode) HasOnlyOneChild() bool {
	counter := 0
	for i := 0; i < len(c.Children); i++ {
		if c.Children[i] != nil {
			counter++
		}
		if counter > 1 { //more than one child
			return false
		}
	}
	return true
}

//UpdateSummary update current summary when adding an object to current SpatNode
func (s *SpatNode) UpdateSummary(obj Object, maxLevel int, nc *Nanocube) {
	hasOnlyOneChild, child := s.HasOnlyOneChild()
	fmt.Println(hasOnlyOneChild)
	if s.Level < maxLevel {
		if s.CatRoot == nil { //if it doesn't have categorical root
			s.CatRoot = child.CatRoot
		} else { //if it has
			if hasOnlyOneChild { //only one child
				s.CatRoot = child.CatRoot
			} else { //need update
				fmt.Println("debug: ", child.CatRoot.Children)
				index := nc.getIndex(obj.Type) //update categorical node
				cpy := make([]*Summary, len(s.CatRoot.Children))
				copy(cpy, s.CatRoot.Children)
				s.CatRoot = &CatNode{Summary: s.CatRoot.Summary.Copy(), Children: cpy} //update cat root
				if s.CatRoot.Children[index] == nil {
					s.CatRoot.Children[index] = &Summary{Count: 1}
				} else {
					s.CatRoot.Children[index] = s.CatRoot.Children[index].Copy()
				}

				s.CatRoot.Summary.Count++

				s.CatRoot.Children[index].Count++
			}
		}
	} else {
		// fmt.Println("leave node")
		// fmt.Println("my cat ROOT:", s.CatRoot)
		if s.CatRoot == nil {
			// fmt.Println("no cat root")
			s.CatRoot = &CatNode{Summary: &Summary{Count: 1}, Children: make([]*Summary, len(nc.Types))}
			index := nc.getIndex(obj.Type)
			s.CatRoot.Children[index] = s.CatRoot.Summary
			// fmt.Println("my cat ROOT now:", s.CatRoot)
		} else { //need update
			index := nc.getIndex(obj.Type)
			// fmt.Println("leave children:", s.CatRoot.Children)
			if s.CatRoot.Children[index] != nil {
				s.CatRoot.Children[index].Count++
			} else {
				// fmt.Println("insert new type")
				s.CatRoot.Children[index] = &Summary{Count: 1}
				for i := 0; i < len(s.CatRoot.Children); i++ {
					if i != index {
						s.CatRoot.Children[i] = s.CatRoot.Children[i].Copy() //deep copy
					}
				}
				s.CatRoot.Summary.Count++
			}
		}
	}
	// fmt.Println(s)
}

//AddObject Add an object
func (nc *Nanocube) AddObject(obj Object) {
	stack := make([]*SpatNode, 0)
	levels := nc.MaxLevel
	currentNode := nc.Root
	currentLevel := 1
	for currentLevel < levels {
		// fmt.Println("currentLevel: ", currentLevel)
		index, b := AssignIndexOnBounds(obj, currentNode.Bounds)
		// fmt.Println("Assignindex: ", index)
		if currentNode.Children[index] == nil { //no nodes on current index
			currentNode.Children[index] = &SpatNode{Bounds: b, Children: make([]*SpatNode, 4)} //create a new node on current index
		}
		currentNode.Level = currentLevel
		stack = append(stack, currentNode)
		currentNode = currentNode.Children[index] //next level node
		currentLevel++
	}
	currentNode.Level = currentLevel
	// fmt.Println("leave level:", currentLevel)
	stack = append(stack, currentNode) //update leaves
	for i := len(stack) - 1; i >= 0; i-- {
		currentNode = stack[i]
		currentNode.UpdateSummary(obj, levels, nc)
	}
}
