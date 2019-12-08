package gormchecker_test

import (
	"testing"

	"github.com/gosagawa/gormchecker"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, gormchecker.Analyzer, "a")
}
