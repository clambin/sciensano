package apihandler

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"strings"
)

func GetTestTargets() []string {
	return []string{
		"tests-total",
		"tests-positive",
		"tests-rate",
	}
}

func GetVaccinationsTargets() []string {
	return []string{
		"vaccinations-first",
		"vaccinations-second",
	}
}

// Targets

func GetVaccinationsByAgeTargets() (targets []string) {
	for _, ageGroup := range sciensano.AgeGroups {
		targets = append(targets, "vaccinations-"+ageGroup+"-first")
		targets = append(targets, "vaccinations-"+ageGroup+"-second")
	}
	return
}

func GetVaccinationsByRegionTargets() (targets []string) {
	for _, region := range sciensano.Regions {
		targets = append(targets, "vaccinations-"+region+"-first")
		targets = append(targets, "vaccinations-"+region+"-second")
	}
	return
}

// Generic functions

func GetModeFromTarget(target string) (mode string) {
	if strings.HasSuffix(target, "-first") {
		mode = "A"
	} else if strings.HasSuffix(target, "-second") {
		mode = "B"
	}
	return
}

// Helpers for Age Group targets

func GetAgeGroupFromTarget(target string) (output string) {
	if strings.HasPrefix(target, "vaccinations-") &&
		(strings.HasSuffix(target, "-first") || strings.HasSuffix(target, "-second")) {
		output = strings.TrimPrefix(target, "vaccinations-")
		output = strings.TrimSuffix(output, "-first")
		output = strings.TrimSuffix(output, "-second")
	}
	return
}

// Helpers for Region targets

func GetRegionFromTarget(target string) (output string) {
	if strings.HasPrefix(target, "vaccinations-") &&
		(strings.HasSuffix(target, "-first") || strings.HasSuffix(target, "-second")) {
		output = strings.TrimPrefix(target, "vaccinations-")
		output = strings.TrimSuffix(output, "-first")
		output = strings.TrimSuffix(output, "-second")
	}
	return
}
