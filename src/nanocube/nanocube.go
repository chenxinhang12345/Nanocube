package nanocube

/*
Nanocube ...
*/
type Nanocube struct {
	Root     *SpatNode
	MaxLevel int //maximum level allowed for spatial attribute
	Types    []string
}

//SpatNode for encoding spatial attribute
type SpatNode struct {
	Bounds   Bounds
	Children []*SpatNode
	Summary  *Summary
	Level    int //current level
}

//CatNode for encoding categorical attribute
type CatNode struct {
	Children []CatNode
	Summary  *Summary
	Type     string //the category
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
	return &Nanocube{&SpatNode{Bounds: MaxBounds, Children: make([]*SpatNode, 4), Level: 1}, MaxLevel, Types}
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

//HasOnlyOneChild check if the SpatNode has only one child
func (s *SpatNode) HasOnlyOneChild() (bool, *SpatNode) {
	counter := 0
	var retptr *SpatNode = nil
	for i := 0; i < 4; i++ {
		if s.Children[i] != nil {
			retptr = s.Children[i]
			counter++
		}
	}
	return (counter == 1), retptr
}

//Copy return a deep copy of a summary
func (s *Summary) Copy() *Summary {
	return &Summary{Count: s.Count}
}

//UpdateSummary update current summary when adding an object to current SpatNode
func (s *SpatNode) UpdateSummary(obj Object, maxLevel int) {
	hasOnlyOneChild, child := s.HasOnlyOneChild()
	if s.Level < maxLevel {
		if s.Summary == nil { //if it doesn't have summary
			s.Summary = child.Summary
		} else {
			if hasOnlyOneChild {
				s.Summary = child.Summary
			} else {
				s.Summary = s.Summary.Copy()
				s.Summary.Count++
			}
		}
	} else {
		if s.Summary == nil {
			s.Summary = &Summary{Count: 1}
		} else {
			s.Summary.Count++
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
	stack = append(stack, currentNode) //update leaves
	for i := len(stack) - 1; i >= 0; i-- {
		currentNode = stack[i]
		currentNode.UpdateSummary(obj, levels)
	}
}
