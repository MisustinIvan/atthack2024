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
    GetNodeById(id int) ([]node.Node, error)
    GetAllPaths() gj.FeatureCollection[gj.Geometry]

    StoreNodes([]*node.Node) error
    StoreGraph(node.Graph) error
    StoreGeoNodes(...conv.GeoNode)error
    StoreGeoPaths(...conv.GeoPath) error
}

type SQLiteDAO struct {
    db *sql.DB
}
