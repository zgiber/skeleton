package main

import (
	"context"
	"fmt"
)

type service struct{}

func (s *service) Greeting(_ context.Context, name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}
