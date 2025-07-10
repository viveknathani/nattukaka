package shared

import (
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ParseBody parses the request body into the provided interface.
func ParseBody(c *fiber.Ctx, v interface{}) error {
	if err := c.BodyParser(v); err != nil {
		return err
	}

	return nil
}

// SendStandardResponse is written to standardize all our API responses
func SendStandardResponse(
	c *fiber.Ctx,
	code int,
	data *map[string]interface{},
	message string,
) error {
	return c.Status(code).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}

// GetUserID use this to retreive the user id from your request context
func GetUserID(c *fiber.Ctx) int {
	value, ok := c.Locals("userID").(int)
	if ok {
		return value
	}
	return 0
}

// IsDNSPointingToVPS checks if the DNS record for the domain points to the expected IP
func IsDNSPointingToVPS(domain string, expectedIP string) (bool, error) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		return false, err
	}
	for _, ip := range ips {
		if strings.TrimSpace(ip) == expectedIP {
			return true, nil
		}
	}
	return false, nil
}
