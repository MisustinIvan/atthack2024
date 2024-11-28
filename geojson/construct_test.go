package geojson

import "testing"

func TestCreateColl(t *testing.T) {
    points := []Geometry{
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