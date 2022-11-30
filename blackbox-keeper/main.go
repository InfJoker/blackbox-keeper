package main

import (
	"blackbox-keeper/configuration"
	"fmt"
)

func main() {
    cgf, err := configuration.NewConfiguration("conf.yml")
    if err != nil {
        fmt.Printf("%v", err)
    }

    fmt.Printf("%v", cgf)
}
