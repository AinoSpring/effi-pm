package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func ParseIni(data string) map[string]map[string]string {
  parsed_ini := make(map[string]map[string]string)

  scanner := bufio.NewScanner(strings.NewReader(data))

  section_regex := regexp.MustCompile("^\\[(?:(.+))\\]$")
  value_regex := regexp.MustCompile("^(?:([^ \t].+?))=(?:(.+?))$")

  var section string = ""

  for scanner.Scan() {
    line := scanner.Text()

    section_match := section_regex.FindStringSubmatch(line)
    value_match := value_regex.FindStringSubmatch(line)

    if section_match != nil {
      section = section_match[1]
      parsed_ini[section] = make(map[string]string)
    }
    if value_match != nil {
      parsed_ini[section][value_match[1]] = value_match[2]
    }
  }

  return parsed_ini
}

func ParseIniFile(path string) (map[string]map[string]string, error) {
  contents, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }
  return ParseIni(string(contents)), nil
}
