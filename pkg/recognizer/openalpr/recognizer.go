package openalpr

import (
	"github.com/opencars/alpr/pkg/recognizer"
	"io"
	"io/ioutil"
	"sync"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/logger"
)

// Recognizer is a pool of OpenALPR instances for recognizing the car plates.
type Recognizer struct {
	pool *sync.Pool
}

// New returns new instance of vehicle plates recognizer.
func New(conf *config.OpenALPR) (*Recognizer, error) {
	var pool = sync.Pool{
		New: func() interface{} {
			alpr := openalpr.NewAlpr(conf.Country, conf.ConfigFile, conf.RuntimeDir)
			if !alpr.IsLoaded() {
				return nil
			}

			alpr.SetTopN(conf.MaxNumber)
			return alpr
		},
	}

	logger.Debug("Using OpenALPR version: %s", openalpr.GetVersion())
	return &Recognizer{
		pool: &pool,
	}, nil
}

// Recognize returns result of car plates recognition.
// Accepts io.Reader from JPEG.
func (r *Recognizer) Recognize(reader io.Reader) ([]recognizer.Result, error) {
	alpr := r.pool.Get().(*openalpr.Alpr)
	if alpr == nil {
		return nil, recognizer.ErrFailedToLoad
	}

	defer r.pool.Put(alpr)

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	res, err := alpr.RecognizeByBlob(content)
	if err != nil {
		return nil, err
	}

	return asRecognizerResult(&res), nil
}

func asRecognizerResult(in *openalpr.AlprResults) []recognizer.Result {
	out := make([]recognizer.Result, len(in.Plates))

	for i, plate := range in.Plates {
		out[i].Plate = plate.BestPlate

		for _, candidate := range plate.TopNPlates {
			out[i].Candidates = append(out[i].Candidates, recognizer.Candidate{
				Confidence: candidate.OverallConfidence,
				Plate:      candidate.Characters,
			})
		}

		for _, point := range plate.PlatePoints {
			out[i].Coordinates = append(out[i].Coordinates, recognizer.Coordinate{
				X: point.X,
				Y: point.Y,
			})
		}
	}

	return out
}
