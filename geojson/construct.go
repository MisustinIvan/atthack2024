package geojson

import (
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


func CreateFeatureColl[G Geometry | MultiGeometry](features ...Feature[G]) FeatureCollection[G] {
    return features
}
