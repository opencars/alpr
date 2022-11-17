package model

// Coordinate represents
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Result represents result of recognition.
type Result struct {
	Coordinates []Coordinate `json:"coordinates"`
	Plate       string       `json:"plate"`
}
