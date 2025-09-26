package server

import (
	"encoding/json"
	"fmt"
	db "github.com/0xdbb/eggsplore/internal/database/sqlc"
	"github.com/0xdbb/eggsplore/token"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

// ConcessionProperties represents the properties for a concession GeoJSON feature
type ConcessionProperties struct {
	Name       string    `json:"name,omitempty" example:"Block A"`
	Owner      string    `json:"owner,omitempty" example:"John Doe"`
	Type       string    `json:"type,omitempty" example:"mining"`
	Status     string    `json:"status,omitempty" example:"OPEN"`
	Assets     string    `json:"assets,omitempty" example:"3D model"`
	StartDate  time.Time `json:"start_date,omitempty" example:"2023-01-01"`
	ExpiryDate time.Time `json:"expiry_date,omitempty" example:"2025-01-01"`
}

// MiningSiteProperties represents the properties for a mining site GeoJSON feature
type MiningSiteProperties struct {
	ID                  string    `json:"id,omitempty" example:"12345"`
	District            string    `json:"district,omitempty" example:"Central District"`
	Area                float64   `json:"area,omitempty" example:"1234.56"`
	Severity            string    `json:"severity,omitempty" example:"HIGH"`
	SeverityType        string    `json:"severity_type,omitempty" example:"RIVER_VIOLATION"`
	Status              string    `json:"status,omitempty" example:"OPEN"`
	SeverityScore       int64     `json:"severity_score,omitempty" example:"5"`
	ProximityToWater    bool      `json:"proximity_to_water" example:"true"`
	InsideForestReserve bool      `json:"inside_forest_reserve,omitempty" example:"false"`
	DetectedDate        string    `json:"detected_date,omitempty" example:"2023-01-01T00:00:00Z"`
	DetectionDate       time.Time `json:"detection_date,omitempty" example:"2023-01-01T00:00:00Z"`
	AllViolationTypes   string    `json:"all_violation_types,omitempty" example:"RIVER_VIOLATION,FOREST_VIOLATION"`
	TaskID              string    `json:"task_id,omitempty" example:"fc82d367-730a-4268-acc6-c35dd60d85e3"`
	DistanceToWaterM    float64   `json:"distance_to_water_m" example:"150.0"`
	DistanceToForestM   float64   `json:"distance_to_forest_m" example:"300.0"`
}

// ForestReserveProperties represents the properties for a forest reserve GeoJSON feature
type ForestReserveProperties struct {
	Name     string `json:"name,omitempty" example:"Reserve A"`
	Category string `json:"category,omitempty" example:"reserve"`
}

// DistrictProperties represents the properties for a district GeoJSON feature
type DistrictProperties struct {
	District string `json:"district,omitempty" example:"Central District"`
	Region   string `json:"region,omitempty" example:"Ashanti"`
}

// RiverProperties represents the properties for a river GeoJSON feature
type RiverProperties struct{}

// PriorityIndexProperties represents the properties for a priority index heatmap GeoJSON feature
type PriorityIndexProperties struct {
	Type   string  `json:"type,omitempty" example:"heatmap_point"`
	Weight float64 `json:"weight,omitempty" example:"0.75"`
}

// GeoJSONFeature is a generic GeoJSON feature structure
type GeoJSONFeature[T any] struct {
	Type       string          `json:"type" example:"Feature"`
	Properties T               `json:"properties"`
	Geometry   json.RawMessage `json:"geometry" swaggertype:"object"`
}

// ConcessionsResponse represents the response for GetConcessions
type ConcessionsResponse struct {
	Type     string                                 `json:"type" example:"FeatureCollection"`
	Features []GeoJSONFeature[ConcessionProperties] `json:"features"`
}

// MiningSitesResponse represents the response for GetMiningSites
type MiningSitesResponse struct {
	Type     string                                 `json:"type" example:"FeatureCollection"`
	Features []GeoJSONFeature[MiningSiteProperties] `json:"features"`
}

// ForestReservesResponse represents the response for GetForestReserves
type ForestReservesResponse struct {
	Type     string                                    `json:"type" example:"FeatureCollection"`
	Features []GeoJSONFeature[ForestReserveProperties] `json:"features"`
}

// DistrictsResponse represents the response for GetDistricts
type DistrictsResponse struct {
	Type     string                               `json:"type" example:"FeatureCollection"`
	Features []GeoJSONFeature[DistrictProperties] `json:"features"`
}

// RiversResponse represents the response for GetRivers
type RiversResponse struct {
	Type     string                            `json:"type" example:"FeatureCollection"`
	Features []GeoJSONFeature[RiverProperties] `json:"features"`
}

// PriorityIndexHeatmapResponse represents the response for GetPriorityIndexHeatmap
type PriorityIndexHeatmapResponse struct {
	Type     string                                    `json:"type" example:"FeatureCollection"`
	Features []GeoJSONFeature[PriorityIndexProperties] `json:"features"`
}

// DistrictSearchResponse represents the response for SearchDistricts
type DistrictSearchResponse struct {
	Result []string `json:"result" example:"Central District,Northern District"`
}

// UpdateMiningStaticStatusRequest defines the request body for updating mining site status
type UpdateMiningStaticStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=UNDER_REVIEW VERIFIED FALSE_POSITIVE" example:"VERIFIED"`
}

