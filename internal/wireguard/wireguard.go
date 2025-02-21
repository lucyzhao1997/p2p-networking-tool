package wireguard

import (
    "fmt"
    "os/exec"
)

func StartWireGuard() error {
    cmd := exec.Command("sudo", "wg-quick", "up", "wg0")
    if err := cmd.Run(); err != nil {
        return err
    }
    fmt.Println("WireGuard started")
    return nil
}