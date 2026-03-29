package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Regex patterns for secret detection
var SecretPatterns = map[string]*regexp.Regexp{
	"OPENAI_KEY":              regexp.MustCompile(`sk-[a-zA-Z0-9_-]{20,}`),
	"EMAIL":                   regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
	"PHONE":                   regexp.MustCompile(`\b(\d{3}[-.]?)?\d{3}[-.]?\d{4}\b`),
	"SSN":                     regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
	"ROUTING_NUMBER":          regexp.MustCompile(`(?i)(?:routing|aba|transit)[\s#:_-]*(?:number|num|no)?[\s#:_-]*\b([0-9]{9})\b`),
	"US_BANK_ACCOUNT":         regexp.MustCompile(`(?i)(?:account|acct|bank)[\s#:_-]*(?:number|num|no)?[\s#:_-]*\b([0-9]{8,17})\b`),
	"CREDIT_CARD":             regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`),
	"GITHUB_TOKEN":            regexp.MustCompile(`ghp_[a-zA-Z0-9]{36,}`),
	"STRIPE_KEY":              regexp.MustCompile(`sk_(live|test)_[a-zA-Z0-9]{24,}`),
	"GOOGLE_KEY":              regexp.MustCompile(`AIza[0-9A-Za-z\-_]{35}`),
	"AWS_KEY":                 regexp.MustCompile(`(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{12,}`),
	"AWS_SECRET_KEY":          regexp.MustCompile(`(?i)(aws_secret_access_key|aws_secret_key|secret_key|SecretAccessKey)(?:\s+is)?[\s:=]+['"]?([A-Za-z0-9/+=]{40})['"]?`),
	"MONGO_URI":               regexp.MustCompile(`mongodb(?:\+srv)?:\/\/[^\s]+`),
	"UUID":                    regexp.MustCompile(`\b[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\b`),
	"TWILIO_SID":              regexp.MustCompile(`\bAC[a-f0-9]{32}\b`),
	"AZURE_CONNECTION_STRING": regexp.MustCompile(`AccountKey=([a-zA-Z0-9+/=]{20,})`),
	"JWT":                     regexp.MustCompile(`eyJ[a-zA-Z0-9_-]+\.eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`),
	"GOOGLE_PRIVATE_KEY":      regexp.MustCompile(`(?s)-----BEGIN PRIVATE KEY-----.*?-----END PRIVATE KEY-----`),
	"PRIVATE_KEY":             regexp.MustCompile(`-----BEGIN [A-Z ]+ PRIVATE KEY(?: BLOCK)?-----[\s\S]*?-----END [A-Z ]+ PRIVATE KEY(?: BLOCK)?-----`),
	"GENERIC_API":             regexp.MustCompile(`(?i)(api[\s_-]?key|secret[\s_-]?key|access[\s_-]?token)(?:\s+is)?[\s:=]+['"]?([^\s'"]{20,})['"]?`),
	"CLIENT_ID":               regexp.MustCompile(`(?i)(client[\s_-]?id)(?:\s+is)?[\s:=]+['"]?([^\s'"]{10,})['"]?`),
	"CLIENT_SECRET":           regexp.MustCompile(`(?i)(client[\s_-]?secret)(?:\s+is)?[\s:=]+['"]?([^\s'"]{10,})['"]?`),
	"DOCKER_AUTH":             regexp.MustCompile(`(?i)"auth"\s*:\s*"([^"]+)"`),
	"JWK_PRIVATE_KEY":         regexp.MustCompile(`"(d|p|q|dp|dq|qi|k|oth)"\s*:\s*"([a-zA-Z0-9_\-]{20,})"`),
}

// RedactSecrets finds secrets and replaces them with tokens
// redisClient can be nil if caching is not needed (e.g. strict simulation)
func RedactSecrets(ctx context.Context, input string, clientID string, rdb *redis.Client) (string, map[string]string) {
	secrets := make(map[string]string)
	output := input

	for label, regex := range SecretPatterns {
		output = regex.ReplaceAllStringFunc(output, func(match string) string {
			// Generate unique token
			timestamp := time.Now().UnixNano()
			token := fmt.Sprintf("<SECRET:%s:%d>", label, timestamp)

			// Store mapping
			secrets[token] = match

			// Cache in Redis (TTL 2 hours — supports long multi-step conversations)
			if rdb != nil {
				rdb.Set(ctx, token, match, 2*time.Hour)
				log.Printf("[%s] Redacted %s: %s -> %s", clientID, label, MaskSecret(match), token)
			}

			return token
		})
	}

	return output, secrets
}

// RehydrateSecrets restores original secrets from tokens
func RehydrateSecrets(ctx context.Context, input string, secretMap map[string]string, rdb *redis.Client) string {
	output := input

	// First try in-memory map
	for token, original := range secretMap {
		if input != replaceAll(input, token, original) {
			output = replaceAll(output, token, original)
			if regexp.MustCompile(regexp.QuoteMeta(token)).MatchString(input) {
				log.Printf("[Rehydration] Restored %s -> %s", token, MaskSecret(original))
			}
		}

		// Handle potential HTML escaping of the token by the AI or frameworks
		escapedToken1 := strings.ReplaceAll(strings.ReplaceAll(token, "<", "&lt;"), ">", "&gt;")
		escapedToken2 := strings.ReplaceAll(strings.ReplaceAll(token, "<", "\\u003c"), ">", "\\u003e")

		if input != replaceAll(input, escapedToken1, original) {
			output = replaceAll(output, escapedToken1, original)
			log.Printf("[Rehydration] Restored (HTML encoded) %s -> %s", token, MaskSecret(original))
		}
		if input != replaceAll(input, escapedToken2, original) {
			output = replaceAll(output, escapedToken2, original)
			log.Printf("[Rehydration] Restored (Unicode encoded) %s -> %s", token, MaskSecret(original))
		}
	}

	// Fallback to Redis for any remaining tokens
	tokenRegex := regexp.MustCompile(`(?:<|&lt;|\\u003c)SECRET:\s*([A-Z_]+)\s*:\s*(\d+)\s*(?:>|&gt;|\\u003e)`)
	output = tokenRegex.ReplaceAllStringFunc(output, func(match string) string {
		submatches := tokenRegex.FindStringSubmatch(match)
		if len(submatches) != 3 {
			return match
		}
		label := submatches[1]
		id := submatches[2]

		strictToken := fmt.Sprintf("<SECRET:%s:%s>", label, id)

		if val, ok := secretMap[strictToken]; ok {
			return val
		}

		if rdb != nil {
			if val, err := rdb.Get(ctx, strictToken).Result(); err == nil {
				log.Printf("[Rehydration] Restored (Redis) %s -> %s", strictToken, MaskSecret(val))
				return val
			}
		}

		return match
	})

	return output
}

func replaceAll(s, old, new string) string {
	return regexp.MustCompile(regexp.QuoteMeta(old)).ReplaceAllString(s, new)
}

func MaskSecret(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "..." + s[len(s)-4:]
}

// SanitizeMap creates a copy of the secret map with masked values
func SanitizeMap(secrets map[string]string) map[string]string {
	sanitized := make(map[string]string)
	for k, v := range secrets {
		sanitized[k] = MaskSecret(v)
	}
	return sanitized
}
