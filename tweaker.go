/*
 * Copyright (C) 2016 Lawrence Woodman <lwoodman@vlifesystems.com>
 */
package main

import (
	"fmt"
	"github.com/lawrencewoodman/dlit_go"
	"sort"
	"strings"
)

func TweakRules(
	sortedRules []*Rule,
	fieldDescriptions map[string]*FieldDescription,
) []*Rule {
	numRulesPerGroup := 3
	groupedRules :=
		groupTweakableRules(sortedRules, numRulesPerGroup)
	return tweakRules(groupedRules, fieldDescriptions)
}

func groupTweakableRules(
	sortedRules []*Rule,
	numPerGroup int,
) map[string][]*Rule {
	groups := make(map[string][]*Rule)
	for _, rule := range sortedRules {
		isTweakable, fieldName, operator, _ := rule.GetTweakableParts()
		if isTweakable {
			groupID := fmt.Sprintf("%s-%s", fieldName, operator)
			if len(groups[groupID]) < numPerGroup {
				groups[groupID] = append(groups[groupID], rule)
			}
		}
	}
	return groups
}

func tweakRules(
	groupedRules map[string][]*Rule,
	fieldDescriptions map[string]*FieldDescription,
) []*Rule {
	newRules := make([]*Rule, 0)
	for _, rules := range groupedRules {
		firstRule := rules[0]
		comparisonPoints := makeComparisonPoints(rules, fieldDescriptions)
		for _, point := range comparisonPoints {
			newRule, err := firstRule.CloneWithValue(point)
			if err != nil {
				panic(fmt.Sprintf("Can't tweak rule: %s - %s", firstRule, err))
			}
			newRules = append(newRules, newRule)
		}
	}
	return newRules
}

func dlitInSlices(needle *dlit.Literal, haystacks ...[]*dlit.Literal) bool {
	for _, haystack := range haystacks {
		for _, v := range haystack {
			if needle.String() == v.String() {
				return true
			}
		}
	}
	return false
}

// TODO: Share similar code with generaters such as generateInt
func makeComparisonPoints(
	rules []*Rule,
	fieldDescriptions map[string]*FieldDescription,
) []string {
	var minInt int64
	var maxInt int64
	var minFloat float64
	var maxFloat float64
	var field string
	var tweakableValue string

	numbers := make([]*dlit.Literal, len(rules))
	newPoints := make([]*dlit.Literal, 0)
	for i, rule := range rules {
		_, field, _, tweakableValue = rule.GetTweakableParts()
		numbers[i] = dlit.MustNew(tweakableValue)
	}

	numNumbers := len(numbers)
	sortNumbers(numbers)

	if fieldDescriptions[field].Kind == INT {
		for numI, numJ := 0, 1; numJ < numNumbers; numI, numJ = numI+1, numJ+1 {
			vI := numbers[numI]
			vJ := numbers[numJ]
			vIint, _ := vI.Int()
			vJint, _ := vJ.Int()
			if vIint < vJint {
				minInt = vIint
				maxInt = vJint
			} else {
				minInt = vJint
				maxInt = vIint
			}

			diff := maxInt - minInt
			step := diff / 10
			if diff < 10 {
				step = 1
			}

			for i := step; i < diff; i += step {
				newNum := dlit.MustNew(minInt + i)
				if !dlitInSlices(newNum, numbers, newPoints) {
					newPoints = append(newPoints, newNum)
				}
			}
		}
	} else {
		for numI, numJ := 0, 1; numJ < numNumbers; numI, numJ = numI+1, numJ+1 {
			vI := numbers[numI]
			vJ := numbers[numJ]
			vIfloat, _ := vI.Float()
			vJfloat, _ := vJ.Float()
			if vIfloat < vJfloat {
				minFloat = vIfloat
				maxFloat = vJfloat
			} else {
				minFloat = vJfloat
				maxFloat = vIfloat
			}

			diff := maxFloat - minFloat
			step := diff / 10.0
			for i := step; i < diff; i += step {
				newNum := dlit.MustNew(minFloat + i)
				if !dlitInSlices(newNum, numbers, newPoints) {
					newPoints = append(newPoints, newNum)
				}
			}
		}
	}
	maxDP := fieldDescriptions[field].MaxDP
	return arrayDlitsToStrings(newPoints, maxDP)
}

func arrayDlitsToStrings(lits []*dlit.Literal, maxDP int) []string {
	r := make([]string, len(lits))
	for i, l := range lits {
		s := l.String()
		// TODO: limit Dec places 0..maxDP
		r[i] = limitStringDecPlaces(s, maxDP)
	}
	return r
}

func limitStringDecPlaces(s string, maxDP int) string {
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return s[:i+maxDP+1]
	}
	return s
}

// byNumber implements sort.Interface for []*dlit.Literal
type byNumber []*dlit.Literal

func (l byNumber) Len() int { return len(l) }
func (l byNumber) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l byNumber) Less(i, j int) bool {
	lI := l[i]
	lJ := l[j]
	iI, lIisInt := lI.Int()
	iJ, lJisInt := lJ.Int()
	if lIisInt && lJisInt {
		if iI < iJ {
			return true
		}
		return false
	}

	fI, lIisFloat := lI.Float()
	fJ, lJisFloat := lJ.Float()

	if lIisFloat && lJisFloat {
		if fI < fJ {
			return true
		}
		return false
	}
	panic(fmt.Sprintf("Can't compare numbers: %s, %s", lI, lJ))
}

func sortNumbers(nums []*dlit.Literal) {
	sort.Sort(byNumber(nums))
}
