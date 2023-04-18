package nextroute

import (
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewModel creates a new empty vehicle routing model. Please use [BuildModel]
// if you want a model which already has all features added to it.
func NewModel() (Model, error) {
	connect.Connect(con, &newModel)
	return newModel()
}

// Model defines routing problem.
type Model interface {
	ModelData
	SolutionObserved

	// AddConstraint adds a constraint to the model. The constraint is
	// checked at the specified violation.
	AddConstraint(constraint ModelConstraint) error

	// Constraints returns all constraints of the model.
	Constraints() ModelConstraints

	// ConstraintsCheckedAt returns all constraints of the model that
	// are checked at the specified time of having calculated the new
	// information for the changed solution.
	ConstraintsCheckedAt(violation CheckedAt) ModelConstraints

	// DistanceUnit returns the unit of distance used in the model. The
	// unit is used to convert distances to values and vice versa. This is
	// also used for reporting.
	DistanceUnit() common.DistanceUnit

	// DurationUnit returns the unit of duration used in the model. The
	// unit is used to convert durations to values and vice versa. This is
	// also used for reporting.
	DurationUnit() time.Duration

	// DurationToValue converts the specified duration to a value as it used
	// internally in the model.
	DurationToValue(duration time.Duration) float64

	// Epoch returns the epoch of the model. The epoch is used to convert
	// time.Time to float64 and vice versa. All float64 values are relative
	// to the epoch.
	Epoch() time.Time

	// Expressions returns all expressions of the model for which a solution
	// has to calculate values. The expressions are sorted by their index. The
	// constraints register their expressions with the model.
	Expressions() ModelExpressions

	// IsLocked returns true if the model is locked. The model is
	// locked after a solution has been created using the model.
	IsLocked() bool

	// NewPlanSequence creates a new plan sequence. A plan sequence
	// is a plan cluster of a collection of stops. A plan cluster is a
	// collection of stops which are always planned and unplanned as a
	// single entity. In this case they have to be planned as a sequence on
	// the same vehicle.
	NewPlanSequence(stops ModelStops) (ModelPlanMultipleStops, error)
	// NewPlanSingleStop creates a new plan single stop. A plan single stop
	// is a plan cluster of a single stop. A plan cluster is a collection of
	// stops which are always planned and unplanned as a single entity.
	NewPlanSingleStop(stop ModelStop) (ModelPlanSingleStop, error)
	// NewPlanMultipleStops creates a new plan of multiple [ModelStops]. A plan
	// of multiple stops is a [ModelPlanCluster] of more than one stop. A plan
	// cluster is a collection of stops which are always planned and unplanned
	// as a single entity. When planned, they are always assigned to the same
	// vehicle. The function takes in a sequence represented by a
	// [DirectedAcyclicGraph] (DAG) which restricts the order in which the
	// stops can be planned on the vehicle. Using an empty DAG means that the
	// stops can be planned in any order and they will always be assigned to
	// the same vehicle. Consider the stops [s1, s2, s3] and the sequence [s1
	// -> s2, s1 -> s3]. This means that we are restricting that the stop s1
	// must come before s2 and s3. However, we are not specifying the order of
	// s2 and s3. This means that we can plan s2 before s3 or s3 before s2.
	NewPlanMultipleStops(
		stops ModelStops,
		sequence DirectedAcyclicGraph,
	) (ModelPlanMultipleStops, error)

	// NewStop creates a new stop. The stop is used to create plan clusters.
	NewStop(location common.Location) (ModelStop, error)

	// NewVehicle creates a new vehicle. The vehicle is used to create
	// solutions.
	NewVehicle(
		vehicleType ModelVehicleType,
		start time.Time,
		first ModelStop,
		last ModelStop,
	) (ModelVehicle, error)
	// NewVehicleType creates a new vehicle type. The vehicle type is used
	// to create vehicles.
	NewVehicleType(
		travelDuration DurationExpression,
		processDuration DurationExpression,
	) (ModelVehicleType, error)

	// NumberOfStops returns the number of stops in the model.
	NumberOfStops() int

	// Objective returns the objective of the model.
	Objective() ModelObjectiveSum

	// PlanClusters returns all plan clusters of the model. A plan cluster
	// is a collection of stops which are always planned and unplanned as a
	// single entity.
	PlanClusters() ModelPlanClusters

	// Random returns a random number generator.
	Random() *rand.Rand

	// SetRandom sets the random number generator of the model.
	SetRandom(random *rand.Rand)

	// SetSeed sets the seed of the random number generator of the model.
	SetSeed(seed int64)

	// Stops returns all stops of the model.
	Stops() ModelStops

	// Stop returns the stop with the specified index.
	Stop(index int) (ModelStop, error)

	// TimeFormat returns the time format used for reporting.
	TimeFormat() string

	// Vehicles returns all vehicles of the model.
	Vehicles() ModelVehicles
	// VehicleTypes returns all vehicle types of the model.
	VehicleTypes() ModelVehicleTypes

	// Vehicle returns the vehicle with the specified index.
	Vehicle(index int) ModelVehicle

	// MaxTime returns the maximum end time (upper bound) for any stop. This
	// function uses the [Model.Epoch()] as a starting point and adds a large
	// number to provide a large enough upper bound.
	MaxTime() time.Time
}