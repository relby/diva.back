package model_test

import (
	"testing"

	"github.com/relby/diva.back/internal/model"
)

func TestValidCustomerID(t *testing.T) {
	_, err := model.NewCustomerID(10)

	if err != nil {
		t.Errorf("ERROR: %v", err)
	}
}

func TestInvalidCustomerID(t *testing.T) {
	_, err := model.NewCustomerID(-10)

	if err == nil {
		t.Errorf("ERROR: %v", err)
	}
}
