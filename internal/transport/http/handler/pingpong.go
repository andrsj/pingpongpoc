package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	pingpongErrors "pingpongpoc/internal/errors"
	"pingpongpoc/internal/transport/responses"
)

type (
	Service interface {
		Serve(context.Context, string) (string, error)
	}

	// Handler holds the use-case to be executed
	Handler struct {
		service Service
		log     *slog.Logger
	}
)

func New(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     logger,
	}
}

// Ping handles the endpoint
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Handling '/' endpoint")

	input, ctx := r.URL.Query().Get("input"), r.Context()

	output, err := h.service.Serve(ctx, input)
	statusCode, resp := http.StatusOK, responses.Response{} //nolint:exhaustruct

	if err != nil {
		h.log.Info("Serving error", "error", err)

		resp.Error, statusCode = responses.InternalServerErrorMessage, http.StatusInternalServerError
		if errors.Is(err, pingpongErrors.ErrNotValidInput) {
			resp.Error, statusCode = responses.InvalidUserResponse, http.StatusBadRequest
		}
	} else {
		resp.Message = fmt.Sprintf("Response: %s", output)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	jsonResp, err := json.Marshal(resp) //nolint:errchkjson
	if err != nil {
		h.log.Error("JSON serialization error", "error", err)

		http.Error(w, responses.InternalServerErrorMessage, http.StatusInternalServerError)

		return
	}

	w.Write(jsonResp)

	h.log.Info(
		"Message was sent",
		"response", resp,
		"status code", statusCode,
	)
}
