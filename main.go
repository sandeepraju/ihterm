package main

import (
	"fmt"

	"github.com/sandeepraju/ihterm/pkg"
)

func main() {
	iht := pkg.NewIHTerm()
	fmt.Println(iht.BitBarString())
}
