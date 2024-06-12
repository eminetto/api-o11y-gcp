//go:build mage
// +build mage

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	log.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "bin/api-o11y-gcp", "cmd/api/main.go")
	return cmd.Run()
}

// Build for Linux
func BuildLinux() error {
	mg.Deps(InstallDeps)
	log.Println("Generating Linux binary...")
	os.Setenv("CGO_ENABLED", "0")
	os.Setenv("GOOS", "linux")
	cmd := exec.Command("go", "build", "-a", "-installsuffix", "cgo", "-tags", `"netgo"`, "-installsuffix", "netgo", "-o", "bin/api-o11y-gcp", "cmd/api/main.go")
	return cmd.Run()
}

// Build a docker image
func BuildDocker() error {
	log.Println("Building...")
	cmd := exec.Command("docker", "build", "-t", "api-o11y-gcp", "-f", "Dockerfile", ".")
	return cmd.Run()
}

// Generate mocks
func GenerateMocks() error {
	log.Println("Installing mockery...")
	cmd := exec.Command("go", "install", "github.com/vektra/mockery/v2@v2.43.1")
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Println("Generating user mocks...")
	cmd = exec.Command("mockery", "--output", "user/mocks", "--dir", "user", "--all")
	err = cmd.Run()
	if err != nil {
		return err
	}
	log.Println("Generating telemetry mocks...")
	cmd = exec.Command("mockery", "--output", "internal/telemetry/mocks", "--dir", "internal/telemetry", "--all")
	return cmd.Run()
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	log.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "tidy")
	return cmd.Run()
}

// Run tests
func Test() error {
	mg.Deps(GenerateMocks)
	cmd := exec.Command("go", "test", "./...")
	return cmd.Run()
}

// Run the docker image
func RunDocker() error {
	mg.Deps(BuildDocker)
	cmd := exec.Command("docker", "run", "-p", "8080:8080", "api-o11y-gcp")
	return cmd.Run()
}

// Clean up after yourself
func Clean() error {
	log.Println("Cleaning...")
	err := removeGlob("user/mocks/*")
	if err != nil {
		return err
	}
	err = removeGlob("internal/telemetry/mocks/*")
	if err != nil {
		return err
	}
	return os.RemoveAll("bin/api-o11y-gcp")
}

func removeGlob(path string) (err error) {
	contents, err := filepath.Glob(path)
	if err != nil {
		return
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return
		}
	}
	return
}
