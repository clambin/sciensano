package sciensano

import "sort"

func getUniqueSortedGroupNames(input map[string]struct{}) (output []string) {
	for group := range input {
		output = append(output, group)
	}
	sort.Strings(output)
	return
}
