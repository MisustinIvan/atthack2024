package graphdb

import (
	"database/sql"
	"errors"
	gj "optitraffic/geojson"
	conv "optitraffic/graphConvertor"
	"optitraffic/node"

	_ "github.com/mattn/go-sqlite3"
)

type GraphDAO interface {
    GetAllPoints() (gj.FeatureCollection[gj.Geometry], error)
    GetAllPaths() (gj.FeatureCollection[gj.Geometry], error)
    GetGraph() (paths, points gj.FeatureCollection[gj.Geometry], err error)

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


func (dao *SQLiteDAO) GetGraph() (paths gj.FeatureCollection[gj.Geometry], points gj.FeatureCollection[gj.FlatGeometry], err error) {
    points, err = dao.GetAllPoints()
    if err != nil {
        return nil, nil, err
    }
    paths, err = dao.GetAllPaths()
    if err != nil {
        return nil, nil, err
    }
    return paths, points, nil
}

func (dao *SQLiteDAO) GetAllPoints() (gj.FeatureCollection[gj.FlatGeometry], error) {
    out, err := dao.db.Query("SELECT id, longitude, latitude FROM graph_nodes ORDER BY id")
    if err != nil {
        return nil, err
    }
    defer out.Close()

    nodesOut := make([]conv.GeoNode, 0)
    var (
        lastID int
        lastLong, lastLat float64
    )
    for out.Next() {
        if err := out.Scan(&lastID, &lastLong, &lastLat); err != nil {
            return nil, err
        }
        nodesOut = append(nodesOut, conv.GeoNode{Id: lastID, Coordinate: gj.CreateCoordinate(lastLong, lastLat)})
    }

    if err := out.Err(); err != nil {
        return nil, err
    }
    return conv.GeoNodesToPointsColl(nodesOut...), nil
}

func (dao *SQLiteDAO) GetAllPaths() (gj.FeatureCollection[gj.Geometry], error) {
    // Magic
    out, err := dao.db.Query(`
        SELECT graph_paths.id, path_ends.longitude, path_ends.latitude,
            graph_paths.state, graph_paths.size, graph_paths.cars
        FROM graph_paths INNER JOIN
            path_ends ON graph_paths.id = path_ends.parent_id
        ORDER BY graph_paths.id, path_ends.id
    `)
    if err != nil {
        return nil, err
    }
    defer out.Close()

    // extract data
    coords, data := make([][2]gj.Coordinate, 0), make([]struct{state, size, cars int}, 0)
    var (
        lastID, currID, coordTupleIndex, i int
        lastLong, LastLat float64
        lastState, lastSize, lastCars int
    )
    for out.Next() {
        if err := out.Scan(&currID, &lastLong, &LastLat, &lastState, &lastSize, &lastCars); err != nil {
            return nil, err
        }
        if currID != lastID {
            data = append(data, struct{state, size, cars int}{
                lastState, lastSize, lastCars})
            coords = append(coords, [2]gj.Coordinate{})
            lastID = currID
            coordTupleIndex--
        } else { coordTupleIndex++ }
        coords[i][coordTupleIndex] = gj.CreateCoordinate(lastLong, LastLat)
        i++
    }
    if err := out.Err(); err != nil {
        return nil, err
    }

    // aggregate to GeoPaths
    leng := len(data)
    if leng != len(coords) {
        return nil, errors.New("neco se posralo")
    }
    paths := make([]conv.GeoPath, 0, leng)
    var dat struct{state, size, cars int}
    for i := 0; i < leng; i++ {
        dat = data[i]
        paths = append(paths, conv.GeoPath{Ends: coords[i],
            State: node.ConnState(dat.state), Size: dat.size, Cars: dat.cars})
    }

    return conv.GeoPathsToLineColl(paths...), nil
}



func (dao *SQLiteDAO) StoreGraph(graph node.Graph) error {
    paths, points := conv.TurnGraphToGeoJSON(graph)

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
