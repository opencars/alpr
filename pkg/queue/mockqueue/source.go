package mockqueue

//go:generate mockgen -destination=./queue.go -package=mockqueue github.com/opencars/alpr/pkg/queue Publisher,Subscriber
