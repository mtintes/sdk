"""
Template for working with Google OR-Tools.
"""

import argparse
import json
import sys
from typing import Any, Dict

from ortools.linear_solver import pywraplp


def main() -> None:
    """Entry point for the template."""

    parser = argparse.ArgumentParser(description="Solve problems with OR-Tools.")
    parser.add_argument(
        "-input",
        default="",
        help="Path to input file. Default is stdin.",
    )
    parser.add_argument(
        "-output",
        default="",
        help="Path to output file. Default is stdout.",
    )
    parser.add_argument(
        "-duration",
        default=30,
        help="Max runtime duration (in seconds). Default is 30.",
        type=int,
    )
    args = parser.parse_args()

    # Read input data, solve the problem and write the solution.
    input_data = read_input(args.input)
    solution = solve(input_data, args.duration)
    write_output(args.output, solution)


def solve(input_data: Dict[str, Any], duration: int) -> Dict[str, Any]:
    """Solves the given problem and returns the solution."""

    # Creates the solver.
    provider = "SCIP"
    solver = pywraplp.Solver.CreateSolver(provider)
    solver.SetTimeLimit(duration * 1000)

    # Initializes the linear sums.
    weights = 0.0
    values = 0.0

    # Creates the decision variables and adds them to the linear sums.
    items = []
    for item in input_data["items"]:
        item_variable = solver.IntVar(0, 1, item["item_id"])
        items.append({"item": item, "variable": item_variable})
        weights += item_variable * item["weight"]
        values += item_variable * item["value"]

    # This constraint ensures the weight capacity of the knapsack will not be
    # exceeded.
    solver.Add(weights <= input_data["weight_capacity"])

    # Sets the objective function: maximize the value of the chosen items.
    solver.Maximize(values)

    # Solves the problem.
    status = solver.Solve()

    # Determines which items were chosen.
    chosen_items = []
    for item in items:
        if item["variable"].solution_value() > 0.9:
            chosen_items.append(item["item"])

    # Creates the statistics.
    statistics = {
        "result": {
            "custom": {
                "constraints": solver.NumConstraints(),
                "provider": provider,
                "status": status,
                "variables": solver.NumVariables(),
            },
            "duration": solver.WallTime() / 1000,
            "value": solver.Objective().Value(),
        },
        "run": {
            "duration": solver.WallTime() / 1000,
        },
        "schema": "v1",
    }

    return {
        "solutions": [{"items": chosen_items}],
        "statistics": statistics,
    }


def read_input(input_path) -> Dict[str, Any]:
    """Reads the input from stdin or a given input file."""

    input_file = {}
    if input_path:
        with open(input_path, "r", encoding="utf-8") as file:
            input_file = json.load(file)
    else:
        input_file = json.load(sys.stdin)

    return input_file


def write_output(output_path, output) -> None:
    """Writes the output to stdout or a given output file."""

    content = json.dumps(output, indent=2)
    if output_path:
        with open(output_path, "w", encoding="utf-8") as file:
            file.write(content + "\n")
    else:
        print(content)


if __name__ == "__main__":
    main()
