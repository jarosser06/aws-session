package main

import (
	"testing"
)

func TestRequestValidation(t *testing.T) {
	type testOptions struct {
		RequiredString string `required:"true"`
		RequiredInt    int    `required:"true"`
		OptionalString string `required:"false"`
		OptionalInt    int
	}

	reqInstance := testOptions{
		RequiredString: "foo",
		RequiredInt:    2,
	}

	if err := Validate(reqInstance); err != nil {
		t.Error(err)
	}

	secondOptInstance := testOptions{
		OptionalInt:    2,
		OptionalString: "foo",
	}

	if err := Validate(secondOptInstance); err == nil {
		t.Error("expected error but got nil")
	}
}

func TestRequestValidation_struct(t *testing.T) {
	type testSubOptions struct {
		RequiredString string `required:"true"`
		OptionalString string `required:"false"`
	}

	type testOptions struct {
		RequiredSubStruct testSubOptions `required:"true"`
	}

	reqInstance := testOptions{
		RequiredSubStruct: testSubOptions{
			RequiredString: "foo",
			OptionalString: "bar",
		},
	}

	if err := Validate(reqInstance); err != nil {
		t.Error(err)
	}
}
