package grpc

import (
	"github.com/opencars/grpc/pkg/alpr"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/opencars/alpr/pkg/domain/model"
)

func convert(r *model.Recognition) *alpr.Recognition {
	return &alpr.Recognition{
		Id:        r.ID,
		ImageUrl:  r.ImageKey,
		Number:    r.Plate,
		CreatedAt: timestamppb.New(r.CreatedAt),
	}
}
