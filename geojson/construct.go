package geojson

import (
	"encoding/json"
	"errors"
)

func CreateCoordinate(long, lat float64) Coordinate {
    return [2]float64{long, lat}
}


func CreatePoint(point Coordinate) Geometry {
    return Geometry{PointT, []Coordinate{point}}
}

func CreateMultiPoint(points ...Geometry) (MultiGeometry, error) {
    if len(points) == 0 {
        return *new(MultiGeometry), errors.New("too few points")
    }
    temp := make([][]Coordinate, len(points))
    for _, v := range points {
        if v.GeometryType != PointT {
            return *new(MultiGeometry), errors.New("it's not all points")
        }
        temp = append(temp, v.Coords)
    }
    return MultiGeometry{MultiPointT, temp}, nil
}

func CreateLineString(joints ...Coordinate) (Geometry, error) {
    if len(joints) < 2 {
        return *new(Geometry), errors.New("a line has at least two points")
    }
    return Geometry{LineStringT, joints}, nil
}

func CreateMultiLineString(lineStrings ...Geometry) (MultiGeometry, error) {
    if len(lineStrings) == 0 {
        return *new(MultiGeometry), errors.New("no args bruh")
    }
    unfolded := make([][]Coordinate, 0)
    for _, v := range lineStrings {
        if v.GeometryType != LineStringT {
            return *new(MultiGeometry), errors.New("where the lines at")
        }
        unfolded = append(unfolded, v.Coords)
    }
    return MultiGeometry{MultiLineStringT, unfolded}, nil
}

func CreatePolygon(verticies ...Coordinate) (Geometry, error) {
    if len(verticies) < 4 {
        return *new(Geometry), errors.New("a polygon has at least three verticies")
    }
    if verticies[0] != verticies[len(verticies)-1] {
        return *new(Geometry), errors.New("the first and last vertex must be the same")
    }
    return Geometry{PolygonT, verticies}, nil
}

func CreateMultiPolygon(polygons ...Geometry) (MultiGeometry, error) {
    if len(polygons) == 0 {
        return *new(MultiGeometry), errors.New("nuh uh")
    }
    all := make([][]Coordinate, len(polygons))
    for _, v := range polygons {
        if v.GeometryType != PolygonT {
            return *new(MultiGeometry), errors.New("polygons must be polygons")
        }
        all = append(all, v.Coords)
    }
    return MultiGeometry{MultiPolygonT, all}, nil
}



func CreateGeometryColl[G Geometry | MultiGeometry](geometries ...G) GeometryCollection[G] {
    return geometries
}

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



func WrapFeature[G Geometry | MultiGeometry](geometry G, props map[string]any) Feature[G] {
    var temp map[string]any
    if props != nil {
        temp = make(map[string]any, len(props))
        for k,v := range props {
            temp[k] = v
        }
    }
    return Feature[G]{geometry, temp}
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



func CreateFeatureColl[G Geometry | MultiGeometry](features ...Feature[G]) FeatureCollection[G] {
    return features
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