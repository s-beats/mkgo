package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/dave/jennifer/jen"
)

func main() {
	targetDir, err := getTargetDir()
	if err != nil {
		fmt.Printf("Failed to get target dir: %s\n", err)
		os.Exit(1)
	}

	if err := makeGoMod(targetDir); err != nil {
		fmt.Printf("Failed to make go.mod:%s\n", err)
		os.Exit(1)
	}

	if err := makeMain(targetDir); err != nil {
		fmt.Printf("Failed to make main.go:%s\n", err)
		os.Exit(1)
	}
}

func makeDir(dir string) error {
	if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(dir, 0700); err != nil {
			return err
		}
	}
	return nil
}

func getTargetDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	targetDir := wd

	if len(os.Args[1:]) > 0 {
		targetDir = path.Join(wd, os.Args[1])
		if err := makeDir(targetDir); err != nil {
			return "", err
		}
	}

	return targetDir, nil
}

func makeMain(dir string) error {
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block()
	return f.Save(filepath.Join(dir, "main.go"))
}

func makeGoMod(dir string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(wd)

	if os.Chdir(dir); err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "init", "example.com")
	cmd.Dir = dir
	return cmd.Run()
}
