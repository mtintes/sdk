// package main holds the implementation of the shift scheduling template.
package main

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/run/schema"
	"github.com/nextmv-io/sdk/run/statistics"
)

const gap = 0.999

func main() {
	runner := run.CLI(solver,
		run.InputValidate[run.CLIRunnerConfig, Input, Options, schema.Output](
			nil,
		),
	)
	err := runner.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func solver(_ context.Context, input Input, opts Options) (out schema.Output, retErr error) {
	// We solve a shift coverage problem using Mixed Integer Programming.
	// We solve this by generating all possible shifts
	// and then selecting a subset of these
	potentialAssignments, potentialAssignmentsPerWorker := potentialAssignments(input, opts)
	demands := demands(input, potentialAssignments)
	m, x := newMIPModel(input, potentialAssignments, potentialAssignmentsPerWorker, demands, opts)

	solver, err := mip.NewSolver("highs", m)
	if err != nil {
		return schema.Output{}, err
	}
	options := mip.NewSolveOptions()
	options.SetVerbosity(mip.Off)
	err = options.SetMIPGapRelative(gap)
	if err != nil {
		return schema.Output{}, err
	}
	err = options.SetMaximumDuration(opts.SolverDuration)
	if err != nil {
		return schema.Output{}, err
	}
	solution, err := solver.Solve(options)
	if err != nil {
		return schema.Output{}, err
	}

	output, err := format(solution, input, x, potentialAssignments)
	if err != nil {
		return schema.Output{}, err
	}

	return output, nil
}

func format(
	solution mip.Solution,
	_ Input,
	x model.MultiMap[mip.Bool, Assignment],
	assignments []Assignment,
) (output schema.Output, err error) {
	o := schema.Output{}

	o.Version = schema.Version{
		Sdk: sdk.VERSION,
	}

	stats := statistics.NewStatistics()

	result := statistics.Result{}

	run := statistics.Run{}

	t := solution.RunTime().Seconds()
	run.Duration = &t
	result.Duration = &t

	nextShiftSolution := Output{}

	nextShiftSolution.Status = "infeasible"

	if solution != nil && solution.HasValues() {
		nextShiftSolution.Status = "suboptimal"
		if solution.IsOptimal() {
			nextShiftSolution.Status = "optimal"
		}

		nextShiftSolution.Value = solution.ObjectiveValue()
		val := statistics.Float64(solution.ObjectiveValue())
		result.Value = &val

		usedWorkers := make(map[int]struct{})

		for _, assignment := range assignments {
			if solution.Value(x.Get(assignment)) >= 0.9 {
				nextShiftSolution.AssignedShifts = append(nextShiftSolution.AssignedShifts, OutputAssignment{
					Start:    assignment.Start,
					End:      assignment.End,
					WorkerID: assignment.Worker.ID,
				})
				if _, ok := usedWorkers[assignment.Worker.ID]; !ok {
					usedWorkers[assignment.Worker.ID] = struct{}{}
				}
			}
		}

		o.Solutions = append(o.Solutions, nextShiftSolution)
		customResultStatistics := CustomResultStatistics{
			NumberOfWorkers: len(usedWorkers),
		}

		result.Custom = customResultStatistics

		stats.Result = &result
		stats.Run = &run
		o.Statistics = stats
	} else {
		retErr := errors.New("no solution found")
		return schema.Output{}, retErr
	}
	return o, nil
}

func newMIPModel(
	input Input,
	potentialAssignments []Assignment,
	potentialAssignmentsPerWorker map[int][]*Assignment,
	demandCovering map[int][]Assignment,
	opts Options,
) (mip.Model, model.MultiMap[mip.Bool, Assignment]) {
	m := mip.NewModel()
	m.Objective().SetMinimize()

	x := model.NewMultiMap(
		func(...Assignment) mip.Bool {
			return m.NewBool()
		}, potentialAssignments)

	underSupplySlack := model.NewMultiMap(
		func(demand ...RequiredWorker) mip.Float {
			return m.NewFloat(0, float64(demand[0].Count))
		}, input.RequiredWorkers)

	overSupplySlack := model.NewMultiMap(
		func(demand ...RequiredWorker) mip.Float {
			return m.NewFloat(0, math.MaxFloat64)
		}, input.RequiredWorkers)

	for _, demand := range input.RequiredWorkers {
		demandCover := demandCovering[demand.RequiredWorkerID]
		// We need to cover all demands
		coverConstraint := m.NewConstraint(mip.Equal, float64(demand.Count))
		coverConstraint.NewTerm(1.0, underSupplySlack.Get(demand))
		coverConstraint.NewTerm(-1.0, overSupplySlack.Get(demand))
		coverPerWorker := map[int]mip.Constraint{}
		for _, assignment := range demandCover {
			constraint, ok := coverPerWorker[assignment.Worker.ID]
			if !ok {
				constraint = m.NewConstraint(mip.LessThanOrEqual, 1.0)
				coverPerWorker[assignment.Worker.ID] = constraint
			}
			constraint.NewTerm(1.0, x.Get(assignment))
			coverConstraint.NewTerm(1.0, x.Get(assignment))
		}
		m.Objective().NewTerm(opts.OverSupplyPenalty, overSupplySlack.Get(demand))
		m.Objective().NewTerm(opts.UnderSupplyPenalty, underSupplySlack.Get(demand))
	}

	// Two shift of a worker have to be at least x hours apart
	for _, worker := range input.Workers {
		for i, a1 := range potentialAssignmentsPerWorker[worker.ID] {
			// A worker can only work y hours per day
			lessThanXhoursPerDay := m.NewConstraint(mip.LessThanOrEqual, opts.MaxHoursPerDay.Hours())
			lessThanXhoursPerDay.NewTerm(a1.Duration.Hours(), x.Get(*a1))
			atLeastYhoursApart := m.NewConstraint(mip.LessThanOrEqual, 1.0)
			atLeastYhoursApart.NewTerm(1.0, x.Get(*a1))
			lessThanZhoursPerWeek := m.NewConstraint(mip.LessThanOrEqual, float64(opts.MaxHoursPerWeek))
			lessThanZhoursPerWeek.NewTerm(a1.Duration.Hours(), x.Get(*a1))
			for _, a2 := range potentialAssignmentsPerWorker[worker.ID][i+1:] {
				durationApart := a1.DurationApart(*a2)
				if durationApart > 0 {
					// if a1 and a2 do not at least have x hours between them, we
					// forbid them to be assigned at the same time
					if durationApart < opts.HoursBetweenShifts {
						atLeastYhoursApart.NewTerm(1.0, x.Get(*a2))
					}

					if durationApart < 24*time.Hour {
						lessThanXhoursPerDay.NewTerm(a2.Duration.Hours(), x.Get(*a2))
					}

					if durationApart < 7*24*time.Hour {
						lessThanZhoursPerWeek.NewTerm(a2.Duration.Hours(), x.Get(*a2))
					}
				}
			}
		}
	}

	return m, x
}

func potentialAssignments(input Input, opts Options) ([]Assignment, map[int][]*Assignment) {
	potentialAssignments := make([]Assignment, 0)
	potentialAssignmentsPerWorker := map[int][]*Assignment{}
	for _, worker := range input.Workers {
		potentialAssignmentsPerWorker[worker.ID] = make([]*Assignment, 0)
		for _, availability := range worker.Availability {
			for start := availability.Start.Time; start.Before(availability.End.Time); start = start.Add(30 * time.Minute) {
				for end := availability.End.Time; start.Before(end); end = end.Add(-30 * time.Minute) {
					// make sure that end-start is not more than 8h
					duration := end.Sub(start)
					if duration > opts.MaxHoursPerShift {
						continue
					}
					// make sure that end-start is not less than 2h - we are
					// only shrinking the end time, so we can break here
					if duration < opts.MinHoursPerShift {
						break
					}
					assignment := Assignment{
						AssignmentID: len(potentialAssignments),
						Start:        start,
						End:          end,
						Worker:       worker,
						Duration:     duration,
					}
					potentialAssignmentsPerWorker[worker.ID] = append(potentialAssignmentsPerWorker[worker.ID], &assignment)
					potentialAssignments = append(potentialAssignments, assignment)
				}
			}
		}
	}
	return potentialAssignments, potentialAssignmentsPerWorker
}

func demands(input Input, potentialAssignments []Assignment) map[int][]Assignment {
	// initialize demand ids
	for i, demand := range input.RequiredWorkers {
		demand.RequiredWorkerID = i
		input.RequiredWorkers[i] = demand
	}

	demandCovering := map[int][]Assignment{}
	for _, demand := range input.RequiredWorkers {
		demandCovering[demand.RequiredWorkerID] = []Assignment{}
		for i, potentialAssignment := range potentialAssignments {
			if (potentialAssignment.Start.Before(demand.Start.Time) || potentialAssignment.Start.Equal(demand.Start.Time)) &&
				(potentialAssignment.End.After(demand.End.Time) || potentialAssignment.End.Equal(demand.End.Time)) {
				potentialAssignments[i].DemandsCovered = append(potentialAssignments[i].DemandsCovered, demand)
				demandCovering[demand.RequiredWorkerID] = append(demandCovering[demand.RequiredWorkerID], potentialAssignment)
			}
		}
	}
	return demandCovering
}
