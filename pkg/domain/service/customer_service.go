package service

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"

	"github.com/opencars/seedwork"
	"github.com/opencars/seedwork/logger"

	"github.com/opencars/alpr/pkg/domain"
	"github.com/opencars/alpr/pkg/domain/command"
	"github.com/opencars/alpr/pkg/domain/model"
)

type CustomerService struct {
	client     *http.Client
	recognizer domain.Recognizer
	obj        domain.ObjectStore
	pub        domain.Publisher
}

func NewCustomerService(rec domain.Recognizer, pub domain.Publisher) *CustomerService {
	httpClient := http.Client{
		Timeout: model.ClientTimeOut,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: model.ClientTimeOut,
			}).DialContext,
		},
	}

	return &CustomerService{
		client:     &httpClient,
		recognizer: rec,
		pub:        pub,
	}
}

func (s *CustomerService) Recognize(ctx context.Context, c *command.Recognize) ([]model.Result, error) {
	if err := seedwork.ProcessCommand(c); err != nil {
		return nil, err
	}

	resp, err := s.client.Get(c.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyWithLimit := io.LimitReader(resp.Body, model.MaxImageSize+1)

	var buff bytes.Buffer
	if _, err = io.CopyN(&buff, bodyWithLimit, bytes.MinRead); err != nil {
		return nil, err
	}

	typ := http.DetectContentType(buff.Bytes())
	if typ != "image/jpeg" {
		return nil, model.ErrUnknownContentType
	}

	_, err = buff.ReadFrom(bodyWithLimit)
	if err != nil {
		return nil, err
	}

	if buff.Len() > model.MaxImageSize {
		return nil, model.ErrImageTooLarge
	}

	reader := bytes.NewReader(buff.Bytes())

	res, err := s.recognizer.Recognize(reader)
	if err != nil {
		return nil, err
	}

	if len(res) > 0 {
		err := s.pub.Publish(&model.Event{
			URL:    c.URL,
			Number: res[0].Plate,
		})

		if err != nil {
			logger.Errorf("publish: %v", err)
		}
	}

	return res, nil
}
