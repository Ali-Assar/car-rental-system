package api

import (
	"github.com/Ali-Assar/car-rental-system/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type ReservationHandler struct {
	store *db.Store
}

func NewReservationHandler(store *db.Store) *ReservationHandler {
	return &ReservationHandler{
		store: store,
	}
}

func (h *ReservationHandler) HandleCancelReservation(c *fiber.Ctx) error {
	id := c.Params("id")
	reservation, err := h.store.Reservation.GetReservationByID(c.Context(), id)
	if err != nil {
		return ErrNotFound("reservation")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrAuthorization()
	}
	if reservation.UserID != user.ID {
		return ErrAuthorization()
	}
	if err := h.store.Reservation.UpdateReservation(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})
}
func (h *ReservationHandler) HandleGetReservations(c *fiber.Ctx) error {
	reservation, err := h.store.Reservation.GetReservation(c.Context(), bson.M{})
	if err != nil {
		return ErrNotFound("reservations")
	}
	return c.JSON(reservation)
}

func (h *ReservationHandler) HandleGetReservation(c *fiber.Ctx) error {
	id := c.Params("id")
	reservation, err := h.store.Reservation.GetReservationByID(c.Context(), id)
	if err != nil {
		return ErrNotFound("reservation")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrAuthorization()
	}

	if reservation.UserID != user.ID {
		return ErrAuthorization()
	}
	return c.JSON(reservation)
}
