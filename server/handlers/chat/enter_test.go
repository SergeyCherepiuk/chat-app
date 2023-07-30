package chathandler_test

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TestUnauthorizedEnterRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/1/enter", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func TestValidEnterRequest(t *testing.T) {
	randomBytes := make([]byte, 16)
	io.ReadFull(rand.Reader, randomBytes)
	clientKey := base64.StdEncoding.EncodeToString(randomBytes)

	request := httptest.NewRequest(http.MethodGet, "/1/enter", nil)
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Sec-WebSocket-Key", clientKey)
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusSwitchingProtocols {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusSwitchingProtocols,
			response.StatusCode,
		)
	}
}
