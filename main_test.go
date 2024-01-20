package main

import (
	"github.com/softleader/memory-calculator/calc"
	"strings"
	"testing"
)

func TestMemoryCalculation_WithLoadedClassCount(t *testing.T) {
	calculator := calc.NewCalculator()
	*calculator.LoadedClassCount = 1000

	jto, err := calculator.Execute()
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}
	options := jto.String()
	for _, o := range calc.ContributeOptions {
		options = strings.ReplaceAll(options, o, "")
	}
	if options == " " {
		t.Errorf("Execute returned an empty string")
	}
}

func TestMemoryCalculation_WithoutLoadedClassCount(t *testing.T) {
	calculator := calc.NewCalculator()
	*calculator.AppPath = "."
	*calculator.JVMClassCount = 1000
	*calculator.JVMClassAdj = "150%"

	jto, err := calculator.Execute()
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}
	options := jto.String()
	for _, o := range calc.ContributeOptions {
		options = strings.ReplaceAll(options, o, "")
	}
	if options == " " {
		t.Errorf("Execute returned an empty string")
	}
}
