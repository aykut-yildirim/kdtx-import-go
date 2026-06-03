package main

import (
	"myapp/internal/api"
	_ "myapp/internal/portals/amazon"
	_ "myapp/internal/portals/etsy"
	_ "myapp/internal/portals/stripe"
)

func main() {
	r := api.SetupRouter()
	r.Run(":8000")
}