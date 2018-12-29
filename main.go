package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var GradeValues = map[string]float64{
	"A":  4.00,
	"A-": 3.67,
	"B+": 3.33,
	"B":  3.00,
	"B-": 2.67,
	"C+": 2.33,
	"C":  2.00,
	"C-": 1.67,
	"D":  1.00,
	"F":  0.00,
}

var usage = `gpa - Calculate the gpa given a series of (grade, weight) pairs.

Usage:
	gpa [(<grade> <credits>)...]
	gpa [<file>]
`

type Grade struct {
	grade   string
	credits int
}

func (g *Grade) QualityPoints() (float64, error) {
	if value, ok := GradeValues[g.grade]; ok {
		return value * float64(g.credits), nil
	} else {
		return -1.00, errors.New("No such grade")
	}
}

func argvGrades(grades chan<- Grade) {
	for i := 1; i < len(os.Args); i += 2 {
		grade := os.Args[i]
		credits, err := strconv.Atoi(os.Args[i+1])
		if err != nil {
			panic(err)
		}

		grades <- Grade{grade, credits}
	}

	close(grades)
}

func fileGrades(file *os.File, grades chan<- Grade) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line_segments := strings.Fields(scanner.Text())
		grade := line_segments[0]

		credits, err := strconv.Atoi(line_segments[1])

		if err != nil {
			panic(err)
		}

		grades <- Grade{grade, credits}
	}

	close(grades)
}

func totalStats(grades <-chan Grade, gpa chan<- float64) {
	total_points := 0.0
	total_credits := 0

	for grade := range grades {
		total_credits += grade.credits
		quality_points, err := grade.QualityPoints()
		if err != nil {
			panic(err)
		}
		total_points += quality_points
	}

	gpa <- total_points / float64(total_credits)
}

func main() {
	grades := make(chan Grade)
	gpa := make(chan float64, 1)

	switch len(os.Args) {
	case 1:
		go fileGrades(os.Stdin, grades)
	case 2:
		arg := os.Args[1]
		switch arg {
		case "-h":
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(0)
		case "--help":
			fmt.Fprintln(os.Stderr, usage)
			os.Exit(0)
		case "-":
			go fileGrades(os.Stdin, grades)
		default:
			file, err := os.Open(os.Args[1])
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "Error: file does not exist: %s\n", os.Args[1])
				}
				os.Exit(1)
			}
			go fileGrades(file, grades)
		}
	default:
		go argvGrades(grades)
	}

	go totalStats(grades, gpa)
	fmt.Printf("%.3f\n", <-gpa)
	close(gpa)
}
