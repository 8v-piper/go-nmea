package pid

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

// systemParameters contains all parameters for an entire system.
type systemParameters map[string]parameters

// ReadJson reads a Json encoded file into this object.
func (s *systemParameters) ReadJson(f string) error {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, s)
	if err != nil {
		return err
	}
	return nil
}

// ReadURLValues reads URL form values into this object.
func (s systemParameters) ReadURLValues(values url.Values) {
	for k, v := range values {
		// The form fields are name <component>_<parameter>.
		parts := strings.SplitN(k, "_", 2)
		if len(parts) == 2 && len(v) == 1 {
			component, param := parts[0], parts[1]
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				continue
			}
			if _, ok := s[component]; !ok {
				s[component] = make(parameters, 0)
			}
			p := &parameter{Name: param, Value: val}
			s[component] = append(s[component], p)
		}
	}
}

// parameters is a collection of parameters for a system component.
type parameters []*parameter

// Get returns the parameter with the specified name.
func (p parameters) Get(n string) (*parameter, bool) {
	for _, param := range p {
		if param.Name == n {
			return param, true
		}
	}
	return &parameter{}, false
}

// SetValue sets the given parameter to the specified value.
func (p parameters) SetValue(n string, v float64) {
	for _, param := range p {
		if param.Name == n {
			param.Value = v
		}
	}
}

// GetValue fetches the value of the specified parameter.
func (p parameters) GetValue(n string) float64 {
	for _, param := range p {
		if param.Name == n {
			return param.Value
		}
	}
	return 0.0
}

// A parameter is a configurable parameter for an IOComponent.
type parameter struct {
	// Name is the shortname of the parameter.
	Name string
	// Title is a more human readable name (for UI display)
	Title string
	// Minimum value of the parameter.
	Minimum float64
	// Maximum value of the parameter.
	Maximum float64
	// The desired step size of the parameter.
	Step float64
	// The default value.
	Default float64
	// The units to display.
	Unit string
	// The value (for setting parameters)
	Value float64
}

// Copy returns a copy of this parameter.
func (p parameter) Copy() *parameter {
	p1 := p
	return &p1
}
