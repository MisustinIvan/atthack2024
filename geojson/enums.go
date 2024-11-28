package geojson

type GeoType string

type GeometryType GeoType

const (
    PointT              GeometryType = "Point"
    MultiPointT         GeometryType = "MultiPoint"
    LineStringT         GeometryType = "LineString"
    MultiLineStringT    GeometryType = "MultiLineString"
    MultiPolygonT       GeometryType = "MultiPolygon"
    PolygonT            GeometryType = "Polygon"
    GeometryCollectionT GeometryType = "GeometryCollection"
)

const (
    FeatureT           GeoType = "Feature"
    FeatureCollectionT GeoType = "FeatureCollection"
)

type GeometryCollection []Geometry

type FeatureCollection []Feature

type Feature struct {
    Geometry `json:"geometry"`
    Props    map[string]any `json:"properties"`
}

type Geometry struct {
    GeometryType `json:"type"`
    Coords  []Coordinate `json:"coordinates"`
}

type Coordinate [2]float64
