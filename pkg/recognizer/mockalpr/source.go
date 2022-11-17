package mockalpr

//go:generate mockgen -destination=./mockalpr.go -package=mockalpr github.com/opencars/alpr/pkg/recognizer Recognizer
