package server

import (
	"errors"
	"fmt"
	"github.com/0xdbb/eggsplore/token"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("tokenMaker", tokenMaker) // Store the token maker in the context

		payload, err := ExtractTokenPayload(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, HandleError(err, http.StatusUnauthorized, "Invalid Token"))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// ExtractTokenPayload extracts the token payload from cookie or header
func ExtractTokenPayload(ctx *gin.Context) (*token.Payload, error) {
	tokenMaker := ctx.MustGet("tokenMaker").(token.Maker)

	// Extract access token from cookie first
	accessToken, err := ctx.Cookie("access_token")
	if err != nil || accessToken == "" {
		log.Println("No access token found in cookie, checking header...")
		// Fallback to Authorization header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return nil, errors.New("access token not provided in cookie or header")
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return nil, errors.New("invalid authorization header format")
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			return nil, fmt.Errorf("unsupported authorization type %s", authorizationType)
		}

		accessToken = fields[1]
	}

	// Verify token
	payload, err := tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return payload, nil
}

// UpdateLastActiveMiddleware updates the last_active and updated_at fields for the authenticated account.
func (s *Server) UpdateLastActiveMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract payload from context
		payload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, "No authorization payload in context"))
			ctx.Abort()
			return
		}

		// Assert payload type
		tokenPayload, ok := payload.(*token.Payload)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, HandleError(nil, http.StatusInternalServerError, "Invalid payload type"))
			ctx.Abort()
			return
		}

		// Validate AccountID
		if tokenPayload.AccountID == uuid.Nil {
			ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, "Invalid account ID"))
			ctx.Abort()
			return
		}

		// Update last_active
		if err := s.db.UpdateAccountLastActive(ctx, tokenPayload.AccountID); err != nil {
			ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to update last active"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
