package graphdb

import (
	"database/sql"
	"optitraffic/node"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestStoreGraph(t *testing.T) {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        t.Error(err.Error())
    }
    _, err = db.Exec("DELETE FROM graph_paths")
    _, err = db.Exec("DELETE FROM path_ends")
    _, err = db.Exec("DELETE FROM graph_nodes")
    if err != nil {
        t.Fatal(err.Error())
    }
    dao := NewDAO(db)

    poses := [...]node.Pos{
        {X: 0, Y: 50},
        {X: 1, Y: 50},
        {X: 1, Y: 51},
    }
    graph := node.NewGraph(poses[0])
    first, second, third := graph.Root, graph.NewNode(poses[1]), graph.NewNode(poses[2])
    graph.LinkBoth(first, second, 2); graph.LinkBoth(second, third, 1)

    err = dao.StoreGraph(graph)
    if err != nil {
        t.Error(err.Error())
    }
}