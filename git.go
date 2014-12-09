package main

import (
	"log"
	"os/exec"
)

type Git struct {
	WorkDir string
}

func (g Git) Commit(message string) {
	err := g.runCommand("commit", "-m", message)
	if err != nil {
		log.Fatal("Error while committing")
	}
}

func (g Git) Add(path string) {
	err := g.runCommand("add", path)
	if err != nil {
		log.Fatal("Error while adding")
	}
}

func (g Git) Push() {
	err := g.runCommand("push")
	if err != nil {
		log.Fatal("Error while pushing")
	}
}

func (g Git) RepoValid() bool {
	return g.runCommand("status") == nil
}

func (g Git) runCommand(arg ...string) error {
	if g.WorkDir == "" {
		log.Fatal("WorkDir not set")
	}
	cmd := exec.Command("git", arg...)
	cmd.Dir = g.WorkDir
	return cmd.Run()
}
