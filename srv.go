package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	re       = regexp.MustCompile(`(?m)(\w+)\.service\s+(\w+)\s+(\w+)\s+(\w+)\s+(.*)`)
	commands = map[string]bool{
		"start":   true,
		"restart": true,
		"stop":    true,
	}
)

type ServiceUnit struct {
	Unit        string
	Load        string
	Active      string
	Sub         string
	Description string
}

func ListServices() ([]ServiceUnit, error) {
	cmd := "/usr/bin/systemctl list-units --type service"
	out, err := exec.Command(cmd).Output()
	if err != nil {
		return nil, err
	}
	units := ExtractServiceUnit(string(out))
	return units, nil
}

func ExtractServiceUnit(s string) []ServiceUnit {
	var units []ServiceUnit
	for _, match := range re.FindAllStringSubmatch(s, -1) {
		units = append(units, ServiceUnit{
			Unit:        match[1],
			Load:        match[2],
			Active:      match[3],
			Sub:         match[4],
			Description: match[5],
		})
	}
	return units
}

func ExistsService(unit string) bool {
	items, err := ListServices()
	if err != nil {
		return false
	}
	for _, item := range items {
		if item.Unit == unit {
			return true
		}
	}
	return false
}

func CommandService(command, unit string) error {
	index := strings.IndexByte(command, '\\')
	if index > 0 {
		command = command[0:index]
	}
	if ok := commands[command]; !ok {
		return fmt.Errorf("command[%s] is unsupported", command)
	}

	index = strings.IndexByte(unit, '\\')
	if index > 0 {
		unit = unit[0:index]
	}

	if !ExistsService(unit) {
		return fmt.Errorf("service[%s] is not exists", unit)
	}

	cmd := "systemctl " + command + " " + unit
	_, err := exec.Command(cmd).Output()
	if err != nil {
		return err
	}
	return nil
}
