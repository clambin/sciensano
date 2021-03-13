package apihandler

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"sort"
	"strings"
)

var (
	// TODO: should be build dynamically based on content. e.g. to provide stats by age group, we'd need:
	// tests-age-18-35-positive
	// tests-ago-18-35-total
	// etc.
	targets = map[string][]string{
		"tests":   getTestTargets(),
		"vaccine": getVaccinationsTargets(),
		"vac-age": getVaccinationsByAgeTargets(),
		"vac-reg": getVaccinationsByRegionTargets(),
	}
)

func findTargetGroup(target string) string {
	for group, entries := range targets {
		for _, entry := range entries {
			if target == entry {
				return group
			}
		}
	}
	return ""
}

func allTargets() (output []string) {
	for _, entries := range targets {
		output = append(output, entries...)
	}
	sort.Strings(output)
	return
}

func getTestTargets() []string {
	return []string{
		"tests-total",
		"tests-positive",
		"tests-rate",
	}
}

func getVaccinationsTargets() []string {
	return []string{
		"vaccinations-first",
		"vaccinations-second",
	}
}

// Targets

func getVaccinationsByAgeTargets() (targets []string) {
	for _, ageGroup := range sciensano.AgeGroups {
		targets = append(targets, "vaccinations-"+ageGroup+"-first")
		targets = append(targets, "vaccinations-"+ageGroup+"-second")
	}
	return
}

func getVaccinationsByRegionTargets() (targets []string) {
	for _, region := range sciensano.Regions {
		targets = append(targets, "vaccinations-"+region+"-first")
		targets = append(targets, "vaccinations-"+region+"-second")
	}
	return
}

// Generic functions

func getModeFromTarget(target string) (mode string) {
	if strings.HasSuffix(target, "-first") {
		mode = "A"
	} else if strings.HasSuffix(target, "-second") {
		mode = "B"
	}
	return
}

// Helpers for Age Group targets

func getAgeGroupFromTarget(target string) (output string) {
	if strings.HasPrefix(target, "vaccinations-") &&
		(strings.HasSuffix(target, "-first") || strings.HasSuffix(target, "-second")) {
		output = strings.TrimPrefix(target, "vaccinations-")
		output = strings.TrimSuffix(output, "-first")
		output = strings.TrimSuffix(output, "-second")
	}
	return
}

// Helpers for Region targets

func getRegionFromTarget(target string) (output string) {
	if strings.HasPrefix(target, "vaccinations-") &&
		(strings.HasSuffix(target, "-first") || strings.HasSuffix(target, "-second")) {
		output = strings.TrimPrefix(target, "vaccinations-")
		output = strings.TrimSuffix(output, "-first")
		output = strings.TrimSuffix(output, "-second")
	}
	return
}
