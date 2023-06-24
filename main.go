package main

import (
	"fmt"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
  "flag"
  "runtime"

	"effi-pm/effi"
)

type EnvConfig struct {
  Depth int
  Editor string
  Dirs []string
}

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"

var run_flag = flag.String("run", "", "run a profile")
var build_flag = flag.String("build", "", "build a profile")

func LoadEnvConfig() (EnvConfig, error) {
  dirs_raw := os.Getenv("EFFI_DIRS")

  if dirs_raw == "" {
    wd, err := os.Getwd()
    if err != nil {
      return EnvConfig{}, err
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
    return EnvConfig{}, err
  }

  return EnvConfig{Depth: depth, Editor: editor_raw, Dirs: dirs}, nil
}

func init() {
  if runtime.GOOS == "windows" {
    reset = ""
    red = ""
    green = ""
	}

  flag.StringVar(run_flag, "r", "", "run a profile")
  flag.StringVar(build_flag, "b", "", "build a profile")
}

func printError(err error) {
  fmt.Printf("%vError%v: %v\n", red, reset, err)
  os.Exit(1)
}

func printSuccess(message string) {
  fmt.Printf("%vSuccess%v: %v\n", green, reset, message)
}

func main() {
  flag.Parse()

  run := *run_flag != ""
  build := *build_flag != ""

  if flag.NArg() < 1 {
    printError(errors.New("no project name was specified"))
  }

  project_name := flag.Arg(0)

  project_regex, err := regexp.Compile(project_name)
  if err != nil {
    printError(err)
  }

  env_config, err := LoadEnvConfig()
  if err != nil {
    printError(err)
  }

  project, err, found := effi.SearchDirsRecursive(env_config.Dirs, project_regex, env_config.Depth)
  if err != nil {
    printError(err)
  }

  if !found {
    printError(errors.New("project not found"))
  }

  printSuccess(fmt.Sprintf("found project %v at %v", project.Name, project.Path))

  if !run && !build {
    command := exec.Command(env_config.Editor, ".")
    command.Dir = project.Path
    command.Stdin = os.Stdin
    command.Stdout = os.Stdout
    command.Stderr = os.Stderr
    if err := command.Run(); err != nil {
      printError(err)
    }

    return
  }

  if build {
    if _, ok := project.Profiles[*build_flag]; !ok {
      printError(errors.New("profile does not exits"))
    }
    printSuccess(fmt.Sprintf("building profile %v", *build_flag))
    if err := project.Build(*build_flag); err != nil {
      printError(err)
    }
  }

  if run {
    if _, ok := project.Profiles[*run_flag]; !ok {
      printError(errors.New("profile does not exits"))
    }
    printSuccess(fmt.Sprintf("running profile %v", *run_flag))
    if err := project.Run(*run_flag); err != nil {
      printError(err)
    }
  }
}
