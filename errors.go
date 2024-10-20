package gor

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Error represents an error that occurred while handling a request.
type Error struct {
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
}

// Error makes it compatible with the `error` interface.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new Error instance with an optional message
func NewError(status Status, message ...string) *Error {
	err := &Error{
		Status: status,
	}

	if len(message) > 0 {
		err.Message = message[0]
	}

	return err
}

var (
	ErrHandler = func(c *fiber.Ctx, err error) error {
		c.Set(
			HeaderContentType,
			ContentTypeTextPlainCharsetUTF8,
		)

		// Tries to parse error as gor.Error or fiber.Error
		var (
			e *Error
			f *fiber.Error
		)

		if !errors.As(err, &e) {
			if !errors.As(err, &f) {
				response := Json{
					"error": err.Error(),
				}

				return c.Status(StatusInternalServerError.Code).JSON(response)
			}

			var (
				code       = f.Code
				code_label = utils.StatusMessage(code)
			)

			// If an invalid status code is provided,
			// replaces with [500] Internal Server Error.
			if f.Code == 0 || utils.StatusMessage(code) == "" {
				code = fiber.StatusInternalServerError
				code_label = "Internal Server Error"
			}

			response := Json{
				"code":    code_label,
				"message": f.Message,
			}

			return c.Status(code).JSON(response)
		}

		return c.Status(e.Status.Code).JSON(e)
	}
)
