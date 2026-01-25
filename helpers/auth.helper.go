package helpers

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/logger"
)

var (
	IsSystemUserKey = &contextKey{"is_system_user"}
)

// AuthMiddleware validates that requests are signed by auth-service
type AuthMiddleware struct {
	publicKey string // JWT public key for signature verification
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(publicKey string) *AuthMiddleware {
	return &AuthMiddleware{
		publicKey: publicKey,
	}
}

// ValidateServiceToken validates that the internal service token is signed by auth-service
// This middleware expects X-Service-Token header (set by APISIX after auth-service validation)
func (m *AuthMiddleware) ValidateServiceToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract service token from X-Service-Token header (set by APISIX/auth-service)
		tokenString := r.Header.Get("X-Service-Token")
		if tokenString == "" {
			logger.Warn("X-Service-Token header is required")
			HandleAppError(w, appError.New(appError.UnauthorizedAccess, "X-Service-Token header is required", 401, nil))
			return
		}

		// Parse and validate JWT token signature
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Parse public key
			block, _ := pem.Decode([]byte(m.publicKey))
			if block == nil {
				return nil, fmt.Errorf("failed to parse PEM block")
			}

			pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse public key: %w", err)
			}

			rsaPubKey, ok := pubKey.(*rsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("not an RSA public key")
			}

			return rsaPubKey, nil
		})

		if err != nil {
			logger.Warn("Failed to verify service token signature", err)
			HandleAppError(w, appError.New(appError.UnauthorizedAccess, "Invalid or expired service token", 401, err))
			return
		}

		if !token.Valid {
			logger.Warn("Invalid service token")
			HandleAppError(w, appError.New(appError.UnauthorizedAccess, "Invalid service token", 401, nil))
			return
		}

		// Verify token type is "service"
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn("Invalid token claims")
			HandleAppError(w, appError.New(appError.UnauthorizedAccess, "Invalid token claims", 401, nil))
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "service" {
			logger.Warn("Invalid token type, expected service token")
			HandleAppError(w, appError.New(appError.UnauthorizedAccess, "Invalid token type, expected service token", 401, nil))
			return
		}

		// Extract X-User-ID header (set by auth-service via forward-auth)
		userID := r.Header.Get("X-User-ID")
		ctx := r.Context()
		if userID != "" {
			ctx = context.WithValue(ctx, UserIDKey, userID)
		}

		// Extract X-Is-System-User header (set by auth-service via forward-auth)
		isSystemUser := false
		if isSystemUserHeader := r.Header.Get("X-Is-System-User"); isSystemUserHeader != "" {
			if parsed, err := strconv.ParseBool(isSystemUserHeader); err == nil {
				isSystemUser = parsed
			}
		}

		// Set system user status in context for downstream handlers
		ctx = context.WithValue(ctx, IsSystemUserKey, isSystemUser)

		// Service token signature is valid - request was validated by auth-service
		// Proceed to next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IsSystemUser retrieves the system user status from the request context
// Returns true if the user is a system user (SuperAdmin or IsSystem flag), false otherwise
func IsSystemUser(ctx context.Context) bool {
	if isSystemUser, ok := ctx.Value(IsSystemUserKey).(bool); ok {
		return isSystemUser
	}
	return false
}
