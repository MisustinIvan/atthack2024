package graphconvertor

import (
	gj "optitraffic/geojson"
	"testing"
)

func TestConverGeoNode(t *testing.T) {
    point := gj.CreatePoint(gj.CreateCoordinate(0, 50))

    result, err := PointToGeoNode(gj.WrapFeature(point, map[string]any{"id":0}))
    if err != nil {
        t.Fail()
    }

    if point.Coords[0][0] != result.Coordinate[0] || point.Coords[0][1] != result.Coordinate[1] {
        t.Fail()
    }
}