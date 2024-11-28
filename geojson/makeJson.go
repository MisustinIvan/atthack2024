package geojson

import (
	"encoding/json"
)

type gcShell[G Geometry | MultiGeometry] struct {
    Typ GeoType `json:"type"`
    GeometryCollection[G] `json:"geometries"`
}

func (gc GeometryCollection[G]) ToJSON() (string, error) {
    shell := gcShell[G]{GeoType(GeometryCollectionT), gc}
    out, err := json.Marshal(shell)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

func GeometryCollFromJSON[G Geometry | MultiGeometry](data string) (GeometryCollection[G], error) {
    var shell gcShell[G]
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return nil, err
    }
    return shell.GeometryCollection, nil
}


type fShell[G Geometry | MultiGeometry] struct {
    Typ GeoType `json:"type"`
    Geometry G `json:"geometry"`
    Props map[string]any `json:"properties"`
}

func (f Feature[G]) ToJSON() (string, error) {
    shell := fShell[G]{FeatureT, f.Geometry, f.Props}
    out, err := json.Marshal(shell)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

func FeatureFromJSON[G Geometry | MultiGeometry](data string) (Feature[G], error) {
    var shell fShell[G]
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return *new(Feature[G]), err
    }
    return Feature[G]{shell.Geometry, shell.Props}, nil
}


type fcShell[G Geometry | MultiGeometry] struct {
    Typ GeoType `json:"type"`
    FeatureCollection[G] `json:"features"`
}

func (fc FeatureCollection[G]) ToJSON() (string, error) {
    shell := fcShell[G]{FeatureCollectionT, fc}
    out, err := json.Marshal(shell)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

func FeatureCollFromJSON[G Geometry | MultiGeometry](data string) (FeatureCollection[G], error) {
    var shell fcShell[G]
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return nil, err
    }
    return shell.FeatureCollection, nil
}