package server

import (
	"errors"
	"fmt"
	db "github.com/0xdbb/eggsplore/internal/database/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// @Summary      Renew Access Token
// @Description  Generates a new access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body renewAccessTokenRequest false "Refresh Token Request"
// @Success      200  {object}  renewAccessTokenResponse
// @Failure      400  {object}  ErrorResponse "Invalid request"
// @Failure      401  {object}  ErrorResponse "Unauthorized or Invalid token"
// @Failure      404  {object}  ErrorResponse "Session not found"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /auth/renew [post]
func (s *Server) RenewAccessToken(ctx *gin.Context) {
	// Try to get the refresh token from cookies first
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		// If not in cookies, try the request body
		var req renewAccessTokenRequest
		if err := ctx.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
			ctx.JSON(http.StatusUnauthorized, HandleError(err, http.StatusUnauthorized, "No refresh token found"))
			return
		}
		refreshToken = req.RefreshToken
	}

	// Verify the refresh token
	refreshPayload, err := s.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, HandleError(err, http.StatusUnauthorized, "Invalid token"))
		return
	}

	// Retrieve the session from the database
	session, err := s.db.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, HandleError(err, http.StatusNotFound, "Session not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Error retrieving session"))
		return
	}

	// Session validation checks
	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, HandleError(fmt.Errorf("blocked session"), http.StatusUnauthorized, "Blocked session"))
		return
	}

	if session.AccountID != refreshPayload.AccountID {
		ctx.JSON(http.StatusUnauthorized, HandleError(fmt.Errorf("incorrect session user"), http.StatusUnauthorized, "Invalid session user"))
		return
	}

	if session.RefreshToken != refreshToken {
		ctx.JSON(http.StatusUnauthorized, HandleError(fmt.Errorf("mismatched session token"), http.StatusUnauthorized, "Invalid session token"))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, HandleError(fmt.Errorf("expired session"), http.StatusUnauthorized, "Session expired"))
		return
	}

	// Generate a new access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.AccountID,
		refreshPayload.Role,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Error creating access token"))
		return
	}

	// Set the access token in a secure, HTTP-only cookie
	ctx.SetCookie(
		"access_token",
		accessToken,
		int(s.config.AccessTokenDuration.Seconds()),
		"/",
		"",
		true, // secure (set to false for local dev over HTTP)
		true, // httpOnly
	)

	// Return the access token in JSON response as well
	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpireAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
