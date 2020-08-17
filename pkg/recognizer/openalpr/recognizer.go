package openalpr

import (
	"errors"
	"io"
	"io/ioutil"
	"sync"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"

	"github.com/opencars/alpr/pkg/config"
	"github.com/opencars/alpr/pkg/logger"
)

type Recognizer struct {
	pool *sync.Pool
}

func New(conf *config.OpenALPR) (*Recognizer, error) {
	var pool = sync.Pool{
		New: func() interface{} {
			alpr := openalpr.NewAlpr(conf.Country, conf.ConfigFile, conf.RuntimeDir)
			if !alpr.IsLoaded() {
				return errors.New("openalpr: failed to load")
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

func (r *Recognizer) Recognize(reader io.Reader) (*openalpr.AlprResults, error) {
	alpr := r.pool.Get().(*openalpr.Alpr)
	defer r.pool.Put(alpr)

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	res, err := alpr.RecognizeByBlob(content)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
