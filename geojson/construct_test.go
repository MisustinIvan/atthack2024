package geojson

import (
	"testing"
)

func TestCreateColl(t *testing.T) {
    points := []FlatGeometry{
        CreatePoint(CreateCoordinate(0, 50)),
        CreatePoint(CreateCoordinate(1, 50)),
        CreatePoint(CreateCoordinate(1, 51)),
    }
    geoColl := CreateGeometryColl(points...)
    json, err := geoColl.ToJSON()
    if err != nil {
        t.Fail()
    }
    print(json)
}

func TestFeature(t *testing.T) {
    point := CreatePoint(CreateCoordinate(0, 50))
    feature := Feature[FlatGeometry]{point, nil}
    json, err := feature.ToJSON()
    if err != nil {
        t.Fail()
    }
    print(json)
    callback, err := FeatureFromJSON[FlatGeometry](json)
    if err != nil {
        t.Fail()
    }
    if !areEqualFeature(feature, callback) {
        t.Errorf("result not equal")
    }
}

func Test(t *testing.T) {

}

func areEqualFeature(a, b Feature[FlatGeometry]) bool {
    aNil, bNil := a.Props == nil, b.Props == nil
    if aNil && !bNil {
        return false
    }
    if !(aNil || bNil) {
        for k, v := range a.Props {
            if b.Props[k] != v {
                return false
            }
        }
    }
    if a.Geometry.GeometryType != b.Geometry.GeometryType {
        return false
    }
    for i, v := range a.Geometry.SingleCoords {
        if b.Geometry.SingleCoords[i] != v {
            return false
        }
    }
    return true
}
