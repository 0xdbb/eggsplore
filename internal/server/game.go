package server

import (
	"net/http"
	"time"

	db "github.com/0xdbb/eggsplore/internal/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DropEggRequest represents a player dropping an egg
type DropEggRequest struct {
	PlayerID string  `json:"player_id" binding:"required"`
	Type     string  `json:"type" binding:"required,oneof=BUNNY GOLDEN LEGENDARY"`
	Message  string  `json:"message"`
	Lat      float64 `json:"lat" binding:"required"`
	Lon      float64 `json:"lon" binding:"required"`
}

// DropEggResponse represents the egg created
type DropEggResponse struct {
	InventoryID string    `json:"inventory_id"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	CreatedAt   time.Time `json:"created_at"`
}

// @Summary		Drop Egg
// @Description	Player drops an egg (adds to inventory + eggs table)
// @Tags		game
// @Accept		json
// @Produce		json
// @Param		request	body		DropEggRequest	true	"Drop Egg Request"
// @Success		200		{object}	DropEggResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/game/eggs [post]
// DropEggRequest represents a player dropping an egg
func (s *Server) DropEgg(ctx *gin.Context) {
	var req DropEggRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if valErr := HandleValidationError(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
		ctx.JSON(http.StatusBadRequest, HandleError(nil, http.StatusBadRequest, "Invalid request format"))
		return
	}

	playerUUID, ok := parseUUID(ctx, req.PlayerID, "player_id")
	if !ok {
		return
	}

	// Step 1: create inventory row
	inv, err := s.db.CreateEgg(ctx, db.CreateEggParams{
		PlayerID:    playerUUID,
		Description: stringToPgtype(req.Message),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to create inventory egg"))
		return
	}

	// Step 2: insert egg details including location
	egg, err := s.db.AddEggDetails(ctx, db.AddEggDetailsParams{
		InventoryID: inv.ID,
		Type:        req.Type,
		Message:     stringToPgtype(req.Message),
		Lat:         req.Lat,
		Lon:         req.Lon,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to add egg details"))
		return
	}

	ctx.JSON(http.StatusOK, DropEggResponse{
		InventoryID: inv.ID.String(),
		Type:        egg.Type,
		Message:     egg.Message.String,
		Lat:         req.Lat,
		Lon:         req.Lon,
		CreatedAt:   inv.CreatedAt,
	})
}

// @Summary		Get Player Eggs
// @Description	Get all eggs belonging to a player
// @Tags		game
// @Produce		json
// @Param		player_id	query		string	true	"Player ID"
// @Success		200		{array}		GetEggsByPlayerResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/game/eggs [get]
func (s *Server) GetPlayerEggs(ctx *gin.Context) {
	playerID, ok := parseUUID(ctx, ctx.Query("player_id"), "player_id")
	if !ok {
		return
	}

	eggs, err := s.db.GetEggsByPlayer(ctx, playerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to fetch eggs"))
		return
	}

	ctx.JSON(http.StatusOK, eggs)
}

// @Summary		Get Player Tools
// @Description	Get all tools belonging to a player
// @Tags		game
// @Produce		json
// @Param		player_id	query		string	true	"Player ID"
// @Success		200		{array}		GetToolsByPlayerResponse
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/game/tools [get]
func (s *Server) GetPlayerTools(ctx *gin.Context) {
	playerID, ok := parseUUID(ctx, ctx.Query("player_id"), "player_id")
	if !ok {
		return
	}

	tools, err := s.db.GetToolsByPlayer(ctx, playerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to fetch tools"))
		return
	}

	ctx.JSON(http.StatusOK, tools)
}

// @Summary		Get Player Inventory
// @Description	Get all inventory items (tools, eggs, boosts) belonging to a player
// @Tags		game
// @Produce		json
// @Param		player_id	query		string	true	"Player ID"
// @Success		200		{array}		Inventory
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/game/inventory [get]
func (s *Server) GetPlayerInventory(ctx *gin.Context) {
	playerID, ok := parseUUID(ctx, ctx.Query("player_id"), "player_id")
	if !ok {
		return
	}

	inv, err := s.db.GetInventoryByPlayer(ctx, playerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to fetch inventory"))
		return
	}

	ctx.JSON(http.StatusOK, inv)
}

// @Summary		Get Player Stats
// @Description	Get player stats by account id
// @Tags		game
// @Produce		json
// @Param		account_id	query		string	true	"Account ID"
// @Success		200		{object}	Players
// @Failure		400		{object}	ErrorResponse
// @Failure		500		{object}	ErrorResponse
// @Router		/game/player [get]
func (s *Server) GetPlayerStats(ctx *gin.Context) {
	accountID, ok := parseUUID(ctx, ctx.Query("account_id"), "account_id")
	if !ok {
		return
	}

	player, err := s.db.GetPlayerByAccount(ctx, accountID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to fetch player stats"))
		return
	}

	ctx.JSON(http.StatusOK, player)
}

func parseUUID(ctx *gin.Context, raw string, fieldName string) (uuid.UUID, bool) {
	parsed, err := uuid.Parse(raw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, HandleError(nil, http.StatusBadRequest, "Invalid "+fieldName+" format"))
		return uuid.Nil, false
	}
	return parsed, true
}
