package data_struct

type Person struct {
	Objects *Objects `json:"object"`
	Time    string   `json:"relative_time"`
}

type Objects struct {
	Pedestrian *Pedestrian `json:"pedestrian"`
}

type Pedestrian struct {
	Rectangle *Rectangle `json:"rectangle"`
	ID        string     `json:"track_id"`
}

type Rectangle struct {
	Vertices []Vertice `json:"vertices"`
}

type Vertice struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PairDistance struct {
	First_ID  string
	Second_ID string
	Distance  float64
}

type DistanceTimeStamp struct {
	Time          string
	PairDistances []PairDistance
}

type People map[string][]Person
