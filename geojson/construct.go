package geojson

import "encoding/json"

func CreateCoordinate(long, lat float64) Coordinate {
    return [2]float64{long, lat}
}


func CreatePoint(point Coordinate) Geometry {
    return Geometry{PointT, []Coordinate{point}}
}


func WrapFeature(geometry Geometry, props map[string]any) Feature {
    var temp map[string]any
    if props != nil {
        temp = make(map[string]any, len(props))
        for k,v := range props {
            temp[k] = v
        }
    }
    return Feature{geometry, temp}
}

type fShell struct {
    Typ GeoType `json:"type"`
    Geometry `json:"geometry"`
    Props map[string]any `json:"properties"`
}

func (f Feature) ToJSON() (string, error) {
    shell := fShell{FeatureT, f.Geometry, f.Props}
    out, err := json.Marshal(shell)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

func FeatureFromJSON(data string) (Feature, error) {
    var shell fShell
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return *new(Feature), err
    }
    return Feature{shell.Geometry, shell.Props}, nil
}



func CreateGeometryColl(geometries ...Geometry) GeometryCollection {
    return geometries
}

type gcShell struct {
    Typ GeoType `json:"type"`
    GeometryCollection `json:"geometries"`
}

func (gc GeometryCollection) ToJSON() (string, error) {
    shell := gcShell{GeoType(GeometryCollectionT), gc}
    out, err := json.Marshal(shell)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

func GeometryCollFromJSON(data string) (GeometryCollection, error) {
    var shell gcShell
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return nil, err
    }
    return shell.GeometryCollection, nil
}



func CreateFeatureColl(features ...Feature) FeatureCollection {
    return features
}

type fcShell struct {
    Typ GeoType `json:"type"`
    FeatureCollection `json:"features"`
}

func (fc FeatureCollection) ToJSON() (string, error) {
    shell := fcShell{FeatureCollectionT, fc}
    out, err := json.Marshal(shell)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

func FeatureCollFromJSON(data string) (FeatureCollection, error) {
    var shell fcShell
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return nil, err
    }
    return shell.FeatureCollection, nil
}