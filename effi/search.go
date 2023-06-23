package effi

import (
	"os"
	"path/filepath"
	"regexp"
)

func SearchDir(path string, regex *regexp.Regexp) (Project, error, bool) {
  contents, err := os.ReadDir(path)
  if err != nil {
    return Project{}, err, false
  }

  for _, file := range contents {
    if file.IsDir() {
      continue
    }
  
    if file.Name() != ".effi" {
      continue
    }

    project, err := ParseProjectFile(filepath.Join(path, file.Name()))
    if err != nil {
      return Project{}, err, false
    }

    if regex.MatchString(project.Name) {
      return project, nil, true
    }
  }
  return Project{}, nil, false
}

func SearchDirsRecursive(dirs []string, regex *regexp.Regexp, depth int) (Project, error, bool) {
  if depth == 0 {
    return Project{}, nil, false
  }

  contents := make(map[string][]os.DirEntry)

  for _, dir := range dirs {
    dir_contents, err := os.ReadDir(dir)
    if err != nil {
      return Project{}, err, false
    }

    contents[dir] = dir_contents
  }

  for ; len(contents) * depth != 0; depth -= 1 {
    new_contents := make(map[string][]os.DirEntry)

    for path, path_contents := range contents {
      project, err, found := SearchDir(path, regex)
      if found {
        return project, nil, true
      }

      if err != nil {
        return Project{}, err, false
      }

      for _, dir := range path_contents {
        if !dir.IsDir() {
          continue
        }

        dir_path := filepath.Join(path, dir.Name())

        dir_contents, err := os.ReadDir(dir_path)

        if err != nil {
          return Project{}, err, false
        }

        new_contents[dir_path] = dir_contents
      }
    }

    contents = new_contents

    if depth < 0 {
      depth = -1
    }
  }

  return Project{}, nil, false
}
