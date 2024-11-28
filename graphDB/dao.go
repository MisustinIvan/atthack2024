package graphdb

import (
	"database/sql"
	gj "optitraffic/geojson"
	conv "optitraffic/graphConvertor"
	"optitraffic/node"

	_ "github.com/mattn/go-sqlite3"
)

type GraphDAO interface {
    GetAllPoints() gj.FeatureCollection[gj.Geometry]
    GetAllNodes() []node.Node
    GetNodeById(id int) (node.Node, error)
    GetAllPaths() gj.FeatureCollection[gj.Geometry]
    GetGraph() (paths, points gj.FeatureCollection[gj.Geometry])

    StoreGraph(node.Graph) error
    StoreGeoNodes(...conv.GeoNode) error
    StoreGeoPaths(...conv.GeoPath) error
}

type SQLiteDAO struct {
    db *sql.DB
}

func NewDAO(db *sql.DB) SQLiteDAO {
    return SQLiteDAO{db}
}

func (dao *SQLiteDAO) StoreGraph(graph node.Graph) error {
    paths, points := conv.GoOverGraph(graph)

    // Extract nodes
    pointArgs := conv.PointsCollToGeoNode(points)
    // Do the thingy
    if err := dao.StoreGeoNodes(pointArgs...); err != nil {
        return err
    }

    // Extract paths
    pathsArgs := conv.LineCollToGeoPath(paths)
    // Do the thingy
    if err := dao.StoreGeoPaths(pathsArgs...); err != nil {
        return err
    }

    return nil
}

func (dao *SQLiteDAO) StoreGeoNodes(nodes ...conv.GeoNode) error {
    // build query
    pointQuery := "INSERT INTO graph_nodes (id, longitude, latitude) VALUES\r\n"
    orgLn := len(pointQuery)
    for i := 0; i < len(nodes); i++ {
        pointQuery += "(?,?,?),\r\n"
    }
    if len := len(pointQuery); len > orgLn {
        pointQuery = pointQuery[:len-3]
    }
    // prepare data
    pointDeconstruct := make([]any, 0, len(nodes)*3)
    for _, v := range nodes {
        pointDeconstruct = append(pointDeconstruct, v.Id, v.Coordinate[0], v.Coordinate[1])
    }
    // execute query
    _, err := dao.db.Exec(pointQuery, pointDeconstruct...)
    return err
}

func (dao *SQLiteDAO) StoreGeoPaths(paths ...conv.GeoPath) error {
    // get last id (to assign coords)
    lastRow := dao.db.QueryRow("SELECT id FROM graph_paths ORDER BY id DESC")
    var lastID int
    if err := lastRow.Scan(&lastID); err != nil {
        lastID = 0
    }
    // prepare data query
    pathDataQuery := "INSERT INTO graph_paths (state, size, cars) VALUES\r\n"
    argLen := len(paths)
    for i := 0; i < argLen; i++ {
        pathDataQuery += "(?,?,?),\r\n"
    }
    if qLen := len(pathDataQuery); qLen > argLen {
        pathDataQuery = pathDataQuery[:qLen-3]
    }
    // prepare data data
    pathDataDecons := make([]any, 0, argLen*3)
    for _, v := range paths {
        pathDataDecons = append(pathDataDecons, v.State, v.Size, v.Cars)
    }
    // execute data query
    _, err := dao.db.Exec(pathDataQuery, pathDataDecons...)
    if err != nil {
        return err
    }

    // prepare coordinates query
    pathCoordQuery := "INSERT INTO path_ends (parent_id, longitude, latitude) VALUES\r\n"
    coordLen := argLen*2
    for i := 0; i < coordLen; i++ {
        pathCoordQuery += "(?,?,?),\r\n"
    }
    if qLen := len(pathCoordQuery); qLen > coordLen {
        pathCoordQuery = pathCoordQuery[:qLen-3]
    }
    // prepare coordinates data
    pathTuples := make([]any, 0, coordLen)
    for _, v := range paths {
        lastID++
        pathTuples = append(pathTuples, lastID, v.Ends[0][0], v.Ends[0][1])
        pathTuples = append(pathTuples, lastID, v.Ends[1][0], v.Ends[1][1])
    }
    // execute coordinates query
    _, err = dao.db.Exec(pathCoordQuery, pathTuples...)
    if err != nil {
        return err
    }

    return nil
}
