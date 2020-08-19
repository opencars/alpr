package recognizer

// Coordinate ...
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Candidate ...
type Candidate struct {
	Confidence float32 `json:"confidence"`
	Plate      string  `json:"plate"`
}

// Result ...
type Result struct {
	Coordinates []Coordinate `json:"coordinates"`
	Candidates  []Candidate  `json:"candidates"`
	Plate       string       `json:"plate"`
}
