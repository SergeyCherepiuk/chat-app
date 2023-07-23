package utils

import (
	"errors"
	"testing"
)

func TestTwoNilErrors(t *testing.T) {
	var err1 error = nil
	var err2 error = nil

	actual := AreErrorsEqual(err1, err2)
	expected := true

	if actual != expected {
		t.Errorf("Expected: %v, got %v\n", expected, actual)
	}
}

func TestTwoDifferentErrors(t *testing.T) {
	err1 := errors.New("first name is empty")
	err2 := errors.New("last name is empty")

	actual := AreErrorsEqual(err1, err2)
	expected := false

	if actual != expected {
		t.Errorf("Expected: %v, got %v\n", expected, actual)
	}
}

func TestTwoDifferentlySizedErrors(t *testing.T) {
	err1 := errors.New("first name is empty")
	err2 := errors.Join(
		errors.New("first name is empty"),
		errors.New("last name is empty"),
	)

	actual := AreErrorsEqual(err1, err2)
	expected := false

	if actual != expected {
		t.Errorf("Expected: %v, got %v\n", expected, actual)
	}
}
