package main

import (
	"context"

	"github.com/jeronimobarea/transaction_parser/internal/platform"
)

func main() {
	platform.Run(context.Background())
}
