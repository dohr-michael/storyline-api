// +build mage

package main

import "github.com/magefile/mage/sh"
import "github.com/magefile/mage/mg"

func init() {

}

func downloadDeps() error {
	/*if err := sh.RunV("go", "get", "-u", "github.com/golang/protobuf/proto"); err != nil {
		return err
	}
	if err := sh.RunV("go", "get", "-u", "github.com/golang/protobuf/protoc-gen-go"); err != nil {
		return err
	}
	if err := sh.RunV("go", "get", "-u", "github.com/jteeuwen/go-bindata/..."); err != nil {
		return err
	}*/
	return sh.RunV("go", "mod", "download")
}

func generate() error {
	return sh.RunV("go", "generate", "./...")
}

func Run() error {
	mg.Deps(downloadDeps, generate)
	return sh.RunV("gin", "--appPort", "8080", "--buildArgs", "main.go", "-i", "run")
}

func Test() error {
	mg.Deps(downloadDeps, generate)
	return sh.RunV("go", "test", "./...")
}

func Build() error {
	mg.Deps(Test)
	return sh.RunV("go", "build", "./...")
}
