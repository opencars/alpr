package mockstore

//go:generate mockgen -destination=./store.go -package=mockstore github.com/opencars/alpr/pkg/store Store,RecognitionRepository
