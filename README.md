## go-interruptible-service

A library to create an interruptible service that has a graceful shutdown function to run after the signals

## Example

```
package main

import (
	"fmt"
	"sync"

	"github.com/eekwong/go-interruptible-service"
	"github.com/gofiber/fiber/v2"
)

type App struct {
	interruptible.Service
}

func (app *App) Run() (interruptible.Stop, error) {
    fmt.Println("service starts")

	var wg sync.WaitGroup

	fiberApp := fiber.New()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fiberApp.Listen(":8080"); err != nil {
			fmt.Errorf("err in fiber: %v", err)
		}
	}()

	return func() error {
		fiberApp.Shutdown()

		wg.Wait()

		fmt.Println("service stopped")
		return nil
	}, nil
}

func main() {
	app := &App{}
	interruptible.Run(app)
}
```
