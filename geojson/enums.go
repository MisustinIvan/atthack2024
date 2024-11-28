package geojson

// All GeoJSON types
type GeoType string

// All Geometric types
type GeometryType GeoType

// Single Geometric types
const (
    PointT              GeometryType = "Point"
    LineStringT         GeometryType = "LineString"
    PolygonT            GeometryType = "Polygon"
)

// Multi-Geometric types
const (
    MultiPointT         GeometryType = "MultiPoint"
    MultiLineStringT    GeometryType = "MultiLineString"
    MultiPolygonT       GeometryType = "MultiPolygon"
)

// Collection of any geometric types
const GeometryCollectionT GeometryType = "GeometryCollection"

// Features
const (
    FeatureT           GeoType = "Feature"
    FeatureCollectionT GeoType = "FeatureCollection"
)

// A random interface
type IJSONable interface {
    ToJSON() (string, error)
}

// Coordinates in format: longtitude, latitude
type Coordinate [2]float64

type FlatGeometry struct {
    GeometryType `json:"type"`
    SingleCoords  Coordinate `json:"coordinates"`
}

// Single geometry struct
type Geometry struct {
    GeometryType `json:"type"`
    Coords  []Coordinate `json:"coordinates"`
}

// Multi-geometry struct
type MultiGeometry struct {
    GeometryType `json:"type"`
    CoordsSets [][]Coordinate `json:"coordinates"`
}

type MultiMultiGeometry struct {
    GeometryType `json:"type"`
    CoordsSetsSets [][][]Coordinate `json:"coordinates"`
}

// Collection of geometries
type GeometryCollection[G Geometry | MultiGeometry | MultiMultiGeometry] []G

// A wrapper around a geometry type
type Feature[G Geometry | MultiGeometry | MultiMultiGeometry] struct {
    Geometry G `json:"geometry"`
    Props    map[string]any `json:"properties"`
}

// Collection of geometries
type FeatureCollection[G Geometry | MultiGeometry | MultiMultiGeometry] []Feature[G]
