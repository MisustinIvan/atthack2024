package traffic

import (
	"log"
	"math/rand"
	"optitraffic/geojson"
	"optitraffic/node"
)

const vehicle_speed = 0.1
const min_delta_perc = 0.05

type Vehicle struct {
    Type     node.VehicleType
    Id       int
    Speed    float64
    Target   *node.Node
    At       *node.Node
    Progress float64
    Route    []*node.Node
}

type TrafficManager struct {
    vehicle_next_id int
    Graph           *node.Graph
    Pathfinder      node.Pathfinder
    Vehicles        []*Vehicle
}

func NewTrafficManager(graph *node.Graph) TrafficManager {

    return TrafficManager{
        vehicle_next_id: 0,
        Graph:           graph,
        Pathfinder:      node.NewPathfinder(graph),
        Vehicles:        []*Vehicle{},
    }
}

func (t *TrafficManager) NewRandomVehicle(typ node.VehicleType) {
    i := int(rand.Float64() * float64(len(t.Graph.Nodes)))
    new := &Vehicle{
        Type:     typ,
        Id:       t.vehicle_next_id,
        Target:   nil,
        At:       t.Graph.Nodes[i],
        Progress: 0,
        Route:    nil,
    }
    switch typ {
    case node.NormalVehicle:
        new.Speed = vehicle_speed
    case node.EmergencyVehicle:
        new.Speed = vehicle_speed * 2
    }
    t.Vehicles = append(t.Vehicles, new)

    t.vehicle_next_id += 1
}

func (t *TrafficManager) Repath(v *Vehicle) error {
    pf := node.NewPathfinder(t.Graph)
    newRoute, err := pf.Path(v.At, v.Target, v.Type)
    if err != nil {
        return err
    }
    v.Route = newRoute[1:]

    //delta := 1 - (node.PathWeightedLength(v.Route, v.Type) / node.PathWeightedLength(newRoute, v.Type))
    //if delta > min_delta_perc {
    //  v.Route = newRoute[1:]
    //}
    return nil
}

func (v *Vehicle) AtTarget() bool {
    return v.Target == v.At
}

func (t *TrafficManager) RandomTarget(v *Vehicle) {
    i := int(rand.Float64() * float64(len(t.Graph.Nodes)))
    v.Target = t.Graph.Nodes[i]
}

func (v *Vehicle) InterpolatePos() node.Pos {
    if len(v.Route) == 0 {
        return v.At.Pos
    }

    diff := v.At.Pos.Diff(v.Route[0].Pos)

    return node.Pos{
        X: v.At.Pos.X - diff.X*v.Progress,
        Y: v.At.Pos.Y - diff.Y*v.Progress,
    }
}

func (t *TrafficManager) VehiclesAsPoints() geojson.FeatureCollection[geojson.FlatGeometry] {
    res := geojson.FeatureCollection[geojson.FlatGeometry]{}
    for _, v := range t.Vehicles {
        pos := v.InterpolatePos()
        res = append(res, geojson.Feature[geojson.FlatGeometry]{
            Geometry: geojson.FlatGeometry{
                GeometryType: geojson.PointT,
                SingleCoords: [2]float64{pos.X, pos.Y},
            },
            Props: map[string]any{
                "type": v.Type,
            },
        })
    }
    return res
}

// Updates the state of all vehicles on the graph
func (t *TrafficManager) Update(dt float64) {
    // Regular traffic lights
    //get ignored nodes
    forbiddenNodes := make([]*node.Node, 0, len(t.Vehicles)*3)
    var (
        leng int
        curr *node.Node
    )
    for _, v := range t.Vehicles {
        leng = min(3, len(v.Route))
        for i := 0; i < leng; i++ {
            curr = v.Route[i]
            if !containsNode(forbiddenNodes, curr) {
                forbiddenNodes = append(forbiddenNodes, curr)
            }
        }
    }
    //generation of traffic lights
    var (
        tempC node.ConnParams
    )
    for _, curr := range t.Graph.Nodes {
        if containsNode(forbiddenNodes, curr) {continue}
        for k := range curr.Conns {
            tempC = k.Conns[curr]
            k.Conns[curr] = node.ConnParams{Dist: tempC.Dist, Size: tempC.Size, State: nextTrafficState(), NCars: tempC.NCars}
        }
    }


    // Update route
    for _, v := range t.Vehicles {
        // assign new random target to vehicle if arrived at target
        if v.AtTarget() || v.Route == nil {
            if v.Target != nil {
                log.Printf("vehicle: %d arrived at target: %d", v.Id, v.Target.Id)
            }
            t.RandomTarget(v)
            v.Progress = 0.0
            // create path based on situation
            err := t.Repath(v)
            if err != nil {
                log.Printf("vehicle: %d cant path to target: %s", v.Id, err)
            }
        }

        // update vehicle progress
        if len(v.Route) == 1 {
            dist := v.At.Pos.DistanceTo(v.Target.Pos)
            v.Progress += v.Speed * dt / dist
            if v.Progress >= 1.0 {
                v.Progress = 0.0
                v.At = v.Target
            }
        } else if len(v.Route) > 0 {
            dist := v.At.Pos.DistanceTo(v.Route[0].Pos)
            v.Progress += v.Speed * dt / dist
            if v.Progress >= 1.0 {
                v.Progress = 0.0
                v.At = v.Route[0]
                err := t.Repath(v)
                if err != nil {
                    log.Printf("vehicle: %d cant path to target: %s", v.Id, err)
                }
            }
        }
    }

    // Override traffics after emergency vehicles
    for _, v := range t.Vehicles {
        if v.Type != node.EmergencyVehicle { continue }
        var (
            prev, curr, next *node.Node
            tempC node.ConnParams
        )
        prev = v.At
        routLen := len(v.Route)
        leng := min(3, routLen)
        for i := 0 ; i < leng; i++ {
            curr = v.Route[i]
            if i + 1 < routLen { next = v.Route[i + 1]} else { next = nil}
            for k := range curr.Conns{
                tempC = k.Conns[curr]
                if next != nil && k == next {
                    k.Conns[curr] = node.ConnParams{Dist: tempC.Dist, Size: tempC.Size, State: node.Open, NCars: tempC.NCars}
                } else if k == prev {
                    k.Conns[curr] = node.ConnParams{Dist: tempC.Dist, Size: tempC.Size, State: node.Danger, NCars: tempC.NCars}
                } else {
                    k.Conns[curr] = node.ConnParams{Dist: tempC.Dist, Size: tempC.Size, State: node.Closed, NCars: tempC.NCars}
                }
            }
        }
    }
}

func containsNode(nodes []*node.Node, target *node.Node) bool {
    for _, v := range nodes {
        if v == target {
            return true
        }
    }
    return false
}

// Returns either Open or Closed
func nextTrafficState() node.ConnState {
    if rand.Float64() > 0.5 {
        return node.Open
    } else {
        return node.Closed
    }
}
