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

func TurnGraphToGeoJSON(graph node.Graph) (paths gj.FeatureCollection[gj.Geometry], joints gj.FeatureCollection[gj.FlatGeometry]) {
    paths = make([]gj.Feature[gj.Geometry], 0)
    joints = make([]gj.Feature[gj.FlatGeometry], 0)
    var (
        lastNode gj.FlatGeometry
        lastConn gj.Geometry
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


func GeoJSONToGraph(paths gj.FeatureCollection[gj.Geometry], joints gj.FeatureCollection[gj.FlatGeometry]) (node.Graph, error) {
    geoPaths := LineCollToGeoPath(paths)
    geoNodes, e := PointsCollToGeoNode(joints)
    if e != nil {
        return *new(node.Graph), e
    }

    // collect relations
    rels := make(map[gj.Coordinate][]struct{
        Coordinate gj.Coordinate
        State node.ConnState
        Size, Cars int
    })
    var (
        exists bool
        ends [2]gj.Coordinate
        val struct{Coordinate gj.Coordinate; State node.ConnState; Size, Cars int}
    )
    for _, v := range geoPaths {
        ends = v.Ends
        val = struct{Coordinate gj.Coordinate; State node.ConnState; Size, Cars int}{
                ends[1], v.State, v.Size, v.Cars}
        _, exists = rels[ends[0]]
        if exists {
            rels[ends[0]] = append(rels[ends[0]], val)
        } else {
            rels[ends[0]] = []struct{Coordinate gj.Coordinate; State node.ConnState; Size, Cars int}{val}
        }
    }

    //half init nodes
    nodes := make([]*node.Node, 0, len(geoNodes))
    for _, v := range geoNodes {
        nodes = append(nodes, &node.Node{Id: v.Id, Pos: CoordToPos(v.Coordinate)})
    }
    // get conns for each node
    var (
        main gj.Coordinate
        mainPos, otherPos node.Pos
        target *node.Node
        connect map[*node.Node]node.ConnParams
    )
    for _, outer := range nodes {
        main = PosToCoord(outer.Pos)
        connect = make(map[*node.Node]node.ConnParams)
        // make conns
        for _, inner := range rels[main] {
            mainPos, otherPos = outer.Pos, CoordToPos(inner.Coordinate)
            // search for target to assign for
            for _, v := range nodes {
                if v.Pos == otherPos {
                    target = v
                    break
                }
            }
            // asign
            connect[target] = node.ConnParams{
                Dist: mainPos.DistanceTo(otherPos),
                Size: inner.Size,
                State: inner.State,
                NCars: inner.Cars,
            }
        }
        // assign conns
        outer.Conns = connect
    }

    // pray to Kerninghan
    return node.Graph{ Nodes: nodes, Root: nodes[0] }, nil
}


type GeoNode struct {
    Id int      `json:"id"`
    gj.Coordinate `json:"coordinates"`
}

func PointToGeoNode(point gj.Feature[gj.FlatGeometry]) (GeoNode, error) {
    if point.Geometry.GeometryType != gj.PointT {
        return *new(GeoNode), errors.New("not a point")
    }
    if val, ok := point.Props["id"]; !ok {
        return *new(GeoNode), errors.New("missing id")
    } else if id, ok := val.(float64); !ok {
        return *new(GeoNode), errors.New("not an id")
    } else {
        return GeoNode{Coordinate: point.Geometry.SingleCoords, Id: int(id)}, nil
    }
}

func (g GeoNode) ToPoint() gj.Feature[gj.FlatGeometry] {
    point := gj.CreatePoint(g.Coordinate)
    return gj.WrapFeature(point, map[string]any{"id":g.Id})
}

func PointsCollToGeoNode(collOfPoints gj.FeatureCollection[gj.FlatGeometry]) ([]GeoNode, error) {
    points := make([]GeoNode, 0, len(collOfPoints))
    var (
        lastN GeoNode
        e error
    )
    for _, v := range collOfPoints {
        lastN, e = PointToGeoNode(v)
        if e != nil {
            return nil, e
        }
        points = append(points, lastN)
    }
    return points, nil
}

func GeoNodesToPointsColl(nodes ...GeoNode) gj.FeatureCollection[gj.FlatGeometry] {
    points := make([]gj.Feature[gj.FlatGeometry], 0, len(nodes))
    for _, v := range nodes {
        points = append(points, v.ToPoint())
    }
    return points
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

func LineCollToGeoPath(collOfLines gj.FeatureCollection[gj.Geometry]) []GeoPath {
    paths := make([]GeoPath, 0, len(collOfLines))
    var (
        lastP GeoPath
        e error
    )
    for _, v := range collOfLines {
        lastP, e = LineStringToGeoPath(v)
        if e == nil {
            paths = append(paths, lastP)
        }
    }
    return paths
}

func GeoPathsToLineColl(paths ...GeoPath) gj.FeatureCollection[gj.Geometry] {
    lines := make([]gj.Feature[gj.Geometry], 0, len(paths))
    for _, v := range paths {
        lines = append(lines, v.ToLineString())
    }
    return lines
}