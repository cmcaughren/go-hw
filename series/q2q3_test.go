package main

import "testing"

func TestBasicExample(t *testing.T) {
	result := countDigitInSeries(1, 11, 1, 1, 1)
	if result != 4 {
		t.Errorf("Expected 4, got %d", result)
	}
}

func TestEvenNumbers(t *testing.T) {
	result := countDigitInSeries(1, 20, 1, 2, 2)
	if result != 3 {
		t.Errorf("Expected 3, got %d", result)
	}
}

func TestOddNumbers(t *testing.T) {
	result := countDigitInSeries(1, 20, 1, 1, 3)
	if result != 7 {
		t.Errorf("Expected 7, got %d", result)
	}
}

func TestWithIncrement(t *testing.T) {
	result := countDigitInSeries(1, 20, 3, 1, 1)
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}

func TestDigitZero(t *testing.T) {
	result := countDigitInSeries(5, 15, 1, 0, 1)
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}
}
