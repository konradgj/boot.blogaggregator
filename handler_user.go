package main

import "fmt"

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected argument: login <name>")
	}

	err := state.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.args[0])

	return nil
}
