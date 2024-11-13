package zeebe

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func MustReadFile(resourceFile string) ([]byte, error) {
	contents, err := res.ReadFile("resources/" + resourceFile)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func toCamelCase(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func unMarshalBpmn(bpmnContent []byte) ([]BPMNProcess, error) {
	var bpmn BPMNDocument
	if err := xml.Unmarshal(bpmnContent, &bpmn); err != nil {
		return nil, fmt.Errorf("failed to parse BPMN content: %w", err)
	}
	return bpmn.Processes, nil
}

func insertGeneratedCode(filePath, generateCode, containString string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if strings.Contains(line, containString) {
			lines = append(lines, generateCode)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(filePath, []byte(output), 0o644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func getModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func readFormFile(filePath string) (*Form, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open form file: %w", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read form file: %w", err)
	}

	var form Form
	if err := json.Unmarshal(byteValue, &form); err != nil {
		return nil, fmt.Errorf("failed to unmarshal form JSON: %w", err)
	}

	return &form, nil
}

func generateStructCode(form *Form, name string) string {
	structCode := "type FormData" + name + " struct {\n"
	for _, component := range form.Components {
		if component.Key == "" {
			continue
		}
		fieldName := toCamelCaseForm(component.Key)
		fieldType := getFieldType(component.Type)
		tag := fmt.Sprintf("`json:\"%s", component.Key)
		if component.Validate.Required {
			tag += "\" validate:\"required\"`"
		} else {
			tag += ",omitempty\"`"
		}
		structCode += fmt.Sprintf("\t%s %s %s\n", fieldName, fieldType, tag)
	}
	structCode += "}\n"
	return structCode
}

func toCamelCaseForm(snakeStr string) string {
	parts := strings.Split(snakeStr, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func getFieldType(fieldType string) string {
	switch fieldType {
	case "textfield", "textarea":
		return "string"
	case "number":
		return "float64"
	case "select":
		return "string"
	default:
		return "interface{}"
	}
}
