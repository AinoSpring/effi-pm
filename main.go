package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"effi-pm/effi"

	"github.com/urfave/cli/v2"
)

func main() {
  app := cli.App{
    Name: "Effi project manager",
    Action: func(ctx *cli.Context) error {
      project_name := ctx.Args().First()
      if project_name == "" {
        return errors.New("no project name was specified")
      }

      project_regex, err := regexp.Compile(project_name)
      if err != nil {
        return err
      }

      dirs_raw := os.Getenv("EFFI_DIRS")

      if dirs_raw == "" {
        wd, err := os.Getwd()
        if err != nil {
          return err
        }

        dirs_raw = wd
      }

      depth_raw := os.Getenv("EFFI_DEPTH")
      
      if depth_raw == "" {
        depth_raw = "-1"
      }

      editor_raw := os.Getenv("EFFI_EDITOR")

      if editor_raw == "" {
        editor_raw = os.Getenv("EDITOR")
      }

      dirs := strings.Split(dirs_raw, ":")
      depth, err := strconv.Atoi(depth_raw)

      if err != nil {
        return err
      }

      project, err, found := effi.SearchDirsRecursive(dirs, project_regex, depth)

      if err != nil {
        return err
      }
      if !found {
        return errors.New("project not found")
      }

      fmt.Printf("opening %s at %v\n", project.Name, project.Path)

      command := exec.Command(editor_raw, project.Path)
      command.Stdin = os.Stdin
      command.Stdout = os.Stdout
      command.Stderr = os.Stderr
      if err := command.Run(); err != nil {
        return err
      }

      return nil
    },
  }

  if err := app.Run(os.Args); err != nil {
    log.Fatal(err);
  }
}
