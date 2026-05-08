// Creates the UI
package main

import (
	"context"
	"era/lazypg/internal/session"
	"era/lazypg/internal/tui"
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	ctx := context.Background()

	sess, err := session.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close(ctx)

	p := tea.NewProgram(tui.InitialModel(sess))

	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
