package geojson

type GeoType string

type GeometryType GeoType

const (
    PointT GeometryType = "Point"
    MultiPointT GeometryType = "MultiPoint"
    LineStringT GeometryType = "LineString"
    MultiLineStringT GeometryType = "MultiLineString"
    MultiPolygonT GeometryType = "MultiPolygon"
    PolygonT GeometryType = "Polygon"
    GeometryCollectionT GeometryType = "GeometryCollection"
)

const (
    FeatureT GeoType = "Feature"
    FeatureCollectionT GeoType = "FeatureCollection"
)

type GeometryCollection struct {
    Items []Geometry    `json:"geometries"`
}

type FeatureCollection struct {
    Items []Feature `json:"features"`
}

type Feature struct {
    Geometry `json:"geometry"`
    Props map[string]any `json:"properties"`
}

type Geometry struct {
    GeoType `json:"type"`
    Coords []Coordinate `json:"coordinates"`
}

type Coordinate struct {
    Pos [2]float64
}
