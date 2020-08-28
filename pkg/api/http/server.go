package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/opencars/alpr/pkg/handler"
	"github.com/opencars/alpr/pkg/logger"
	"github.com/opencars/alpr/pkg/objectstore"
	"github.com/opencars/alpr/pkg/recognizer"
	"github.com/opencars/alpr/pkg/version"
)

type server struct {
	router      *mux.Router
	client      *http.Client
	recognizer  recognizer.Recognizer
	objectStore objectstore.ObjectStore
}

func newServer(recognizer recognizer.Recognizer, objectStore objectstore.ObjectStore) *server {
	s := server{
		router:      mux.NewRouter(),
		recognizer:  recognizer,
		objectStore: objectStore,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	s.configureRouter()

	return &s
}

func (s *server) configureRouter() {
	router := s.router.Methods("GET", "OPTIONS").Subrouter()

	router.Handle("/api/v1/alpr/public/version", s.Version())
	router.Handle("/api/v1/alpr/private/recognize", s.Recognize())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Api-Key", "X-Api-Key"})

	cors := handlers.CORS(origins, methods, headers)(s.router)
	compress := handlers.CompressHandler(cors)
	compress.ServeHTTP(w, r)
}

func (*server) Swagger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yml")
	}
}

func (*server) Version() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		v := struct {
			Version string `json:"version"`
			Go      string `json:"go"`
		}{
			Version: version.Version,
			Go:      runtime.Version(),
		}

		return json.NewEncoder(w).Encode(v)
	}
}

// Note: Later we could publish event of recognition into the NATS queue and prepare.
func (s *server) Recognize() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		imageURL := r.URL.Query().Get("image_url")
		if imageURL == "" {
			return handler.ErrRequiredImageURL
		}

		_, err := url.ParseRequestURI(imageURL)
		if err != nil {
			return handler.ErrInvalidImageURL
		}

		resp, err := s.client.Get(imageURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var buf bytes.Buffer
		tee := io.TeeReader(resp.Body, &buf)

		res, err := s.recognizer.Recognize(tee)
		if err != nil {
			return err
		}

		if len(res) > 0 && s.objectStore != nil {
			ext := path.Ext(imageURL)
			err := s.objectStore.Put(r.Context(), res[0].Plate+ext, &buf)
			if err != nil {
				logger.Errorf("failed to put: %v", err)
			}
		}

		return json.NewEncoder(w).Encode(res)
	}
}
