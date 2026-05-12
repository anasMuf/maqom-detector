package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const SessionContextKey = "session_id"

// SessionMiddleware validates X-Session-ID header and injects into context
func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionHeader := c.Request().Header.Get("X-Session-ID")
		if sessionHeader == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "VALIDATION_ERROR",
					"message": "Header X-Session-ID wajib diisi",
				},
			})
		}

		sessionID, err := uuid.Parse(sessionHeader)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "VALIDATION_ERROR",
					"message": "X-Session-ID harus berupa UUID yang valid",
				},
			})
		}

		c.Set(SessionContextKey, sessionID)
		return next(c)
	}
}

// GetSessionID retrieves the session ID from echo context
func GetSessionID(c echo.Context) uuid.UUID {
	return c.Get(SessionContextKey).(uuid.UUID)
}
