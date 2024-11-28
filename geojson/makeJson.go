package geojson

import (
	"encoding/json"
)

func (g Geometry) ToJSON() (string, error) {
    out, err := json.Marshal(g)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

// Creates a single geometry from JSON input
func GeometryFromJSON(data string) (Geometry, error) {
    var out Geometry
    err := json.Unmarshal([]byte(data), &out)
    if err != nil {
        return *new(Geometry), err
    }
    return out, nil
}


func (g MultiGeometry) ToJSON() (string, error) {
    out, err := json.Marshal(g)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

// Creates a multi-geometry from JSON input
func MultiGeometryFromJSON(data string) (MultiGeometry, error) {
    var out MultiGeometry
    err := json.Unmarshal([]byte(data), &out)
    if err != nil {
        return *new(MultiGeometry), err
    }
    return out, nil
}


type gcShell[G Geometry | MultiGeometry | MultiMultiGeometry] struct {
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

// Creates a collection geometry from JSON input
func GeometryCollFromJSON[G Geometry | MultiGeometry | MultiMultiGeometry](data string) (GeometryCollection[G], error) {
    var shell gcShell[G]
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return nil, err
    }
    return shell.GeometryCollection, nil
}


type fShell[G Geometry | MultiGeometry | MultiMultiGeometry] struct {
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

// Creates a feature of the geometry type from JSON input
func FeatureFromJSON[G Geometry | MultiGeometry | MultiMultiGeometry](data string) (Feature[G], error) {
    var shell fShell[G]
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return *new(Feature[G]), err
    }
    return Feature[G]{shell.Geometry, shell.Props}, nil
}


type fcShell[G Geometry | MultiGeometry | MultiMultiGeometry] struct {
    Typ GeoType `json:"type"`
    features []fShell[G] `json:"features"`
}

func (fc FeatureCollection[G]) ToJSON() (string, error) {
    features := make([]fShell[G], 0, len(fc))
    for _, v := range fc {
        features = append(features, fShell[G]{FeatureT, v.Geometry, v.Props})
    }
    shell := fcShell[G]{FeatureCollectionT, features}
    out, err := json.Marshal(shell)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

// Creates a collection of features of the geometry type from JSON input
func FeatureCollFromJSON[G Geometry | MultiGeometry | MultiMultiGeometry](data string) (FeatureCollection[G], error) {
    var shell fcShell[G]
    err := json.Unmarshal([]byte(data), &shell)
    if err != nil {
        return nil, err
    }
    FeatureColl := make([]Feature[G], 0, len(shell.features))
    for _, v := range shell.features {
        FeatureColl = append(FeatureColl, Feature[G]{v.Geometry, v.Props})
    }
    return FeatureColl, nil
}