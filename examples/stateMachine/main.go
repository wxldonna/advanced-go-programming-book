package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// State the FSM state for turnstile
type State uint32

const (
	// Locked locked state
	Locked State = iota
	// Unlocked unlocked state
	Unlocked
)

const (
	// CmdCoin command coin
	CmdCoin = "coin"
	// CmdPush command push
	CmdPush = "push"
)

// Turnstile the finite state machine
type Turnstile struct {
	State State
}

// ExecuteCmd execute command
func (p *Turnstile) ExecuteCmd(cmd string) {
	// get function from transition table
	tupple := CmdStateTupple{strings.TrimSpace(cmd), p.State}
	if f := StateTransitionTable[tupple]; f == nil {
		fmt.Println("unknown command, try again please")
	} else {
		f(&p.State)
	}
}

func main() {
	machine := &Turnstile{State: Locked}
	prompt(machine.State)
	reader := bufio.NewReader(os.Stdin)

	for {
		// read command from stdin
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		machine.ExecuteCmd(cmd)
	}
}

// CmdStateTupple tupple for state-command combination
type CmdStateTupple struct {
	Cmd   string
	State State
}

// TransitionFunc transition function
type TransitionFunc func(state *State)

// StateTransitionTable trsition table
var StateTransitionTable = map[CmdStateTupple]TransitionFunc{
	{CmdCoin, Locked}: func(state *State) {
		fmt.Println("unlocked, ready for pass through")
		*state = Unlocked
	},
	{CmdPush, Locked}: func(state *State) {
		fmt.Println("not allowed, unlock first")
	},
	{CmdCoin, Unlocked}: func(state *State) {
		fmt.Println("well, don't waste your coin")
	},
	{CmdPush, Unlocked}: func(state *State) {
		fmt.Println("pass through, shift back to locked")
		*state = Locked
	},
}

func prompt(s State) {
	m := map[State]string{
		Locked:   "Locked",
		Unlocked: "Unlocked",
	}
	fmt.Printf("current state is [%s], please input command [coin|push]\n", m[s])
}
