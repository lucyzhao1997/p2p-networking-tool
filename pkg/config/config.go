package config

import (
    "fmt"
    "os"
)

func GetConfig() string {
    return os.Getenv("CONFIG")
}