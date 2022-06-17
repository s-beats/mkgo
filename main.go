package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dave/jennifer/jen"
)

func main() {
	targetDir := os.Args[1]

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path := filepath.Join(dir, targetDir)

	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(path, 0700); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := makeGoMod(path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := makeMain(path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func makeDir(dir string) error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path := filepath.Join(wd, dir)

	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(path, 0700); err != nil {
			os.Exit(1)
		}
	}
	return nil
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
