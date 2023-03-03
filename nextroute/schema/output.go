package schema

import (
	"time"
)

type JsonSolution struct {
	Epoch       time.Time              `json:"epoch"`
	Unplanned   []JsonModelStop        `json:"unplanned"`
	Vehicles    []JsonVehicle          `json:"vehicles"`
	Progression []JsonObjectiveElapsed `json:"progression"`
	Objective   JsonObjective          `json:"objective"`
}

type JsonObjectiveElapsed struct {
	ElapsedSeconds float64 `json:"elapsed_seconds"`
	Value          float64 `json:"value"`
}

type JsonObjective struct {
	Name       string          `json:"name"`
	Objectives []JsonObjective `json:"objectives"`
	Value      float64         `json:"value"`
}

type JsonVehicle struct {
	Start time.Time  `json:"start"`
	End   time.Time  `json:"end"`
	Name  string     `json:"name"`
	Stops []JsonStop `json:"stops"`
	Index int        `json:"index"`
}

type JsonModelStop struct {
	Name     string       `json:"name"`
	Location JsonLocation `json:"location"`
	Index    int          `json:"index"`
}

type JsonStop struct {
	Arrival           time.Time             `json:"arrival"`
	EarliestStart     time.Time             `json:"earliest_start"`
	Start             time.Time             `json:"start"`
	End               time.Time             `json:"end"`
	ExpressionValues  []JsonExpressionValue `json:"expression_values"`
	ConstraintValues  []JsonConstraintValue `json:"constraint_values"`
	ConstraintReports []map[string]any      `json:"constraint_reports"`
	JsonModelStop
	Position       int           `json:"position"`
	TravelDuration time.Duration `json:"travel_duration"`
}
type JsonExpressionValue struct {
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Cumulative float64 `json:"cumulative"`
}

type JsonConstraintValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type JsonLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
