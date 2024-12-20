package geojson

import (
	"errors"
)

// Creates a coordinate
func CreateCoordinate(long, lat float64) Coordinate {
    return [2]float64{long, lat}
}


// Creates a point
func CreatePoint(point Coordinate) FlatGeometry {
    return FlatGeometry{PointT, point}
}

// Creates a group of points
func CreateMultiPoint(points ...FlatGeometry) (Geometry, error) {
    if len(points) == 0 {
        return *new(Geometry), errors.New("too few points")
    }
    temp := make([]Coordinate, len(points))
    for _, v := range points {
        if v.GeometryType != PointT {
            return *new(Geometry), errors.New("it's not all points")
        }
        temp = append(temp, v.SingleCoords)
    }
    return Geometry{MultiPointT, temp}, nil
}


// Creates a lines between points
func CreateLineString(joints ...Coordinate) (Geometry, error) {
    if len(joints) < 2 {
        return *new(Geometry), errors.New("a line has at least two points")
    }
    return Geometry{LineStringT, joints}, nil
}

// Creates a group of lines between own points
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


// Creates lines that connect into a polygon
func CreatePolygon(verticies ...Coordinate) (MultiGeometry, error) {
    if len(verticies) < 4 {
        return *new(MultiGeometry), errors.New("a polygon has at least three verticies")
    }
    if verticies[0] != verticies[len(verticies)-1] {
        return *new(MultiGeometry), errors.New("the first and last vertex must be the same")
    }
    return MultiGeometry{PolygonT, [][]Coordinate{verticies}}, nil
}

// Creates a group of polygons
func CreateMultiPolygon(polygons ...MultiGeometry) (MultiMultiGeometry, error) {
    if len(polygons) == 0 {
        return *new(MultiMultiGeometry), errors.New("nuh uh")
    }
    all := make([][][]Coordinate, len(polygons))
    for _, v := range polygons {
        if v.GeometryType != PolygonT {
            return *new(MultiMultiGeometry), errors.New("polygons must be polygons")
        }
        all = append(all, v.CoordsSets)
    }
    return MultiMultiGeometry{MultiPolygonT, all}, nil
}


// Creates a collection of same-type different geometric kinds
func CreateGeometryColl[G FlatGeometry | Geometry | MultiGeometry | MultiMultiGeometry](geometries ...G) GeometryCollection[G] {
    return geometries
}


// Wraps a geometric type as a Feature
func WrapFeature[G FlatGeometry | Geometry | MultiGeometry | MultiMultiGeometry](geometry G, props map[string]any) Feature[G] {
    var temp map[string]any
    if props != nil {
        temp = make(map[string]any, len(props))
        for k,v := range props {
            temp[k] = v
        }
    }
    return Feature[G]{geometry, temp}
}


// Creates a collection of different features with the same underlying geometric type
func CreateFeatureColl[G FlatGeometry | Geometry | MultiGeometry | MultiMultiGeometry](features ...Feature[G]) FeatureCollection[G] {
    return features
}
