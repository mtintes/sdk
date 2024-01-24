package nextroute

import (
	"context"
	"math/rand"

	"github.com/nextmv-io/sdk/connect"
)

// NewSolution creates a new solution. The solution is created from the given
// model. The solution starts with all plan units unplanned. Once a solution
// has been created the model can no longer be changed, it becomes immutable.
func NewSolution(
	m Model,
) (Solution, error) {
	connect.Connect(con, &newSolution)
	return newSolution(m)
}

// NewRandomSolution creates a new solution. The solution is created from the
// given model. The solution starts with an empty solution and will assign
// a random plan unit to a random vehicle. The remaining plan units
// are added to the solution in a random order at the best possible position.
func NewRandomSolution(
	ctx context.Context,
	m Model,
) (Solution, error) {
	connect.Connect(con, &newRandomSolution)
	return newRandomSolution(ctx, m)
}

// NewSweepSolution creates a new solution. The solution is created from the
// given model using a sweep construction heuristic.
func NewSweepSolution(
	ctx context.Context,
	m Model,
) (Solution, error) {
	connect.Connect(con, &newSweepSolution)
	return newSweepSolution(ctx, m)
}

// NewClusterSolution creates a new solution. The solution is created from the
// given model using a cluster construction heuristic. The number of clusters
// is maximized to the number of empty vehicles in the solution and the
// maximumClusters parameter.
func NewClusterSolution(
	ctx context.Context,
	m Model,
	maximumClusters int,
) (Solution, error) {
	connect.Connect(con, &newClusterSolution)
	return newClusterSolution(ctx, m, maximumClusters)
}

// Solution is a solution to a model.
type Solution interface {
	// BestMove returns the best move for the given solution plan unit. The
	// best move is the move that has the lowest score. If there are no moves
	// available for the given solution plan unit, a move is returned which
	// is not executable, SolutionMoveStops.IsExecutable.
	BestMove(context.Context, SolutionPlanUnit) SolutionMove

	// ConstraintData returns the data of the constraint for the solution. The
	// constraint data of a solution is set by the
	// ConstraintSolutionDataUpdater.UpdateConstraintSolutionData method of the
	// constraint.
	ConstraintData(constraint ModelConstraint) any
	// Copy returns a deep copy of the solution.
	Copy() Solution

	// FixedPlanUnits returns the solution plan units that are fixed.
	// Fixed plan units are plan units that are not allowed to be planned or
	// unplanned. The union of fixed, planned and unplanned plan units
	// is the set of all plan units in the model.
	FixedPlanUnits() ImmutableSolutionPlanUnitCollection

	// Model returns the model of the solution.
	Model() Model

	// ObjectiveData returns the value of the objective for the solution. The
	// objective value of a solution is set by the
	// ObjectiveSolutionDataUpdater.UpdateObjectiveSolutionData method of the
	// objective. If the objective is not set on the solution, nil is returned.
	ObjectiveData(objective ModelObjective) any
	// ObjectiveValue returns the objective value for the objective in the
	// solution. Also returns 0.0 if the objective is not part of the solution.
	ObjectiveValue(objective ModelObjective) float64

	// PlannedPlanUnits returns the solution plan units that are planned as
	// a collection of solution plan units.
	PlannedPlanUnits() ImmutableSolutionPlanUnitCollection

	// Random returns the random number generator of the solution.
	Random() *rand.Rand

	// Score returns the score of the solution.
	Score() float64
	// SolutionPlanStopsUnit returns the [SolutionPlanStopsUnit] for the given
	// model plan unit.
	SolutionPlanStopsUnit(planUnit ModelPlanStopsUnit) SolutionPlanStopsUnit
	// SolutionPlanUnit returns the [SolutionPlanUnit] for the given
	// model plan unit.
	SolutionPlanUnit(planUnit ModelPlanUnit) SolutionPlanUnit
	// SolutionStop returns the solution stop for the given model stop.
	SolutionStop(stop ModelStop) SolutionStop
	// SolutionVehicle returns the solution vehicle for the given model vehicle.
	SolutionVehicle(vehicle ModelVehicle) SolutionVehicle

	// UnPlannedPlanUnits returns the solution plan units that are not
	// planned.
	UnPlannedPlanUnits() ImmutableSolutionPlanUnitCollection

	// Vehicles returns the vehicles of the solution.
	Vehicles() SolutionVehicles
}

// Solutions is a slice of solutions.
type Solutions []Solution
