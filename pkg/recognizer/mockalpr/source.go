package mockalpr

//go:generate mockgen -destination=./store.go -package=mockalpr github.com/opencars/alpr/pkg/recognizer Recognizer
