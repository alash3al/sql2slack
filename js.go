package main

import (
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

type JSVM struct {
	src    string
	vm     *otto.Otto
	script *otto.Script
}

func NewJSVM(name, src string) (*JSVM, error) {
	vm := otto.New()
	script, err := vm.Compile(name, src)
	if err != nil {
		return nil, err
	}
	return &JSVM{
		src:    src,
		vm:     vm,
		script: script,
	}, nil
}

func (vm *JSVM) Exec(ctx map[string]interface{}) error {
	for k, v := range ctx {
		if err := vm.vm.Set(k, v); err != nil {
			return err
		}
	}
	_, err := vm.vm.Run(vm.script)
	return err
}
