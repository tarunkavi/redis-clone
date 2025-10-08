package core

type Cmd struct {
	Cmd  string
	Args []string
}

type Cmds []*Cmd
