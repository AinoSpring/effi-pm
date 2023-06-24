package effi

import (
	"os"
  "os/exec"
	"path/filepath"
	"effi-pm/parser"
)

type Project struct {
  Name string
  Profiles map[string]Profile
  Path string
}

type Profile struct {
  RunCommand string
  BuildCommand string
}

func ParseProject(data string) Project {
  project := Project{}
  
  parsed_ini := parser.ParseIni(data)

  profiles := make(map[string]Profile)

  for section, values := range parsed_ini {
    if section == "project" {
      project.Name = values["name"]
      continue
    }

    profiles[section] = Profile{
      RunCommand: values["run"],
      BuildCommand: values["build"],
    }
  }

  project.Profiles = profiles

  return project
}

func ParseProjectFile(path string) (Project, error) {
  contents, err := os.ReadFile(path)
  if err != nil {
    return Project{}, err
  }

  project := ParseProject(string(contents))
  project.Path = filepath.Dir(path)

  return project, nil
}

func (project Project) Run(profile string) error {
  command := exec.Command("sh", "-c", project.Profiles[profile].RunCommand)
  command.Dir = project.Path
  command.Stdin = os.Stdin
  command.Stdout = os.Stdout
  command.Stderr = os.Stderr
  return command.Run()
}

func (project Project) Build(profile string) error {
  command := exec.Command("sh", "-c", project.Profiles[profile].BuildCommand)
  command.Dir = project.Path
  command.Stdin = os.Stdin
  command.Stdout = os.Stdout
  command.Stderr = os.Stderr
  return command.Run()
}

