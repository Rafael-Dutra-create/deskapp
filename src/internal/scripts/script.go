package main

type IScript interface {
	Name() string
    Description() string
	Execute(args []string) error
}

type ScriptBase struct{}

func (s *ScriptBase) Execute(args []string) error {
    return nil
}