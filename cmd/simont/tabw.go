package main

import (
	"os"
	"text/tabwriter"
)

var tabw = tabwriter.NewWriter(os.Stdout, 5, 3, 1, ' ', 0)
