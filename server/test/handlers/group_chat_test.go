package handlers_test

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Test_Unauthorized_EnterChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/api/chat/markwatson", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_UserNotFound_EnterChatRequest(t *testing.T) {
	randomBytes := make([]byte, 16)
	io.ReadFull(rand.Reader, randomBytes)
	clientKey := base64.StdEncoding.EncodeToString(randomBytes)

	request := httptest.NewRequest(http.MethodGet, "/api/chat/unknownusername", nil)
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
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func TestValidEnterRequest(t *testing.T) {
	randomBytes := make([]byte, 16)
	io.ReadFull(rand.Reader, randomBytes)
	clientKey := base64.StdEncoding.EncodeToString(randomBytes)

	request := httptest.NewRequest(http.MethodGet, "/api/chat/markwatson", nil)
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

func Test_Unauthorized_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/api/chat/markwatson/1", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_UnparsableBody_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodPut, "/api/chat/markwatson/1", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnprocessableEntity,
			response.StatusCode,
		)
	}
}

func Test_InvalidBody_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/api/chat/markwatson/1",
		strings.NewReader(`{
			"message": ""
		}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func Test_UserNotAnAuthor_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/api/chat/markwatson/2",
		strings.NewReader(`{
			"message": "New message"
		}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_Valid_UpdateMessageRequest(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodPut,
		"/api/chat/markwatson/1",
		strings.NewReader(`{
			"message": "New message"
		}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusOK {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusOK,
			response.StatusCode,
		)
	}
}

func Test_Unauthorized_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/chat/johndoe/1", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_UserNotFound_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/chat/unknownusername/1", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func Test_MessageNotBelongsToChat_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/chat/markwatson/3", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func Test_Valid_DeleteMessageRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/chat/markwatson/2", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusOK {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusOK,
			response.StatusCode,
		)
	}
}

func Test_Unauthorized_DeleteChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/chat/johndoe", nil)

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusUnauthorized {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusUnauthorized,
			response.StatusCode,
		)
	}
}

func Test_Invalid_DeleteChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/chat/unknownusername", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusBadRequest {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusBadRequest,
			response.StatusCode,
		)
	}
}

func Test_Valid_DeleteChatRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodDelete, "/api/chat/markwatson", nil)
	request.AddCookie(&http.Cookie{
		Name:     "session_id",
		Value:    uuid.NewString(),
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	response, _ := app.Test(request)
	if response.StatusCode != fiber.StatusOK {
		t.Errorf(
			"expected status code: %v, actual status code: %v\n",
			fiber.StatusOK,
			response.StatusCode,
		)
	}
}
