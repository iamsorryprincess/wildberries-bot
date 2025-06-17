package service

import (
	"context"

	"github.com/iamsorryprincess/wildberries-bot/cmd/api/model"
	"github.com/iamsorryprincess/wildberries-bot/internal/pkg/log"
)

type TrackingRepository interface {
	FindMatchTracking(ctx context.Context, categoryID uint64) ([]model.TrackingResult, error)
}

type NotificationSender interface {
	Send(ctx context.Context, message model.TrackingResult) error
}

type TrackingService struct {
	logger             log.Logger
	trackingRepository TrackingRepository
	notificationSender NotificationSender
}

func NewTrackingService(logger log.Logger, trackingRepository TrackingRepository, notificationSender NotificationSender) *TrackingService {
	return &TrackingService{
		logger:             logger,
		trackingRepository: trackingRepository,
		notificationSender: notificationSender,
	}
}

func (s *TrackingService) SendNotifications(ctx context.Context, categoryID uint64) error {
	trackingResults, err := s.trackingRepository.FindMatchTracking(ctx, categoryID)
	if err != nil {
		return err
	}

	if len(trackingResults) == 0 {
		return nil
	}

	for _, tracking := range trackingResults {
		if err = s.notificationSender.Send(ctx, tracking); err != nil {
			s.logger.Error().Err(err).
				Int64("chat_id", tracking.ChatID).
				Msg("failed send notification about tracking result")
		}
	}

	return nil
}
