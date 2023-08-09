package annotations

import (
	"fmt"
	"regexp"
	"strings"
)

func Parse(annotationCode string) *Annotation {
	re := regexp.MustCompile(`@\w+(\([^\)]*\))?`)
	matches := re.FindStringSubmatch(annotationCode)
	if len(matches) < 1 {
		return nil
	}

	nameAndData := strings.TrimPrefix(matches[0], "@")
	parts := strings.SplitN(nameAndData, "(", 2)
	name := parts[0]
	var data string
	if len(parts) > 1 {
		data = strings.TrimSuffix(parts[1], ")")
	}

	dataMap := make(map[string]interface{})
	if data != "" {
		for i, pair := range strings.Split(data, ",") {
			if strings.Contains(pair, "=") {
				keyValue := strings.SplitN(pair, "=", 2)
				key := strings.TrimSpace(keyValue[0])
				value := strings.TrimSpace(keyValue[1])
				if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
					arrayValues := strings.Trim(value, "[]")
					dataMap[key] = strings.Split(arrayValues, ",")
				} else if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
					mapValues := strings.Trim(value, "{}")
					subMap := make(map[string]string)
					for _, subPair := range strings.Split(mapValues, ",") {
						subKeyValue := strings.SplitN(subPair, "=", 2)
						subMap[subKeyValue[0]] = subKeyValue[1]
					}
					dataMap[key] = subMap
				} else {
					dataMap[key] = value
				}
			} else {
				dataMap[fmt.Sprintf("%d", i)] = pair
			}
		}
	}

	return &Annotation{
		Name: name,
		Data: dataMap,
	}
}
