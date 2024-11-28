package graphcreator

import (
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
                    "size": conn.State,
                    "cars": conn.NCars})
            paths = append(paths, lastFeat)
        }
    }

    return paths, joints
}
