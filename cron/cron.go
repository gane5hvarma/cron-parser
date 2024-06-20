package cron

import (
	"fmt"
	"strconv"
	"strings"
)

type CronParser struct {
	minutes     []string
	hours       []string
	daysOfMonth []string
	months      []string
	daysOfWeek  []string
	command     string
}

func ParseCronExp(cronExp string) (cp CronParser, err error) {
	fields := strings.Fields(cronExp)
	if len(fields) != 6 {
		fmt.Println("not a valid expression")
	}
	minute := fields[0]
	hour := fields[1]
	dayOfMonth := fields[2]
	month := fields[3]
	dayOfWeek := fields[4]
	cp.command = fields[5]
	if cp.minutes, err = parseField(minute, 0, 59); err != nil {
		return cp, err
	}

	if cp.hours, err = parseField(hour, 0, 23); err != nil {
		return cp, err
	}

	if cp.daysOfMonth, err = parseField(dayOfMonth, 1, 31); err != nil {
		return cp, err
	}

	if cp.months, err = parseField(month, 1, 12); err != nil {
		return cp, err
	}

	if cp.daysOfWeek, err = parseField(dayOfWeek, 0, 6); err != nil {
		return cp, err
	}
	return cp, nil
}

func (cp CronParser) Print() {
	data := [][]string{
		{"minute", strings.Join(cp.minutes, " ")},
		{"hour", strings.Join(cp.hours, " ")},
		{"day of month", strings.Join(cp.daysOfMonth, " ")},
		{"month", strings.Join(cp.months, " ")},
		{"day of week", strings.Join(cp.daysOfWeek, " ")},
		{"command", cp.command},
	}
	for i := 0; i < len(data); i++ {
		fmt.Printf("%-14s %s\n", data[i][0], data[i][1])
	}
}

func genStringRangeWithStep(start int, end int, skip int) []string {
	var numRange []string
	for i := start; i <= end; i = i + skip {
		numRange = append(numRange, strconv.Itoa(i))
	}
	return numRange
}

func addIfDoesntExist(existingElementsMap map[string]bool, elements []string, toBeAdded []string) []string {
	for _, element := range toBeAdded {
		if ok := existingElementsMap[element]; !ok {
			existingElementsMap[element] = true
			elements = append(elements, element)
		}
	}
	return elements
}

func parseRange(part string, min int, max int, step int) (values []string, err error) {
	start := min
	end := max
	bound := strings.Split(part, "-") // 0-5
	if len(bound) != 2 {
		return []string{}, fmt.Errorf("not a valid expression")
	}
	start, err = strconv.Atoi(bound[0]) // 0
	if err != nil {
		return []string{}, err
	}
	end, err = strconv.Atoi(bound[1]) // 5
	if err != nil {
		return []string{}, err
	}
	if start > end {
		return []string{}, fmt.Errorf("not a valid cron expression")
	}
	if start < min {
		return []string{}, fmt.Errorf("not a valid cron expression")
	}
	if end > max {
		return []string{}, fmt.Errorf("not a valid cron expression")
	}
	return genStringRangeWithStep(start, end, step), nil
}

// if start > end when generating range then its not a valid cron.
// *,* this is a valid expression. handle duplicates when joining results from each part separated by ',' ,
// for ex: */2,5-59/4. some values can overlap
func parseField(field string, min int, max int) ([]string, error) {
	if field == "*" { // all possible values between min and max
		return genStringRangeWithStep(min, max, 1), nil
	}
	var fieldValues []string
	var err error
	fieldSubParts := strings.Split(field, ",") // ex: *, 2-5/1 -> 0, 1, 2, ..., 59
	existentElementsMap := make(map[string]bool)
	for _, part := range fieldSubParts {
		step := 1 // valid min valu
		startValue := part
		hasSplit := strings.Contains(part, "/")
		if hasSplit {
			stepSplit := strings.Split(part, "/")
			if len(stepSplit) != 2 {
				return []string{}, fmt.Errorf("not a valid expression")
			}
			step, err = strconv.Atoi(stepSplit[1])
			if err != nil || step == 0 {
				return []string{}, fmt.Errorf("not a valid expression")
			}
			startValue = stepSplit[0] // can be 2 or 2-5
		}
		end := max
		if startValue == "*" {
			startValue = fmt.Sprintf("%d-%d", min, max) // min: 0, max: 59 -> 0-59
		}
		if strings.Contains(startValue, "-") { // range of values will be generated because of hypen
			values, err := parseRange(startValue, min, max, step)
			if err != nil {
				return []string{}, fmt.Errorf("not a valid expression")
			}
			fieldValues = addIfDoesntExist(existentElementsMap, fieldValues, values)
		} else {
			start, err := strconv.Atoi(startValue) // 5
			if err != nil {
				return []string{}, err
			}
			if start <= min || start >= max {
				return []string{}, fmt.Errorf("not a valid expression")
			}
			if hasSplit { // has step operator, range is generated
				fieldValues = addIfDoesntExist(existentElementsMap, fieldValues, genStringRangeWithStep(start, end, step))
			} else { // individual value
				fieldValues = addIfDoesntExist(existentElementsMap, fieldValues, []string{strconv.Itoa(start)})
			}
		}

	}
	return fieldValues, nil
}
