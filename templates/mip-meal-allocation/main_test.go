package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/store"
)

type output struct {
	Status  string  `json:"status,omitempty"`
	Runtime string  `json:"runtime,omitempty"`
	Binkies float64 `json:"binkies,omitempty"`
	Meals   []struct {
		Name     string `json:"name,omitempty"`
		Quantity int    `json:"quantity,omitempty"`
	} `json:"meals,omitempty"`
}

func TestTemplate(t *testing.T) {
	// Read the input from the file.
	input := input{}
	b, err := os.ReadFile("input.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}

	// Declare the options.
	opt := store.DefaultOptions()
	opt.Limits.Duration = 5 * time.Second
	opt.Diagram.Expansion.Limit = 1
	opt.Limits.Solutions = 1

	// Declare the solver.
	solver, err := solver(input, opt)
	if err != nil {
		t.Fatal(err)
	}

	// Get the solution.
	last := solver.Last(context.Background())
	b, err = json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	gotOutput := output{}
	if err := json.Unmarshal(b, &gotOutput); err != nil {
		t.Fatal(err)
	}

	// Get the expected solution.
	wantOutput := output{}
	b, err = os.ReadFile("testdata/output.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &wantOutput); err != nil {
		t.Fatal(err)
	}

	got := gotOutput.Binkies
	want := wantOutput.Binkies

	// Compare against expected.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