// MiningStaticResponse defines the response for updating mining site status
type MiningStaticResponse struct {
	ID     string `json:"id" example:"12345"`
	Status string `json:"status" example:"IN_REVIEW"`
}

// Cache constants
const (
	cacheTTL          = 1 * time.Hour
	concessionsKey    = "cache:concessions"
	miningKey         = "cache:mining_sites"
	forestReservesKey = "cache:forest_reserves"
	districtsKey      = "cache:districts"
	riversKey         = "cache:rivers"
)

// GetConcessions retrieves a GeoJSON FeatureCollection of all concessions
// @Summary      Get Concessions
// @Description  Retrieve a GeoJSON FeatureCollection of all mining concessions with details like name, owner, type, status, assets, start date, and expiry date
// @Tags         data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} ConcessionsResponse
// @Failure      401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure      500 {object} ErrorResponse "Internal Server Error - Failed to retrieve concessions"
// @Router       /data/concessions [get]
func (s *Server) GetConcessions(ctx *gin.Context) {
	// Try to get from cache
	cacheKey := getUserCacheKey(concessionsKey, ctx)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response ConcessionsResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			ctx.JSON(http.StatusOK, response)
			return
		}
	}
	if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	// Fetch from database
	concessions, err := s.db.ListConcessions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to retrieve concessions"))
		return
	}

	features := make([]GeoJSONFeature[ConcessionProperties], len(concessions))
	for i, c := range concessions {
		props := ConcessionProperties{
			Name:       c.Name.String,
			Owner:      c.Owner.String,
			Type:       c.Type.String,
			Status:     c.Status.String,
			Assets:     c.Assets.String,
			StartDate:  c.StartDate.Time,
			ExpiryDate: c.ExpiryDat.Time,
		}
		features[i] = GeoJSONFeature[ConcessionProperties]{
			Type:       "Feature",
			Properties: props,
			Geometry:   c.Geometry,
		}
	}

	response := ConcessionsResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Cache the response
	data, err := json.Marshal(response)
	if err == nil {
		err = s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err()
		if err != nil {
			log.Printf("Failed to cache concessions: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary		Get Rivers
// @Description	Retrieve a GeoJSON FeatureCollection of all rivers
// @Tags		data
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Success		200	{object}	RiversResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router		/data/rivers [get]
func (s *Server) GetRivers(ctx *gin.Context) {
	cacheKey := getUserCacheKey(riversKey, ctx)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response RiversResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			ctx.JSON(http.StatusOK, response)
			return
		}
	}
	if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	// Fetch from database
	rivers, err := s.db.ListRivers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to retrieve rivers"))
		return
	}

	features := make([]GeoJSONFeature[RiverProperties], len(rivers))

	for i, geometry := range rivers {
		features[i] = GeoJSONFeature[RiverProperties]{
			Type:       "Feature",
			Properties: RiverProperties{},
			Geometry:   geometry,
		}
	}

	response := RiversResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Cache the response
	data, err := json.Marshal(response)
	if err == nil {
		err = s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err()
		if err != nil {
			log.Printf("Failed to cache rivers: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMiningSites retrieves a GeoJSON FeatureCollection of all mining sites
// @Summary      Get Mining Sites
// @Description  Retrieve a GeoJSON FeatureCollection of all mining sites with details like ID, district, area, severity, and detection date
// @Tags         data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} MiningSitesResponse
// @Failure      401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure      500 {object} ErrorResponse "Internal Server Error - Failed to retrieve mining sites"
// @Router       /data/mining-sites [get]
func (s *Server) GetMiningSites(ctx *gin.Context) {
	// Try to get from cache
	cacheKey := getUserCacheKey(miningKey, ctx)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response MiningSitesResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			ctx.JSON(http.StatusOK, response)
			return
		}
	}
	if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	// Fetch from database
	miningSites, err := s.db.ListMiningStatic(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to retrieve mining sites"))
		return
	}

	features := make([]GeoJSONFeature[MiningSiteProperties], len(miningSites))
	for i, d := range miningSites {
		props := MiningSiteProperties{
			DetectedDate:        d.DetectionDate.Time.Format(time.RFC3339),
			District:            d.District.String,
			ID:                  d.ID,
			Area:                d.Area.Float64,
			Severity:            d.Severity.String,
			SeverityType:        d.SeverityType.String,
			Status:              d.Status.String,
			SeverityScore:       d.SeverityScore.Int64,
			ProximityToWater:    d.ProximityToWater.Bool,
			InsideForestReserve: d.InsideForestReserve.Bool,
			TaskID:              d.TaskID.String(),
			DetectionDate:       d.DetectionDate.Time,
			AllViolationTypes:   d.AllViolationTypes.String,
			DistanceToWaterM:    d.DistanceToWaterM.Float64,
			DistanceToForestM:   d.DistanceToForestM.Float64,
		}
		features[i] = GeoJSONFeature[MiningSiteProperties]{
			Type:       "Feature",
			Properties: props,
			Geometry:   d.Geometry,
		}
	}

	response := MiningSitesResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Cache the response
	data, err := json.Marshal(response)
	if err == nil {
		err = s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err()
		if err != nil {
			log.Printf("Failed to cache mining sites: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// GetForestReserves retrieves a GeoJSON FeatureCollection of all forest reserves
// @Summary      Get Forest Reserves
// @Description  Retrieve a GeoJSON FeatureCollection of all forest reserves with name and category
// @Tags         data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} ForestReservesResponse
// @Failure      401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure      500 {object} ErrorResponse "Internal Server Error - Failed to retrieve forest reserves"
// @Router       /data/forest-reserves [get]
func (s *Server) GetForestReserves(ctx *gin.Context) {
	// Try to get from cache
	cacheKey := getUserCacheKey(forestReservesKey, ctx)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response ForestReservesResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			ctx.JSON(http.StatusOK, response)
			return
		}
	}
	if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	// Fetch from database
	reserves, err := s.db.ListForestReserves(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to retrieve forest reserves"))
		return
	}

	features := make([]GeoJSONFeature[ForestReserveProperties], len(reserves))
	for i, r := range reserves {
		props := ForestReserveProperties{
			Name:     r.Name.String,
			Category: r.Category.String,
		}
		features[i] = GeoJSONFeature[ForestReserveProperties]{
			Type:       "Feature",
			Properties: props,
			Geometry:   r.Geometry,
		}
	}

	response := ForestReservesResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Cache the response
	data, err := json.Marshal(response)
	if err == nil {
		err = s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err()
		if err != nil {
			log.Printf("Failed to cache forest reserves: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// GetDistricts retrieves a GeoJSON FeatureCollection of all districts
// @Summary      Get Districts
// @Description  Retrieve a GeoJSON FeatureCollection of all districts with district and region information
// @Tags         data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} DistrictsResponse
// @Failure      401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure      500 {object} ErrorResponse "Internal Server Error - Failed to retrieve districts"
// @Router       /data/districts [get]
func (s *Server) GetDistricts(ctx *gin.Context) {
	// Try to get from cache
	cacheKey := getUserCacheKey(districtsKey, ctx)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response DistrictsResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			ctx.JSON(http.StatusOK, response)
			return
		}
	}
	if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	// Fetch from database
	districts, err := s.db.ListDistricts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to retrieve districts"))
		return
	}

	features := make([]GeoJSONFeature[DistrictProperties], len(districts))
	for i, d := range districts {
		props := DistrictProperties{
			District: d.District.String,
			Region:   d.Region.String,
		}
		features[i] = GeoJSONFeature[DistrictProperties]{
			Type:       "Feature",
			Properties: props,
			Geometry:   d.Geometry,
		}
	}

	response := DistrictsResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Cache the response
	data, err := json.Marshal(response)
	if err == nil {
		err = s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err()
		if err != nil {
			log.Printf("Failed to cache districts: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// SearchDistricts retrieves a list of district names matching the provided name
// @Summary      Search Districts by Name
// @Description  Retrieve a list of district names matching the provided name (case-insensitive partial match) or first 10 districts if no name provided
// @Tags         data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        name query string false "Name to search for" example:"central"
// @Success      200 {object} DistrictSearchResponse
// @Failure      400 {object} ErrorResponse "Bad Request - Invalid query parameter"
// @Failure      401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure      500 {object} ErrorResponse "Internal Server Error - Failed to search districts"
// @Router       /data/districts/search [get]
func (s *Server) SearchDistricts(ctx *gin.Context) {
	name := ctx.Query("name")
	var (
		districts []pgtype.Text
		err       error
		cacheKey  string
	)
	if name == "" {
		cacheKey = "cache:districts:first10"
	} else {
		cacheKey = fmt.Sprintf("cache:districts:search:%s", strings.ToLower(name))
	}

	cacheKey = getUserCacheKey(cacheKey, ctx)

	// Try to get from cache
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response DistrictSearchResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			ctx.JSON(http.StatusOK, response)
			return
		}
	}
	if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	// Fetch from database
	if name == "" {
		districts, err = s.db.ListFirst10Districts(ctx)
	} else {
		districts, err = s.db.ListDistrictsByName(ctx, name)
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to search districts"))
		return
	}

	// Convert to []string
	districtStrings := make([]string, len(districts))
	for i, d := range districts {
		districtStrings[i] = d.String
	}

	response := DistrictSearchResponse{
		Result: districtStrings,
	}

	// Cache the response
	data, err := json.Marshal(response)
	if err == nil {
		if err := s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err(); err != nil {
			log.Printf("Failed to cache districts search: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateMiningStaticStatus updates the status of a mining site
// @Summary      Update Mining Site Status
// @Description  Update the status of a mining site to  UNDER_REVIEW, FALSE_POSITIVE, or VERIFIED
// @Tags         data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      string                           true  "Mining Site ID" example:"12345"
// @Param        request body      UpdateMiningStaticStatusRequest true  "Update Mining Status Request"
// @Success      200     {object}  MiningStaticResponse
// @Failure      400     {object}  ErrorResponse "Bad Request - Invalid ID or request body"
// @Failure      401     {object}  ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure      500     {object}  ErrorResponse "Internal Server Error - Failed to update mining site status"
// @Router       /data/mining-sites/{id}/status [patch]
func (s *Server) UpdateMiningStaticStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "ID parameter is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req UpdateMiningStaticStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if valErr := HandleValidationError(err); valErr != nil {
			ctx.JSON(http.StatusBadRequest, valErr)
			return
		}
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Update the DB
	arg := db.UpdateMiningStaticStatusParams{
		Status: stringToPgtype(req.Status),
		ID:     id,
	}
	_, err := s.db.UpdateMiningStaticStatus(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to update mining site status"))
		return
	}

	cacheKey := getUserCacheKey(miningKey, ctx)
	// Try to update the cached data
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response MiningSitesResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			updated := false
			for i, feature := range response.Features {
				if feature.Properties.ID == id {
					response.Features[i].Properties.Status = req.Status
					updated = true
					break
				}
			}
			if updated {
				data, err := json.Marshal(response)
				if err == nil {
					err := s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err()
					if err != nil {
						log.Printf("Failed to update cache after status update: %v", err)
					} else {
						log.Printf("Cache updated for mining site ID: %s", id)
					}
				} else {
					log.Printf("Failed to marshal updated cache: %v", err)
				}
			} else {
				log.Printf("ID %s not found in cache — skipping cache update", id)
			}
		} else {
			log.Printf("Failed to unmarshal cached data — skipping cache update")
		}
	} else if err != redis.Nil {
		log.Printf("Redis error while fetching cache: %v", err)
	}

	ctx.JSON(http.StatusOK, MiningStaticResponse{
		ID:     id,
		Status: req.Status,
	})
}

type PriorityIndexRequest struct {
	StartDate string `form:"start_date" binding:"required" example:"2024-01-01"`
	EndDate   string `form:"end_date" binding:"required" example:"2024-12-31"`
}

// GetPriorityIndexHeatmap retrieves a GeoJSON FeatureCollection of priority index heatmap points
// @Summary      Get Priority Index Heatmap
// @Description  Retrieve a GeoJSON FeatureCollection of priority index heatmap points for a given date range with type and weight
// @Tags         data
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        start_date query string true "Start date (YYYY-MM-DD)" example:"2024-01-01"
// @Param        end_date   query string true "End date (YYYY-MM-DD)" example:"2024-12-31"
// @Success      200 {object} PriorityIndexHeatmapResponse
// @Failure      400 {object} ErrorResponse "Bad Request - Invalid or missing date parameters"
// @Failure      401 {object} ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure      500 {object} ErrorResponse "Internal Server Error - Failed to retrieve priority index heatmap"
// @Router       /data/heatmap-data [get]
func (s *Server) GetPriorityIndexHeatmap(ctx *gin.Context) {
	var req PriorityIndexRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "Missing or invalid query parameters: start_date and end_date are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "Invalid start_date format. Expected YYYY-MM-DD",
			Code:    http.StatusBadRequest,
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "Invalid end_date format. Expected YYYY-MM-DD",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if endDate.Before(startDate) {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  "error",
			Message: "end_date cannot be before start_date",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Build a unique cache key based on the date params
	cacheKey := fmt.Sprintf("cache:priority_index:%s:%s", req.StartDate, req.EndDate)

	cacheKey = getUserCacheKey(cacheKey, ctx)

	// Check cache
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var response PriorityIndexHeatmapResponse
		if json.Unmarshal([]byte(cached), &response) == nil {
			ctx.JSON(http.StatusOK, response)
			return
		}
	}
	if err != nil && err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	// Query DB
	rows, err := s.db.CalculatePriorityIndex(ctx, db.CalculatePriorityIndexParams{
		StartDate: pgtype.Timestamp{
			Time:  startDate,
			Valid: true,
		},
		EndDate: pgtype.Timestamp{
			Time:  endDate,
			Valid: true,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, HandleError(err, http.StatusInternalServerError, "Failed to retrieve priority index heatmap"))
		return
	}

	features := make([]GeoJSONFeature[PriorityIndexProperties], len(rows))
	for i, r := range rows {
		props := PriorityIndexProperties{
			Type:   r.Type,
			Weight: r.Weight,
		}
		features[i] = GeoJSONFeature[PriorityIndexProperties]{
			Type:       "Feature",
			Properties: props,
			Geometry:   r.Geometry,
		}
	}

	response := PriorityIndexHeatmapResponse{
		Type:     "FeatureCollection",
		Features: features,
	}

	// Cache response
	data, err := json.Marshal(response)
	if err == nil {
		err := s.redisClient.Set(ctx, cacheKey, data, cacheTTL).Err()
		if err != nil {
			log.Printf("Failed to cache priority index heatmap: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// getMetricsCacheKey generates a user-specific cache key
func getUserCacheKey(cacheKey string, ctx *gin.Context) string {
	accountID := ctx.MustGet(authorizationPayloadKey).(*token.Payload).AccountID
	return fmt.Sprintf("cache:%s:%s", cacheKey, accountID)
}
