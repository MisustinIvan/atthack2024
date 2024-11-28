package graphconvertor

import (
	"errors"
	gj "optitraffic/geojson"
	"optitraffic/node"
)

func PosToCoord(pos node.Pos) gj.Coordinate {
    return gj.CreateCoordinate(pos.X, pos.Y)
}

func CoordToPos(coord gj.Coordinate) node.Pos {
    return node.Pos{X: coord[0], Y: coord[1]}
}

func GoOverGraph(graph node.Graph) (paths, joints gj.FeatureCollection[gj.Geometry]) {
    paths = make([]gj.Feature[gj.Geometry], 0)
    joints = make([]gj.Feature[gj.Geometry], 0)
    var (
        lastNode, lastConn gj.Geometry
        lastFeat gj.Feature[gj.Geometry]
    )

    for _, parent := range graph.Nodes {
        lastNode = gj.CreatePoint(PosToCoord(parent.Pos))
        joints = append(joints, gj.WrapFeature(lastNode,
            map[string]any{"id": parent.Id}))

        for node, conn := range parent.Conns {
            lastConn, _ =  gj.CreateLineString(
                PosToCoord(parent.Pos),
                PosToCoord(node.Pos))
            lastFeat = gj.WrapFeature(
                lastConn,
                map[string]any{
                    "state": conn.State,
                    "size": conn.Size,
                    "cars": conn.NCars})
            paths = append(paths, lastFeat)
        }
    }

    return paths, joints
}


type GeoNode struct {
    gj.Coordinate `json:"coordinates"`
    Id int      `json:"id"`
}

func PointToGeoNode(point gj.Feature[gj.Geometry]) (GeoNode, error) {
    if point.Geometry.GeometryType != gj.PointT {
        return *new(GeoNode), errors.New("not a point")
    }
    if val, ok := point.Props["id"]; !ok {
        return *new(GeoNode), errors.New("missing id")
    } else if id, ok := val.(int); !ok {
        return *new(GeoNode), errors.New("not an id")
    } else {
        return GeoNode{gj.Coordinate(point.Geometry.Coords[0]), id}, nil
    }
}

func (g GeoNode) ToPoint() gj.Feature[gj.Geometry] {
    point := gj.CreatePoint(g.Coordinate)
    return gj.WrapFeature(point, map[string]any{"id":g.Id})
}


type GeoPath struct {
    Ends [2]gj.Coordinate   `json:"coordinates"`
    State node.ConnState    `json:"state"`
    Size int                `json:"size"`
    Cars int                `json:"cars"`
}

func LineStringToGeoPath(line gj.Feature[gj.Geometry]) (GeoPath, error) {
    if line.Geometry.GeometryType != gj.LineStringT {
        return *new(GeoPath), errors.New("line")
    }
    if len(line.Geometry.Coords) != 2 {
        return *new(GeoPath), errors.New("doesnt have only two ends")
    }

    var (
        state node.ConnState
        size, cars int
    )

    if val, ok := line.Props["state"]; !ok {
        return *new(GeoPath), errors.New("missing state")
    } else if state, ok = val.(node.ConnState); !ok {
        return *new(GeoPath), errors.New("not a state")
    }

    if val, ok := line.Props["size"]; !ok {
        return *new(GeoPath), errors.New("missing size")
    } else if size, ok = val.(int); !ok {
        return *new(GeoPath), errors.New("not a size")
    }

    if val, ok := line.Props["cars"]; !ok {
        return *new(GeoPath), errors.New("missing cars")
    } else if cars, ok = val.(int); !ok {
        return *new(GeoPath), errors.New("not a cars")
    }

    return GeoPath{[2]gj.Coordinate(line.Geometry.Coords), state, size, cars}, nil
}

func (g GeoPath) ToLineString() gj.Feature[gj.Geometry] {
    line, _ := gj.CreateLineString(g.Ends[0], g.Ends[1])
    return gj.WrapFeature(line, map[string]any{"state": g.State, "size": g.Size, "cars": g.Cars})
}