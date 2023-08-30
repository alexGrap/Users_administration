package usecase

import (
	models "avito/internal"
	"errors"
	"strconv"
	"time"
)

type repository interface {
	GetSubs(id int64) ([]models.UserSubscription, error)
	CreateSegment(body models.SegmentBody) (models.SegmentBody, error)
	DeleteSegment(segmentName string) error
	GetSegments() ([]models.SegmentBody, error)
	Subscriber(subscriber models.Subscriber) error
	SubWIthTimeout(id int64, name string, timeToDie time.Time) error
	TimeOutDeleter(sec chan int)
}

type UseCase struct {
	repository
}

func InitUsecase(rep repository) *UseCase {

	return &UseCase{rep}
}

func (useCase *UseCase) GetById(key string) ([]models.UserSubscription, error) {
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return []models.UserSubscription{}, errors.New("not valid user id")
	}
	return useCase.repository.GetSubs(id)
}

func (useCase *UseCase) CreateSegment(body models.SegmentBody) (models.SegmentBody, error) {
	if body.Name == "" {
		return models.SegmentBody{}, errors.New("empty segment name")
	}
	if body.Percent > 100 || body.Percent < 0 {
		return models.SegmentBody{}, errors.New("non valid users percent")
	}
	return useCase.repository.CreateSegment(body)
}

func (useCase *UseCase) DeleteSegment(body models.SegmentBody) error {
	if body.Name == "" {
		return errors.New("empty segment name")
	}
	return useCase.repository.DeleteSegment(body.Name)
}

func (useCase *UseCase) GetSegments() ([]models.SegmentBody, error) {
	result, err := useCase.repository.GetSegments()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (useCase *UseCase) Subscriber(body models.Subscriber) ([]models.UserSubscription, error) {
	err := useCase.repository.Subscriber(body)
	if err != nil {
		return []models.UserSubscription{}, err
	}
	return useCase.repository.GetSubs(body.UserId)
}

func (useCase *UseCase) SubWithTime(body models.SubscribeWithTimeout) ([]models.UserSubscription, error) {
	if body.TimeOut < 0 {
		return []models.UserSubscription{}, errors.New("not valid count of subscription day")
	}
	currentTime := time.Now()
	newTime := currentTime.Add(24 * time.Hour * time.Duration(body.TimeOut))
	err := useCase.repository.SubWIthTimeout(body.UserId, body.SegmentName, newTime)
	if err != nil {
		return []models.UserSubscription{}, err
	}
	return useCase.repository.GetSubs(body.UserId)
}
