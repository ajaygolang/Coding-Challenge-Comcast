package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Input represents the input JSON structure
type Input map[string]interface{}

// Output represents the desired output JSON structure
type Output []map[string]interface{}

func main() {
	// Read input JSON from stdin
	var inputJSON Input
	err := json.NewDecoder(os.Stdin).Decode(&inputJSON)
	if err != nil {
		log.Fatalf("error decoding input JSON: %v", err)
	}

	// Transform input JSON to desired output format
	output := transformInput(inputJSON)

	// Print output JSON to stdout
	printOutput(output)
}

// transformInput transforms the input JSON to the desired output format
func transformInput(input Input) Output {
	var output Output

	// Iterate through input keys and transform each field
	for key, value := range input {
		// Skip fields with empty keys
		if key == "" {
			continue
		}

		// Sanitize key by trimming leading and trailing whitespace
		key = strings.TrimSpace(key)

		// Transform value based on data type
		switch v := value.(type) {
		case map[string]interface{}:
			outputMap := transformMap(v)
			if len(outputMap) > 0 {
				output = append(output, outputMap)
			}
		case string:
			if ts, err := time.Parse(time.RFC3339, v); err == nil {
				output = append(output, map[string]interface{}{key: ts.Unix()})
			} else {
				output = append(output, map[string]interface{}{key: strings.TrimSpace(v)})
			}
		case []interface{}:
			outputList := transformList(v)
			if len(outputList) > 0 {
				output = append(output, map[string]interface{}{key: outputList})
			}
		default:
			fmt.Printf("Warning: Skipping unsupported data type for key %q\n", key)
		}
	}

	return output
}

// transformMap transforms a map[string]interface{} to the desired output format
func transformMap(m map[string]interface{}) map[string]interface{} {
	outputMap := make(map[string]interface{})

	// Sort map keys lexically
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate through sorted keys and transform each field
	for _, k := range keys {
		// Sanitize key by trimming leading and trailing whitespace
		key := strings.TrimSpace(k)

		// Transform value based on data type
		switch v := m[k].(type) {
		case map[string]interface{}:
			outputMap[key] = transformMap(v)
		case string:
			outputMap[key] = strings.TrimSpace(v)
		case []interface{}:
			outputList := transformList(v)
			if len(outputList) > 0 {
				outputMap[key] = outputList
			}
		default:
			fmt.Printf("Warning: Skipping unsupported data type for key %q\n", key)
		}
	}

	return outputMap
}

// transformList transforms a []interface{} to the desired output format
func transformList(l []interface{}) []interface{} {
	var outputList []interface{}

	// Iterate through list elements and transform each item
	for _, item := range l {
		switch v := item.(type) {
		case map[string]interface{}:
			outputMap := transformMap(v)
			if len(outputMap) > 0 {
				outputList = append(outputList, outputMap)
			}
		case string:
			if ts, err := time.Parse(time.RFC3339, v); err == nil {
				outputList = append(outputList, ts.Unix())
			} else if isNumeric(v) {
				outputList = append(outputList, parseNumber(v))
			} else {
				outputList = append(outputList, strings.TrimSpace(v))
			}
		default:
			fmt.Printf("Warning: Skipping unsupported data type in list\n")
		}
	}

	return outputList
}

// isNumeric checks if a string represents a numeric value
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// parseNumber parses a numeric string and returns the corresponding number
func parseNumber(s string) interface{} {
	// Strip leading zeros
	trimmed := strings.TrimLeft(s, "0")
	// Parse integer or float
	if strings.Contains(trimmed, ".") {
		f, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return nil
		}
		return f
	}
	i, err := strconv.Atoi(trimmed)
	if err != nil {
		return nil
	}
	return i
}

// printOutput prints the output JSON to stdout
func printOutput(output Output) {
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("error encoding output JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}
