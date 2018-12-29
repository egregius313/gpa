# gpa
Command line application for calculating a GPA

## Usage

`gpa` provides a simple interface

```
gpa - Calculate the gpa given a series of (grade, weight) pairs.

Usage:
	gpa [(<grade> <credits>)...]
	gpa [<file>]

```

Either invoke it with a letter grade/credit amount pairs as command line parameters, or pass it a file from which to read the pairs.

If neither is specified, it will read from stdin like it would from a file.
