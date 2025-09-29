package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/0xdbb/eggsplore/util"

	db "github.com/0xdbb/eggsplore/internal/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AccountLoginRequest struct {
	Email    string `json:"email" binding:"required" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,min=8,StrongPassword" example:"password123{#Pbb"`
}

type AccountLoginResponse struct {
	// SessionID             uuid.UUID    `json:"session_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	AccessToken           string          `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	AccessTokenExpiresAt  time.Time       `json:"access_token_expires_at" example:"2025-02-05T13:15:08Z"`
	RefreshToken          string          `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	RefreshTokenExpiresAt time.Time       `json:"refresh_token_expires_at" example:"2025-02-06T13:15:08Z"`
	Account               AccountResponse `json:"user"`
}

// Account Response
type AccountResponse struct {
	ID         uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserName   string    `json:"user_name"  example:"John_doe142"`
	FirstName  string    `json:"first_name"  example:"John_doe11"`
	LastName   string    `json:"last_name"  example:"John_doe11"`
	Email      string    `json:"email" example:"john.doe@example.com"`
	Role       string    `json:"role" example:"ADMIN"`
	CreatedAt  time.Time `json:"created_at" example:"2025-01-01T12:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2025-01-02T12:00:00Z"`
	Status     string    `json:"status" example:"ACTIVE" bson:"status"`
	LastActive time.Time `json:"last_active" example:"2025-01-02T12:00:00Z"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// @BasePath /api/v1

// setCookie sets a secure HTTP-only cookie

func setCookie(ctx *gin.Context, name, value string, maxAge int) {
	isLocal := ctx.Request.Host == "localhost:8080" || ctx.Request.Host == "127.0.0.1:8080"

	domain := ""
	secure := false

	if !isLocal {
		domain = ".blvcksapphire.com"
		secure = true
	}

	ctx.SetCookie(
		name,
		value,
		maxAge, // in seconds
		"/",    // path
		domain, // domain
		secure, // secure
		true,   // httpOnly
	)
}

// @Summary		Login Account
// @Description	Login account with email and password
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		request	body		AccountLoginRequest	true	"Account Login Request"
// @Success		200		{object}	Message
// @Failure		400		{object}	ErrorResponse
// @Failure		404		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/auth/login [post]
func (s *Server) Login(ctx *gin.Context) {
	var req AccountLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Try to extract and return validation errors
		if valErr := HandleValidationError(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}

		// Fallback for non-validation binding errors
		ctx.JSON(http.StatusBadRequest, &ErrorResponse{
			Status:  "error",
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	account, err := s.db.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, HandleError(nil, http.StatusNotFound, "Invalid email or password"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Error retrieving account"))
		return
	}

	// Verify password
	if err := util.VerifyPassword(account.Password, req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, "Invalid email or password"))
		return
	}

	s.issueTokensAndRespond(ctx, account)
}

type RegisterAccountRequest struct {
	Email     string `json:"email" binding:"required" example:"john.doe@example.com"`
	FirstName string `json:"first_name"  example:"John"`
	LastName  string `json:"last_name"  example:"Doe"`
	UserName  string `json:"username" binding:"required,min=3,max=30,alphanum" example:"John_doe11"`
	Password  string `json:"password" binding:"required,min=8,StrongPassword" example:"password123{#Pbb"`
}

// @Summary		Register Account
// @Description	Register account with email and password
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		request	body		RegisterAccountRequest	true	"Account Login Request"
// @Success		200		{object}	Message
// @Failure		400		{object}	ErrorResponse
// @Failure		404		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/auth/register [post]
func (s *Server) Register(ctx *gin.Context) {
	var req RegisterAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Try to extract and return validation errors
		if valErr := HandleValidationError(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}

		// Fallback for non-validation binding errors
		ctx.JSON(http.StatusBadRequest, &ErrorResponse{
			Status:  "error",
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if req.UserName == "" {
		req.UserName = req.FirstName + "_" + req.LastName + util.RandomString(3)
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
		ctx.JSON(http.StatusInternalServerError, HandleError(nil, http.StatusInternalServerError, "Error hashing password"))
		return
	}

	arg := db.CreateAccountParams{
		Email:     req.Email,
		Password:  hashPassword,
		FirstName: stringToPgtype(req.FirstName),
		LastName:  stringToPgtype(req.LastName),
		Username:  stringToPgtype(req.UserName),
	}

	_, err = s.db.CreateAccount(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrUniqueViolation) {
			ctx.JSON(http.StatusNotFound, HandleError(nil, http.StatusNotFound, "User with this email or username already exists"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Error creating User"))
		return
	}

	ctx.JSON(http.StatusOK, Message{
		Message: "Account created successfully",
	})
}

// @Summary		Logout Account
// @Description	Logout account by deleting session and clearing cookies
// @Tags		auth
// @Produce		json
// @Success		200	{object}	Message
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/auth/logout [post]
func (h *Server) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		// If not in cookies, check the Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, "No refresh token found in cookies or header"))
			return
		}

		// Extract token from "Bearer <token>" format
		const bearerPrefix = "Bearer "
		if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
			refreshToken = authHeader[len(bearerPrefix):]
		} else {
			ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, "Invalid Authorization header format"))
			return
		}
	}

	// Validate refresh token to get session ID
	refreshPayload, err := h.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, HandleError(err, http.StatusUnauthorized, "Invalid refresh token"))
		return
	}

	// Delete session from DB
	err = h.db.DeleteSession(ctx, refreshPayload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to delete session"))
		return
	}

	// Clear cookies by setting expired cookies
	clearCookie(ctx, "access_token")
	clearCookie(ctx, "refresh_token")

	ctx.JSON(http.StatusOK, HandleMessage("Logged out successfully"))
}

// clearCookie clears a cookie by setting it with MaxAge -1
func clearCookie(ctx *gin.Context, name string) {
	ctx.SetCookie(
		name,
		"",
		-1, // MaxAge -1 deletes the cookie
		"/",
		"",
		ctx.Request.TLS != nil,
		true,
	)
}

// issueTokensAndRespond generates tokens and creates a session
func (s *Server) issueTokensAndRespond(ctx *gin.Context, account db.Accounts) {
	// Generate tokens
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(account.ID, account.Role, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to create access token"))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(account.ID, account.Role, s.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to create refresh token"))
		return
	}

	// Create session
	_, err = s.db.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		AccountID:    account.ID,
		RefreshToken: refreshToken,
		AccountAgent: ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpireAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to create session"))
		return
	}

	// Set cookies
	setCookie(ctx, "access_token", accessToken, int(s.config.AccessTokenDuration.Seconds()))
	setCookie(ctx, "refresh_token", refreshToken, int(s.config.RefreshTokenDuration.Seconds()))

	ctx.JSON(http.StatusOK, AccountLoginResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpireAt,
	})
}

func pgtypeToString(p pgtype.Text) string {
	if p.Valid {
		return p.String
	}
	return ""
}

func stringToPgtype(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	} else {
		return pgtype.Text{String: s, Valid: true}
	}
}
