package main

import "encoding/json"

type Module struct {
	Path    string
	Version string
}

type GoMod struct {
	Module  Module
	Require []Require
	Exclude []Module
	Replace []Replace
}

type Require struct {
	Path     string
	Version  string
	Indirect bool
}

type Replace struct {
	Old Module
	New Module
}

func unmarshalGoMod(data []byte) (*GoMod, error) {
	r := new(GoMod)
	err := json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
