package geojson

type GeoType string

type GeometryType GeoType

const (
    PointT              GeometryType = "Point"
    LineStringT         GeometryType = "LineString"
    PolygonT            GeometryType = "Polygon"
)

const (
    MultiPointT         GeometryType = "MultiPoint"
    MultiLineStringT    GeometryType = "MultiLineString"
    MultiPolygonT       GeometryType = "MultiPolygon"
    GeometryCollectionT GeometryType = "GeometryCollection"
)

const (
    FeatureT           GeoType = "Feature"
    FeatureCollectionT GeoType = "FeatureCollection"
)

type IJSONable interface {
    ToJSON() (string, error)
}

type Coordinate [2]float64

type Geometry struct {
    GeometryType `json:"type"`
    Coords  []Coordinate `json:"coordinates"`
}

type MultiGeometry struct {
    GeometryType `json:"type"`
    CoordsSets [][]Coordinate `json:"coordinates"`
}

type GeometryCollection[G Geometry | MultiGeometry] []G

type Feature[G Geometry | MultiGeometry] struct {
    Geometry G `json:"geometry"`
    Props    map[string]any `json:"properties"`
}

type FeatureCollection[G Geometry | MultiGeometry] []Feature[G]
