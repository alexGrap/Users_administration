package handlers

import (
	models "avito/internal"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type usecase interface {
	GetById(id string) ([]models.UserSubscription, error)
	CreateSegment(body models.SegmentBody) (models.SegmentBody, error)
	DeleteSegment(body models.SegmentBody) error
	GetSegments() ([]models.SegmentBody, error)
	Subscriber(body models.Subscriber) ([]models.UserSubscription, error)
	SubWithTime(body models.SubscribeWithTimeout) ([]models.UserSubscription, error)
	History(userId int64, from string, to string) (string, error)
}

type Handlers struct {
	usecase
}

func InitHandlers(service usecase) *Handlers {
	return &Handlers{service}
}

// GetById godoc
// @Description Get list of user`s subscriptions.
// @Summary Get subscriptions
// @Tags User
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 "OK" models.UserSubscription
// @Failure 404 "Bad request"
// @Failure 503 "Not found"
// @Router /getById [get]
func (h *Handlers) GetById(ctx *fiber.Ctx) error {
	key := ctx.Query("id")
	if key == "" {
		ctx.SendStatus(404)
	}
	body, err := h.usecase.GetById(key)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	return ctx.JSON(body)
}

// CreateSegment godoc
// @Description Create a segment. You need to post name of segment and if this need percent of user
// @Summary Create a new segment
// @Tags Segment
// @Accept json
// @Produce json
// @Param input body models.SegmentBody true "Body of creation segment"
// @Success 204 "OK" models.SegmentBody
// @Failure 503 "Bad request"
// @Router /createSegment [post]
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

// DeleteSegment godoc
// @Description Delete a segment. You need to post name of segment
// @Summary Delete a segment
// @Tags Segment
// @Accept json
// @Produce json
// @Param input body models.SegmentBody true "Name of segment"
// @Success 200 "OK"
// @Failure 503 "Bad request"
// @Router /deleteSegment [delete]
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

// GetSegments godoc
// @Description Get a segment.
// @Summary Get all existed segments
// @Tags Segment
// @Produce json
// @Success 204 "OK" models.SegmentBody
// @Failure 503 "Bad request"
// @Router /getSegment [get]
func (h *Handlers) GetSegments(ctx *fiber.Ctx) error {
	result, err := h.usecase.GetSegments()
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	return ctx.JSON(result)
}

// Subscriber godoc
// @Description Create new and delete existed subscription for user
// @Summary Make a subscription
// @Tags User
// @Accept json
// @Produce json
// @Param input body models.Subscriber true "Body of subscription"
// @Success 203 "OK" models.UserSubscription
// @Failure 503 "Bad request"
// @Router /subscription [put]
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

// SubscribeWithTimeOut godoc
// @Description Create new subscription for user with ttl
// @Summary Make a subscription with timeout
// @Tags User
// @Accept json
// @Produce json
// @Param input body models.SubscribeWithTimeout true "Body of subscription"
// @Success 203 "OK" models.UserSubscription
// @Failure 503 "Bad request"
// @Router /timeoutSubscribe [put]
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

// History godoc
// @Description Get user subscriptions history
// @Summary Get history
// @Tags History
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param from query string true "Time start. ex: 2020-03-20"
// @Param to query string true "Time end. ex: 2024-03-20"
// @Success 203 "OK" models.UserSubscription
// @Failure 400 "Bad request"
// @Router /timeoutSubscribe [put]
func (h *Handlers) History(ctx *fiber.Ctx) error {
	userId, err := strconv.ParseInt(ctx.Query("userId"), 10, 64)
	if err != nil {
		ctx.Status(400)
		return ctx.Send([]byte(err.Error()))
	}
	from := ctx.Query("from")
	to := ctx.Query("to")
	result, err := h.usecase.History(userId, from, to)
	if err != nil {
		ctx.Status(503)
		return ctx.Send([]byte(err.Error()))
	}
	ctx.Status(200)
	return ctx.SendFile(result)
}
