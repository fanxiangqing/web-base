package config

import (
	"fmt"
	"os"
)

func Init() {
	dir, _ := os.Getwd()
	fmt.Println(dir)
}
