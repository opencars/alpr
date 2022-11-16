package openalpr

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
	"github.com/opencars/seedwork/logger"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/recognizer"
)

// Recognizer is a pool of OpenALPR instances for recognizing the car plates.
type Recognizer struct {
	workers chan *openalpr.Alpr
}

// New returns new instance of vehicle plates recognizer.
func New(conf *config.OpenALPR) (*Recognizer, error) {
	workers := make(chan *openalpr.Alpr, conf.Pool)
	for i := 0; i < conf.Pool; i++ {
		alpr := openalpr.NewAlpr(conf.Country, conf.ConfigFile, conf.RuntimeDir)
		if !alpr.IsLoaded() {
			return nil, fmt.Errorf("failed to load alpr")
		}

		alpr.SetTopN(conf.MaxNumber)
		workers <- alpr
	}

	logger.Debugf("Using OpenALPR version: %s", openalpr.GetVersion())
	return &Recognizer{
		workers: workers,
	}, nil
}

// Recognize returns result of car plates recognition.
// Accepts io.Reader from JPEG.
func (r *Recognizer) Recognize(reader io.Reader) ([]recognizer.Result, error) {
	w := <-r.workers
	defer func() {
		r.workers <- w
	}()

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	res, err := w.RecognizeByBlob(content)
	if err != nil {
		return nil, err
	}

	return asRecognizerResult(&res), nil
}

func asRecognizerResult(in *openalpr.AlprResults) []recognizer.Result {
	out := make([]recognizer.Result, len(in.Plates))

	for i, plate := range in.Plates {
		out[i].Plate = plate.BestPlate

		//for _, candidate := range plate.TopNPlates {
		//	out[i].Candidates = append(out[i].Candidates, recognizer.Candidate{
		//		Confidence: candidate.OverallConfidence,
		//		Plate:      candidate.Characters,
		//	})
		//}

		for _, point := range plate.PlatePoints {
			out[i].Coordinates = append(out[i].Coordinates, recognizer.Coordinate{
				X: point.X,
				Y: point.Y,
			})
		}
	}

	return out
}
