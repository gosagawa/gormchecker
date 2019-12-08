package main

import (
	"github.com/gosagawa/gormchecker"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(gormchecker.Analyzer) }
