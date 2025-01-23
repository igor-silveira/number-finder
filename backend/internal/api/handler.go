package api

import (
	"log/slog"
	"number-finder-api/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	finder service.FinderService
	logger *slog.Logger
}

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Index         int  `json:"index"`
	Value         int  `json:"value"`
	IsApproximate bool `json:"is_approximate"`
}

func NewHandler(finder service.FinderService, logger *slog.Logger) *Handler {
	return &Handler{
		finder: finder,
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/number/:value", h.handleFind)
	api.Get("/health", h.handleHealthCheck)
}

func (h *Handler) handleFind(c *fiber.Ctx) error {
	value := c.Params("value")
	thresholdStr := c.Query("thresholdPercentage", "0")

	h.logger.Debug("Received find request", "value", value)

	target, err := strconv.Atoi(value)

	if err != nil {
		h.logger.Error("Invalid value parameter", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid value parameter",
		})
	}

	thresholdPercentage, err := strconv.ParseFloat(thresholdStr, 64)

	if err != nil {
		h.logger.Error("Invalid thresholdPercentage parameter", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid thresholdPercentage parameter",
		})
	}

	result, err := h.finder.Find(target, thresholdPercentage)

	if err != nil {
		h.logger.Debug(
			"Find operation error",
			"target", target,
			"thresholdPercentage", thresholdPercentage,
			"error", err,
		)

		errResponse := Error{Message: err.Error()}
		return c.Status(fiber.StatusNotFound).JSON(errResponse)
	}

	response := Response{
		Index:         result.Index,
		Value:         result.Value,
		IsApproximate: result.IsApproximate,
	}

	h.logger.Debug("Find operation completed",
		"target", target,
		"result_index", result.Index,
		"result_value", result.Value,
		"is_approximate", result.IsApproximate,
	)

	return c.JSON(response)
}

func (h *Handler) handleHealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
