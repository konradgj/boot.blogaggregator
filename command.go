package main

import (
	"fmt"
	"sync"
)

type command struct {
	name string
	args []string
}

type commands struct {
	mu             sync.RWMutex
	registeredCmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	c.mu.RLock()
	f, ok := c.registeredCmds[cmd.name]
	c.mu.RUnlock()

	if !ok {
		return fmt.Errorf("unkown command: %s", cmd.name)
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.registeredCmds[name]; exists {
		return fmt.Errorf("command %q already registered", name)
	}

	c.registeredCmds[name] = f
	return nil
}
