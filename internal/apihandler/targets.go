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
		"tests":   GetTestTargets(),
		"vaccine": GetVaccinationsTargets(),
		"vac-age": GetVaccinationsByAgeTargets(),
		"vac-reg": GetVaccinationsByRegionTargets(),
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

func getModeFromTarget(target string) (mode string) {
	if strings.HasSuffix(target, "-first") {
		mode = "A"
	} else if strings.HasSuffix(target, "-second") {
		mode = "B"
	}
	return
}

// Helpers for Age Group targets

func getAgeGroupFromTarget(target string) (ageGroup string) {
	if strings.HasPrefix(target, "vaccinations-") &&
		(strings.HasSuffix(target, "-first") || strings.HasSuffix(target, "-second")) {
		ageGroup = strings.TrimPrefix(target, "vaccinations-")
		ageGroup = strings.TrimSuffix(ageGroup, "-first")
		ageGroup = strings.TrimSuffix(ageGroup, "-second")
	}
	return
}

// Helpers for Region targets

func getRegionFromTarget(target string) (region string) {
	if strings.HasPrefix(target, "vaccinations-") &&
		(strings.HasSuffix(target, "-first") || strings.HasSuffix(target, "-second")) {
		region = strings.TrimPrefix(target, "vaccinations-")
		region = strings.TrimSuffix(region, "-first")
		region = strings.TrimSuffix(region, "-second")
	}
	return
}
