package graphconvertor

import (
	gj "optitraffic/geojson"
	"testing"
)

func TestConverGeoNode(t *testing.T) {
    point := gj.CreatePoint(gj.CreateCoordinate(0, 50))

    result, err := PointToGeoNode(gj.WrapFeature(point, map[string]any{"id": float64(0)}))
    if err != nil {
        t.Error(err.Error())
    }

    if point.SingleCoords[0] != result.Coordinate[0] || point.SingleCoords[1] != result.Coordinate[1] {
        t.Error("neco se nerovna")
    }
}

func TestConvertGeoNodesM(t *testing.T) {
    points := []gj.FlatGeometry{
        gj.CreatePoint(gj.CreateCoordinate(0, 50)),
        gj.CreatePoint(gj.CreateCoordinate(1, 50)),
        gj.CreatePoint(gj.CreateCoordinate(1, 51)),
    }
    feats := gj.CreateFeatureColl(
        gj.WrapFeature(points[0], map[string]any{"id":0.0}),
        gj.WrapFeature(points[1], map[string]any{"id":1.0}),
        gj.WrapFeature(points[2], map[string]any{"id":2.0}),
    )

    result, err := PointsCollToGeoNode(feats)
    if err != nil {
        t.Fatal(err.Error())
    }
    print(len(result))
}