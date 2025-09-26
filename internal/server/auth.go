package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/0xdbb/eggsplore/util"

	db "github.com/0xdbb/eggsplore/internal/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// TODO: Forgot password
var (
	SuperADMIN            = "SUPER_ADMIN"
	Admin                 = "ADMIN"
	Standard              = "STANDARD"
	AccountActive         = "ACTIVE"
	AccountInActive       = "INACTIVE"
	AccountPending        = "PENDING"
	SignupTokenDuration   = 48 * time.Hour
	OTPExpirationDuration = 1 * time.Minute
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
	Department string    `json:"department" example:"Minerals Department"`
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

	// Check if account is active
	if account.Status != AccountActive {
		ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, "Account not activated. Please complete setup."))
		return
	}

	// Verify password
	if err := util.VerifyPassword(account.Password.String, req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, "Invalid email or password"))
		return
	}
	//
	// Generate and store new OTP
	otp := util.GenerateOTP()
	expiresAt := time.Now().Add(OTPExpirationDuration)

	args := db.UpdateAccountOTPParams{
		ID:           account.ID,
		OtpCode:      stringToPgtype(otp),
		OtpExpiresAt: stringToPgtype(expiresAt.Format(time.RFC3339)),
	}

	if err := s.db.UpdateAccountOTP(ctx, args); err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to store OTP"))
		return
	}

	if err := util.SendVerificationEmail(account.Email, otp, s.config.ResendApiKey); err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to send OTP"))
		return
	}

	ctx.JSON(http.StatusOK, HandleMessage("An OTP has been sent to your registered email address. Please check your email and enter the OTP to complete the login process."))
}

// @Summary      Verify 2FA OTP to complete login
// @Description  Verifies the OTP sent to the account's email and generates access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body VerifyOTPRequest true "OTP Verification Request"
// @Success      200  {object}  AccountLoginResponse "Successful verification and login"
// @Failure      400  {object}  ErrorResponse "Invalid request"
// @Failure      401  {object}  ErrorResponse "Invalid or expired OTP"
// @Failure      404  {object}  ErrorResponse "Invalid email or OTP"
// @Failure      500  {object}  ErrorResponse "Internal server error"
// @Router       /auth/verify-otp [post]
func (s *Server) VerifyOTP(ctx *gin.Context) {
	var req VerifyOTPRequest

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
			ctx.JSON(http.StatusNotFound, HandleError(nil, http.StatusNotFound, "Invalid email or OTP"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Error retrieving account"))
		return
	}

	if err := s.verifyOTP(ctx, account, req.OTP); err != nil {
		ctx.JSON(http.StatusUnauthorized, HandleError(nil, http.StatusUnauthorized, err.Error()))
		return
	}

	// Generate tokens
	s.issueTokensAndRespond(ctx, account)
}

// @Summary      Resend OTP
// @Description  Resend OTP to a verified account's email if previous OTP is expired
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body SendOTPRequest true "Send OTP Request"
// @Success      200 {object} Message "OTP resent successfully"
// @Failure      400 {object} ErrorResponse "Invalid request"
// @Failure      401 {object} ErrorResponse "Email not verified or OTP not expired"
// @Failure      404 {object} ErrorResponse "Account not found"
// @Failure      429 {object} ErrorResponse "OTP still valid, wait before requesting new one"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /auth/send-otp [post]
func (s *Server) SendOTP(ctx *gin.Context) {
	var req SendOTPRequest

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
			ctx.JSON(http.StatusNotFound, HandleError(nil, http.StatusNotFound, "Account not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to retrieve account"))
		return
	}

	// Check if existing OTP is still valid
	if !account.OtpCode.Valid || strings.TrimSpace(account.OtpCode.String) == "" {
		otpExpiresAt, err := time.Parse(time.RFC3339, account.OtpExpiresAt.String)
		if err == nil && time.Now().Before(otpExpiresAt) {
			ctx.JSON(http.StatusTooManyRequests, HandleError(nil, http.StatusTooManyRequests, "OTP already sent. Please wait before requesting a new one."))
			return
		}
	}

	// Generate and store new OTP
	otp := util.GenerateOTP()
	expiresAt := time.Now().Add(OTPExpirationDuration)

	args := db.UpdateAccountOTPParams{
		ID: account.ID,
		OtpCode: pgtype.Text{
			String: otp,
			Valid:  true,
		},
		OtpExpiresAt: pgtype.Text{
			String: expiresAt.Format(time.RFC3339),
			Valid:  true,
		},
	}

	if err := s.db.UpdateAccountOTP(ctx, args); err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to store OTP"))
		return
	}

	if err := util.SendVerificationEmail(account.Email, otp, s.config.ResendApiKey); err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to send OTP"))
		return
	}

	ctx.JSON(http.StatusOK, HandleMessage("OTP resent successfully. Please check your email."))
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

// verifyOTP is a reusable method to verify OTP
func (s *Server) verifyOTP(ctx *gin.Context, account db.Accounts, otp string) error {
	if !account.OtpCode.Valid || account.OtpCode.String == "" {
		return errors.New("OTP not found")
	}

	otpExpiry, err := time.Parse(time.RFC3339, account.OtpExpiresAt.String)
	if err != nil || time.Now().After(otpExpiry) {
		return errors.New("OTP has expired, please request a new one")
	}

	if account.OtpCode.String != otp {
		return errors.New("invalid OTP")
	}

	// Clear OTP after successful use
	if err := s.db.ClearAccountOTP(ctx, account.ID); err != nil {
		return errors.New("failed to clear OTP")
	}
	return nil
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
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
