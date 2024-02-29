package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	pingpongErrors "pingpongpoc/internal/errors"
	"pingpongpoc/internal/transport/responses"
)

type (
	Service interface {
		Serve(context.Context, string) (string, error)
	}

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

func (h *Handler) Ping(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	h.log.Info("Handling new Unix socket connection")

	inputBytes, err := io.ReadAll(conn)
	if err != nil {
		h.log.Error("Reading error", "error", err)
		h.respondWithError(conn, "Error reading input")

		return
	}

	input := string(inputBytes)

	output, err := h.service.Serve(ctx, input)
	if err != nil {
		errorMessage := "Error: " + err.Error()

		if errors.Is(err, pingpongErrors.ErrNotValidInput) {
			h.log.Info("User sent invalid input", "input", input)

			errorMessage = responses.InvalidUserResponse
		} else {
			h.log.Error("Serving error", "error", err)
		}

		h.respondWithError(conn, errorMessage)

		return
	}

	fmt.Fprintf(conn, "Response: %s\n", output)

	h.log.Info("Finished successfully")
}

func (Handler) respondWithError(conn net.Conn, errMsg string) {
	resp := responses.Response{Error: errMsg} //nolint:exhaustruct

	jsonResp, err := json.Marshal(resp) //nolint:errchkjson
	if err != nil {
		conn.Write([]byte(responses.InternalServerErrorMessage + "\n"))

		return
	}

	conn.Write(jsonResp)
}
