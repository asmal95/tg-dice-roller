package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var compRegEx = regexp.MustCompile(`^(?P<mod>[+-])?(?P<count>[0-9]+)?d(?P<edge>[0-9]+)(?P<additional>[+-][0-9]+)?$`)
var randomSource = rand.NewSource(time.Now().UnixNano())
var randGenerator = rand.New(randomSource)

// d20
// 2d20 = (d20 + d20)
// 2d20+3 = (d20 + d20) + 3
// +d20 = max(d20, d20)
// -d20 = min(d20, d20)
// returns result number, explanation string and error for error case
// result number is calculated result
// explanation string is visualisation of the calculation, example for 2d20+3: (10 + 9) + 3
func Roll(str string) (int, string, error) {
	dict := getParams(str)
	count := 1
	edge := 0
	explanation := ""
	if val := dict["count"]; val != "" {
		count, _ = strconv.Atoi(val)
	}
	if val := dict["edge"]; val != "" {
		edge, _ = strconv.Atoi(val)
	} else {
		return 0, "", fmt.Errorf("incorrect exression: %v", str)
	}
	toAdd := 0
	if additional := dict["additional"]; additional != "" {
		toAdd, _ = strconv.Atoi(additional)
	}
	mod := dict["mod"]

	res, explanation := roll(edge, count, toAdd, mod != "")

	if mod != "" {
		modRes, modExplanation := roll(edge, count, toAdd, true)
		switch mod {
		case "-":
			if modRes < res {
				res = modRes
			}
		case "+":
			if modRes > res {
				res = modRes
			}
		}

		explanation += " | " + modExplanation
	}

	return res, explanation, nil
}

func roll(edge, count, additional int, explainSingle bool) (int, string) {
	res := 0
	explanation := ""
	if count > 1 {
		explanation += "("
	}
	for i := 0; i < count; i++ {
		r := randomDice(edge)
		res += r
		if count > 1 || additional != 0 || explainSingle {
			explanation += strconv.Itoa(r)
		}
		if i < count-1 {
			explanation += " + "
		}
	}
	if count > 1 {
		explanation += ")"
	}
	res += additional
	if additional > 0 {
		explanation += " + " + strconv.Itoa(additional)
	} else if additional < 0 {
		explanation += " - " + strconv.Itoa(-additional)
	}
	return res, explanation
}

func randomDice(edge int) int {
	if edge <= 0 {
		return edge
	}
	return randGenerator.Intn(edge) + 1
}

func getParams(url string) (paramsMap map[string]string) {
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
