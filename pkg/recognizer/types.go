package recognizer

// Coordinate represents
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Candidate represents recognized car plate candidate.
type Candidate struct {
	Confidence float32 `json:"confidence"`
	Plate      string  `json:"plate"`
}

// Result represents result of recognition.
type Result struct {
	// Candidates  []Candidate  `json:"candidates"`
	Coordinates []Coordinate `json:"coordinates"`
	Plate       string       `json:"plate"`
}
