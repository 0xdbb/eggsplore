package util

import (
	"fmt"
	"math"
)

type Coord struct {
	Lat float64 `json:"lat" binding:"required" example:"5.89980"`
	Lon float64 `json:"lon" binding:"required" example:"-2.03874"`
}

type BoundingBox struct {
	XMin float64 `json:"xmin" example:"-3.10904826799998091"`
	YMin float64 `json:"ymin" example:"4.7387739720000468"`
	XMax float64 `json:"xmax" example:"-1.39979249599997502"`
	YMax float64 `json:"ymax" example:"6.15031057300007"`
}

func PointToWKTPolygon(coord Coord, halfSize float64) string {
	// Earth's radius in meters
	const earthRadius = 6378137.0

	// Convert meters to degrees
	dLat := halfSize / earthRadius
	dLon := halfSize / (earthRadius * math.Cos(math.Pi*coord.Lat/180))

	// Offsets in degrees
	latOffset := dLat * (180 / math.Pi)
	lonOffset := dLon * (180 / math.Pi)

	// Calculate bounding box
	latMin := coord.Lat - latOffset
	latMax := coord.Lat + latOffset
	lonMin := coord.Lon - lonOffset
	lonMax := coord.Lon + lonOffset

	// Create WKT polygon
	return fmt.Sprintf("POLYGON((%f %f, %f %f, %f %f, %f %f, %f %f))",
		lonMin, latMin, lonMax, latMin, lonMax, latMax, lonMin, latMax, lonMin, latMin)
}

// CoordToWKT returns a WKT Point representation of a coordinate
func CoordToWKT(coord Coord) string {
	return fmt.Sprintf("POINT(%f %f)", coord.Lon, coord.Lat)
}

func BBoxToWKT(bbox BoundingBox) string {
	minLon, minLat := bbox.XMin, bbox.YMin
	maxLon, maxLat := bbox.XMax, bbox.YMax

	return fmt.Sprintf(
		"POLYGON((%f %f, %f %f, %f %f, %f %f, %f %f))",
		minLon, minLat,
		maxLon, minLat,
		maxLon, maxLat,
		minLon, maxLat,
		minLon, minLat, // closing the ring
	)
}
