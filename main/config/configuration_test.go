package config

import (
	"flag"
	"reflect"
	"testing"
)

func TestFieldToEnv(t *testing.T) {

	tests := []struct {
		title    string
		in       string
		expected string
	}{
		{
			title:    "empty string",
			in:       "",
			expected: "PROXEUS_",
		},
		{
			title:    "one lowercase char",
			in:       "a",
			expected: "PROXEUS_A",
		},
		{
			title:    "one uppercase char",
			in:       "A",
			expected: "PROXEUS_A",
		},
		{
			title:    "lowercase chars",
			in:       "abc",
			expected: "PROXEUS_ABC",
		},
		{
			title:    "uppercase chars",
			in:       "ABC",
			expected: "PROXEUS_ABC",
		},
		{
			title:    "starting underscore",
			in:       "_abc",
			expected: "PROXEUS_ABC",
		},
		{
			title:    "ending underscore",
			in:       "abc_",
			expected: "PROXEUS_ABC_",
		},
		{
			title:    "underscores",
			in:       "_abc_",
			expected: "PROXEUS_ABC_",
		},
		{
			title:    "underscores",
			in:       "_a_b_c_",
			expected: "PROXEUS_A_B_C_",
		},
		{
			title:    "camel case",
			in:       "aAbBc",
			expected: "PROXEUS_A_AB_BC",
		},
		{
			title:    "underscores",
			in:       "aA-bB-c",
			expected: "PROXEUS_A_A_B_B_C",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if result := fieldToEnv(test.in); result != test.expected {
				t.Errorf("Exptected %s --> %s but go %s", test.in, test.expected, result)
			}

		})
	}
}

func TestFlagFromStruct(t *testing.T) {
	type inner struct {
		InnerValue string `default:"in"`
	}
	type outer struct {
		OuterValue string `default:"out"`
		Inner      inner
	}
	tests := []struct {
		title    string
		in       outer
		params   []string
		env      []string
		expected outer
	}{
		{
			title:    "no params",
			in:       outer{},
			expected: outer{OuterValue: "out", Inner: inner{InnerValue: "in"}},
		},
		{
			title: "env",
			in:    outer{},
			env: []string{
				"PROXEUS_OUTER_VALUE=hello",
				"PROXEUS_INNER_VALUE=world",
			},
			expected: outer{OuterValue: "hello", Inner: inner{InnerValue: "world"}},
		},
		{
			title:    "params",
			in:       outer{OuterValue: "foo", Inner: inner{InnerValue: "bar"}},
			params:   []string{"-OuterValue=abc", "-InnerValue=xyz"},
			expected: outer{OuterValue: "abc", Inner: inner{InnerValue: "xyz"}},
		},
		{
			title: "env and params",
			in:    outer{OuterValue: "foo", Inner: inner{InnerValue: "bar"}},
			env: []string{
				"PROXEUS_OUTER_VALUE=hello",
				"PROXEUS_INNER_VALUE=world",
			},
			params:   []string{"-OuterValue=abc", "-InnerValue=xyz"},
			expected: outer{OuterValue: "abc", Inner: inner{InnerValue: "xyz"}},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			flagFromStruct(fs, test.env, &test.in)
			fs.Parse(test.params)
			fs.VisitAll(func(f *flag.Flag) {
				t.Logf("%#v\n", f)
			})
			if !reflect.DeepEqual(test.in, test.expected) {
				t.Errorf("Expected %#v got %#v", test.expected, test.in)
			}
		})
	}

}
