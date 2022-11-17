package grpc

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/opencars/seedwork"
)

var (
	ErrValidationFailed = status.New(codes.InvalidArgument, "request.validation_failed")
)

func handleErr(err error) error {
	var vErr seedwork.ValidationError

	if errors.As(err, &vErr) {
		br := errdetails.BadRequest{
			FieldViolations: make([]*errdetails.BadRequest_FieldViolation, 0),
		}

		for k, vv := range vErr.Messages {
			for _, v := range vv {
				br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
					Field:       k,
					Description: v,
				})
			}
		}

		st, err := ErrValidationFailed.WithDetails(&br)
		if err != nil {
			panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
		}

		return st.Err()
	}

	var e seedwork.Error
	if errors.As(err, &e) {
		return status.Error(codes.FailedPrecondition, e.Error())
	}

	return err
}
