package main

import (
	"github.com/softleader/memory-calculator/calc"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMemoryCalculation_WithLoadedClassCount(t *testing.T) {
	// 避免在不同環境中, 計算出來的結果超過那個環境實際的記憶體, 進而造成測試錯誤, 所以這邊我們假裝有足夠大的記憶體
	file, err := createTempFile("10g")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	calculator := calc.NewCalculator()
	calculator.MemoryLimitPath.V1 = file.Name()
	*calculator.ThreadCount = 10
	*calculator.LoadedClassCount = 100

	jto, err := calculator.Execute()
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}
	options := jto.String()
	for _, o := range calc.ContributeOptions {
		options = strings.ReplaceAll(options, o, "")
	}
	options = strings.ReplaceAll(options, calc.OptionsSeparator, "")

	if options == "" {
		t.Errorf("Execute returned an empty string")
	}
}

func TestMemoryCalculation_WithoutLoadedClassCount(t *testing.T) {
	// 避免在不同環境中, 計算出來的結果超過那個環境實際的記憶體, 進而造成測試錯誤, 所以這邊我們假裝有足夠大的記憶體
	file, err := createTempFile("10g")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	calculator := calc.NewCalculator()
	calculator.MemoryLimitPath.V1 = file.Name()
	*calculator.ThreadCount = 10
	*calculator.AppPath = "."
	*calculator.JVMClassCount = 100
	*calculator.JVMClassAdj = "10%"

	jto, err := calculator.Execute()
	if err != nil {
		t.Fatalf("Execute returned an error: %v", err)
	}
	options := jto.String()
	for _, o := range calc.ContributeOptions {
		options = strings.ReplaceAll(options, o, "")
	}
	options = strings.ReplaceAll(options, calc.OptionsSeparator, "")

	if options == "" {
		t.Errorf("Execute returned an empty string")
	}
}

func createTempFile(content string) (*os.File, error) {
	file, err := os.CreateTemp("", "testing-temp-file-")
	if err != nil {
		return nil, err
	}

	_, err = io.WriteString(file, content)
	if err != nil {
		return nil, err
	}
	return file, nil
}
