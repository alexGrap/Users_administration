package handlers

import (
	models "avito/internal"
	"github.com/gofiber/fiber/v2"
)

type usecase interface {
	GetById(id string) ([]models.UserSubscription, error)
	CreateSegment(body models.SegmentBody) (models.SegmentBody, error)
	DeleteSegment(body models.SegmentBody) error
	GetSegments() ([]models.SegmentBody, error)
	Subscriber(body models.Subscriber) ([]models.UserSubscription, error)
	SubWithTime(body models.SubscribeWithTimeout) ([]models.UserSubscription, error)
}

type Handlers struct {
	usecase
}

func InitHandlers(service usecase) *Handlers {
	return &Handlers{service}
}

func (h *Handlers) GetById(ctx *fiber.Ctx) error {
	key := ctx.Query("id")
	body, err := h.usecase.GetById(key)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	return ctx.JSON(body)
}

func (h *Handlers) CreateSegment(ctx *fiber.Ctx) error {
	var body models.SegmentBody
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(503)
		return ctx.JSON("not valid segment body for parsing")
	}
	result, err := h.usecase.CreateSegment(body)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	ctx.Status(203)
	return ctx.JSON(result)
}

func (h *Handlers) DeleteSegment(ctx *fiber.Ctx) error {
	var body models.SegmentBody
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte("not valid segment body for parsing"))
	}
	err = h.usecase.DeleteSegment(body)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	ctx.Status(200)
	return ctx.Send([]byte("segment was successful deleted"))
}

func (h *Handlers) GetSegments(ctx *fiber.Ctx) error {
	result, err := h.usecase.GetSegments()
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	return ctx.JSON(result)
}

func (h *Handlers) Subscriber(ctx *fiber.Ctx) error {
	var body models.Subscriber
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(503)
		return ctx.JSON("not valid body")
	}
	result, err := h.usecase.Subscriber(body)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	ctx.Status(203)
	return ctx.JSON(result)
}

func (h *Handlers) SubscribeWithTimeOut(ctx *fiber.Ctx) error {
	var body models.SubscribeWithTimeout
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(503)
		return ctx.JSON("not valid body")
	}
	result, err := h.usecase.SubWithTime(body)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	ctx.Status(203)
	return ctx.JSON(result)
}
