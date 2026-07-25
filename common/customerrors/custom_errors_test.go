package customerrors

import (
	"errors"
	"testing"
)

func TestErrModuleNotRegistered(t *testing.T) {
	if ErrModuleNotRegistered == nil {
		t.Fatal("ErrModuleNotRegistered should not be nil")
	}
	if ErrModuleNotRegistered.Error() != "can't register module" {
		t.Errorf("unexpected message: %v", ErrModuleNotRegistered.Error())
	}
}

func TestErrModuleNameAlreadyExists(t *testing.T) {
	if ErrModuleNameAlreadyExists == nil {
		t.Fatal("ErrModuleNameAlreadyExists should not be nil")
	}
	if ErrModuleNameAlreadyExists.Error() != "module name already exists" {
		t.Errorf("unexpected message: %v", ErrModuleNameAlreadyExists.Error())
	}
}

func TestErrSeedingModule(t *testing.T) {
	if ErrSeedingModule == nil {
		t.Fatal("ErrSeedingModule should not be nil")
	}
}

func TestErrTypeAssersionFailed(t *testing.T) {
	if ErrTypeAssersionFailed == nil {
		t.Fatal("ErrTypeAssersionFailed should not be nil")
	}
}

func TestErrInvalidObjectId(t *testing.T) {
	if ErrInvalidObjectId == nil {
		t.Fatal("ErrInvalidObjectId should not be nil")
	}
	if ErrInvalidObjectId.Error() != "invalid object id" {
		t.Errorf("unexpected message: %v", ErrInvalidObjectId.Error())
	}
}

func TestErrInsufficientReady(t *testing.T) {
	if ErrInsufficientReady == nil {
		t.Fatal("ErrInsufficientReady should not be nil")
	}
}

func TestErrorsAreUnique(t *testing.T) {
	errs := []error{
		ErrModuleNotRegistered,
		ErrModuleNameAlreadyExists,
		ErrSeedingModule,
		ErrTypeAssersionFailed,
		ErrInvalidObjectId,
		ErrInsufficientReady,
	}

	seen := make(map[string]bool)
	for _, err := range errs {
		msg := err.Error()
		if seen[msg] {
			t.Errorf("duplicate error message: %q", msg)
		}
		seen[msg] = true
	}
}

func TestErrorsAreComparable(t *testing.T) {
	if !errors.Is(ErrModuleNotRegistered, ErrModuleNotRegistered) {
		t.Error("ErrModuleNotRegistered should be comparable to itself")
	}
	if errors.Is(ErrModuleNotRegistered, ErrModuleNameAlreadyExists) {
		t.Error("ErrModuleNotRegistered should not equal ErrModuleNameAlreadyExists")
	}
}
