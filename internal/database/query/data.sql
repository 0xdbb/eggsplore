-- name: ListRivers :many
SELECT ST_AsGeoJSON(geom)::jsonb AS geometry
FROM rivers;

-- name: ListDistricts :many
SELECT  district, region, ST_AsGeoJSON(geom)::jsonb AS geometry
FROM districts;

-- name: ListFirst10Districts :many
SELECT district
FROM districts
LIMIT 10;

-- name: ListDistrictsByName :many
SELECT district
FROM districts
WHERE district ILIKE '%' || @text::text || '%';

-- name: ListConcessions :many
SELECT name, owner, type, status, assets, start_date, expiry_dat, ST_AsGeoJSON(geom)::jsonb AS geometry
FROM concessions;

-- name: ListForestReserves :many
SELECT name, category, ST_AsGeoJSON(geom)::jsonb AS geometry
FROM forest_reserves;

-- name: ListMiningStatic :many
SELECT id, district,severity_type, status, severity_score, severity, all_violation_types, area, proximity_to_water, inside_forest_reserve, distance_to_water_m, distance_to_forest_m,  detection_date, task_id, ST_AsGeoJSON(geom)::jsonb AS geometry
FROM mining_static;

-- name: UpdateMiningStaticStatus :one
UPDATE mining_static
SET status = @status
WHERE id = @id
RETURNING id, district, status, severity_type, severity_score, severity, all_violation_types, area, proximity_to_water, inside_forest_reserve, distance_to_water_m, distance_to_forest_m, detection_date, ST_AsGeoJSON(geom)::jsonb AS geometry;

-- name: CalculatePriorityIndex :many
WITH reports_transformed AS (
    SELECT
        ST_Transform(location, 3857) AS geom,
        'report' AS type,
        1.0 AS weight
    FROM reports
    WHERE created_at >= @start_date::timestamp
      AND created_at <= @end_date::timestamp
),
segments_raw AS (
    SELECT
        ST_Centroid(ST_Transform(geom, 3857)) AS geom,
        'segment' AS type,
        ST_Area(ST_Transform(geom, 3857)) AS area
    FROM mining_static
    WHERE detection_date >= @start_date::timestamp
      AND detection_date <= @end_date::timestamp
),
area_stats AS (
    SELECT 
        MIN(area) AS min_area,
        MAX(area) AS max_area
    FROM segments_raw
),
segments_transformed AS (
    SELECT
        geom,
        type,
        CASE
            WHEN (area_stats.max_area - area_stats.min_area) = 0 THEN 1.0
            ELSE (s.area - area_stats.min_area) / (area_stats.max_area - area_stats.min_area)
        END AS weight
    FROM segments_raw s, area_stats
)

-- Final SELECT projecting back to 4326
SELECT
    ST_AsGeoJSON(ST_Transform(geom, 4326))::jsonb AS geometry,
    type,
    weight
FROM reports_transformed

UNION ALL

SELECT
    ST_AsGeoJSON(ST_Transform(geom, 4326))::jsonb AS geometry,
    type,
    weight
FROM segments_transformed;
