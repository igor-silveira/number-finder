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

func NewHandler(finder service.FinderService, logger *slog.Logger) *Handler {
	return &Handler{
		finder: finder,
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/number/:number", h.handleFind)
	api.Get("/health", h.handleHealthCheck)
}

func (h *Handler) handleFind(c *fiber.Ctx) error {
	number := c.Params("number")
	thresholdStr := c.Query("thresholdPercentage", "0")

	h.logger.Debug("Received find request", "number", number)

	target, err := strconv.Atoi(number)

	if err != nil {
		h.logger.Error("Invalid number parameter", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid number parameter",
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

	h.logger.Debug("Find operation completed",
		"target", target,
		"result_index", result.Index,
		"result_value", result.Number,
		"is_approximate", result.IsApproximate,
	)

	return c.JSON(result)
}

func (h *Handler) handleHealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
