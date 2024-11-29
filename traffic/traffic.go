package traffic

import (
	"log"
	"math/rand"
	"optitraffic/node"
)

const vehicle_speed = 0.1

type Vehicle struct {
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
	Vehicles        []*Vehicle
}

func NewTrafficManager(graph *node.Graph) TrafficManager {
	return TrafficManager{
		vehicle_next_id: 0,
		Graph:           graph,
		Vehicles:        []*Vehicle{},
	}
}

func (tm *TrafficManager) NewRandomVehicle() {
	i := int(rand.Float64() * float64(len(tm.Graph.Nodes)))
	tm.Vehicles = append(tm.Vehicles, &Vehicle{
		Id:       tm.vehicle_next_id,
		Speed:    vehicle_speed,
		Target:   nil,
		At:       tm.Graph.Nodes[i],
		Progress: 0,
		Route:    nil,
	})

	tm.vehicle_next_id += 1
}

func (v *Vehicle) AtTarget() bool {
	return v.Target == v.At
}

func (tm *TrafficManager) RandomTarget(v *Vehicle) {
	i := int(rand.Float64() * float64(len(tm.Graph.Nodes)))
	v.Target = tm.Graph.Nodes[i]
}

func (t *TrafficManager) Update(dt float64) {
	for _, v := range t.Vehicles {
		// assign new random target to vehicle if arrived at target
		if v.AtTarget() || v.Route == nil {
			if v.Target != nil {
				log.Printf("vehicle: %d arrived at target: %d", v.Id, v.Target.Id)
			}
			t.RandomTarget(v)
			v.Progress = 0
			// create path based on situation
			pf := node.NewPathfinder(t.Graph)
			route, err := pf.Path(v.At, v.Target)
			if err != nil {
				log.Printf("vehicle: %d cant path to target: %s", v.Id, err)
			} else {
				v.Route = route[1:]
			}
		}

		// update vehicle progress
		if len(v.Route) == 1 {
			dist := v.At.Pos.DistanceTo(v.Target.Pos)
			v.Progress += v.Speed * dt / dist
			if v.Progress >= 1.0 {
				v.At = v.Target
			}
		} else if len(v.Route) > 0 {
			dist := v.At.Pos.DistanceTo(v.Route[0].Pos)
			v.Progress += v.Speed * dt / dist
			if v.Progress >= 1.0 {
				v.At = v.Route[0]
				pf := node.NewPathfinder(t.Graph)
				new_route, err := pf.Path(v.At, v.Target)
				if err != nil {
					log.Printf("vehicle: %d cant path to target: %s", v.Id, err)
				} else {
					v.Route = new_route[1:]
				}
			}
		}
	}
}
