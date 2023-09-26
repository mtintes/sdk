package nextroute

import "github.com/nextmv-io/sdk/connect"

// Arc is a directed connection between two nodes ([ModelStops]) that specifies
// that the origin stop must be planned before the destination stop on a
// vehicle's route.
type Arc interface {
	// Origin returns the origin node ([ModelStop]) of the arc.
	Origin() ModelStop
	// Destination returns the destination node ([ModelStop]) of the arc.
	Destination() ModelStop
}

// Arcs is a collection of [Arc]s.
type Arcs []Arc

// DirectedAcyclicGraph is a set of nodes (of type [ModelStop]) connected by
// arcs that does not contain cycles. It restricts the sequence in which the
// stops can be planned on the vehicle. An arc (u -> v) indicates that the stop
// u must be planned before the stop v on the vehicle's route.
type DirectedAcyclicGraph interface {
	// Arcs returns all [Arcs] in the graph.
	Arcs() Arcs

	// IndependentDirectedAcyclicGraphs returns all the independent
	// [DirectedAcyclicGraph]s in the graph. An independent
	// [DirectedAcyclicGraph] is a [DirectedAcyclicGraph] that does not share
	// any [ModelStop]s with any other [DirectedAcyclicGraph]s.
	IndependentDirectedAcyclicGraphs() ([]DirectedAcyclicGraph, error)

	// IsAllowed returns true if the sequence of stops is allowed by the DAG,
	// otherwise returns false.
	IsAllowed(stops ModelStops) (bool, error)

	// ModelStops returns all [ModelStops] in the graph.
	ModelStops() ModelStops
	// AddArc adds a new [Arc] in the graph if it was not already added. The new
	// [Arc] should not cause a cycle.
	AddArc(origin, destination ModelStop) error
	// OutboundArcs returns all [Arcs] that have the given [ModelStop] as their
	// origin.
	OutboundArcs(stop ModelStop) Arcs
}

// NewDirectedAcyclicGraph creates a new [DirectedAcyclicGraph].
func NewDirectedAcyclicGraph() DirectedAcyclicGraph {
	connect.Connect(con, &newDirectedAcyclicGraph)
	return newDirectedAcyclicGraph()
}
