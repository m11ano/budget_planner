package middleware

import (
	"net"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(corsAllowOrigins []string) func(*fiber.Ctx) error {
	regexps := make([]*regexp.Regexp, 0, len(corsAllowOrigins))

	for _, raw := range corsAllowOrigins {
		raw = strings.TrimSpace(raw)

		if raw == "*" {
			return cors.New(cors.Config{
				AllowCredentials: true,
				AllowOriginsFunc: func(origin string) bool {
					return true
				},
			})
		}

		if ip := net.ParseIP(raw); ip != nil {
			pattern := `^https?://` +
				regexp.QuoteMeta(raw) +
				`$`

			regexps = append(regexps, regexp.MustCompile(pattern))
			continue
		}

		pattern := `^https?://([a-zA-Z0-9-]+\.)?` +
			regexp.QuoteMeta(raw) +
			`$`

		regexps = append(regexps, regexp.MustCompile(pattern))
	}

	return cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			for _, r := range regexps {
				if r.MatchString(origin) {
					return true
				}
			}
			return false
		},
	})
}
