package geojson

import (
	"testing"
)

func TestCreateGeometry(t *testing.T) {
    source, err := CreateLineString(CreateCoordinate(0, 50), CreateCoordinate(1, 50))
    if err != nil {
        t.Fatal(err.Error())
    }

    result, err := source.ToJSON()
    if err != nil {
        t.FailNow()
    }
    _, err = GeometryFromJSON(result)
    if err != nil {
        t.FailNow()
    }
}
