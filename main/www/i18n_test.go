package www

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestAcceptLanguageHeaderParse(t *testing.T) {
	a := struct {
		Lang   string
		Weight float64
	}{}
	log.Println(a)
	v := " de-CH,de;q=0.9,en;q=0.8,en-US;q=0.7"
	splittedList := strings.Split(v, ",")
	langList := make([]struct {
		Lang   string
		Weight float64
	}, 0)
	for i, a := range splittedList {
		b := strings.Split(strings.TrimSpace(a), ";")
		if len(b) > 0 {
			if len(b[0]) > 1 {
				lang := struct {
					Lang   string
					Weight float64
				}{Lang: b[0], Weight: 0}
				//try to read the weight
				if len(b) > 1 {
					var re = regexp.MustCompile(`q=([0-9\.]+)`)
					abb := re.FindAllStringSubmatch(b[1], -1)
					if len(abb) > 0 && len(abb[0]) > 1 {
						if s, err := strconv.ParseFloat(abb[0][1], 64); err == nil {
							lang.Weight = s
						}
					}
				}
				if lang.Weight == 0 && i == 0 {
					lang.Weight = 1
				}
				langList = append(langList, lang)
			}
		}
	}
	sort.Slice(langList, func(i, j int) bool {
		return langList[i].Weight > langList[j].Weight
	})
	log.Println(langList)
	langs := ParseAcceptLanguageHeader(`de-CH,de,en;q=0.9`)
	for i, match := range langs {
		fmt.Println(match, "found at index", i)
	}
}
