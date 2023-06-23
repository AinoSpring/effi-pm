package effi

import (
	"effi-pm/parser"
	"os"
	"path/filepath"
)

type Project struct {
  Name string;
  Profiles map[string]Profile;
  Path string;
}

type Profile struct {
  RunCommand string;
  BuildCommand string;
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

