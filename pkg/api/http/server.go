package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/opencars/httputil"

	"github.com/opencars/alpr/pkg/domain"
	"github.com/opencars/alpr/pkg/domain/command"
)

const (
	// MaxImageSize equals to 5 MB.
	MaxImageSize = 5 << 20

	// ClientTimeOut equals to 10 seconds.
	ClientTimeOut = 10 * time.Second
)

type server struct {
	router *mux.Router
	svc    domain.CustomerService
}

func newServer(svc domain.CustomerService) *server {
	s := server{
		router: mux.NewRouter(),
		svc:    svc,
	}

	s.configureRouter()

	return &s
}

func (s *server) Recognize() httputil.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		c := command.Recognize{
			URL: r.URL.Query().Get("image_url"),
		}

		result, err := s.svc.Recognize(r.Context(), &c)
		if err != nil {
			return handleErr(err)
		}

		return json.NewEncoder(w).Encode(result)
	}
}
