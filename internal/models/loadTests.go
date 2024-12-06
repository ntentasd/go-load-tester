package internal

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

type Definition struct {
	Host      string     `yaml:"host"`
	LoadTests []LoadTest `yaml:"load_tests"`
}

type LoadTest struct {
	Method   string `yaml:"method"`
	Endpoint string `yaml:"endpoint"`
}

var methods = map[string]struct{}{
	"GET":    {},
	"POST":   {},
	"PUT":    {},
	"DELETE": {},
}

func (d *Definition) Dbg() {
	fmt.Println("Host:", d.Host)
	fmt.Println("Tests:")
	for _, test := range d.LoadTests {
		fmt.Printf("- %s\n  %s\n", test.Method, test.Endpoint)
	}
}

func UnmarshalYaml(r io.Reader) (Definition, error) {
	var root yaml.Node
	err := yaml.NewDecoder(r).Decode(&root)
	if err != nil {
		return Definition{}, fmt.Errorf("failed to decode YAML: %w", err)
	}

	testsNode, err := findTestsKey(root.Content[0])
	if err != nil {
		return Definition{}, err
	}

	var t Definition
	err = root.Decode(&t)
	if err != nil {
		return Definition{}, err
	}

	for i, test := range t.LoadTests {
		if i >= len(testsNode.Content) {
			return Definition{}, fmt.Errorf("unexpected number of tests")
		}
		testNode := testsNode.Content[i]
		var methodNode *yaml.Node
		for j := 0; j < len(testNode.Content); j += 2 {
			if testNode.Content[j].Value == "method" {
				methodNode = testNode.Content[j+1]
				break
			}
		}
		if methodNode == nil {
			return Definition{}, fmt.Errorf("'method' key missing in test on line %d", testNode.Line)
		}

		method := stringChain(test.Method, strings.ToUpper, strings.TrimSpace)
		t.LoadTests[i].Method = method
		if _, valid := methods[method]; !valid {
			return Definition{}, fmt.Errorf("invalid HTTP method '%s' on line %d", test.Method, methodNode.Line)
		}
	}

	return t, nil
}

func findTestsKey(node *yaml.Node) (*yaml.Node, error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected a mapping node, got %v", node.Kind)
	}

	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]

		if key.Value == "load_tests" {
			if value.Kind != yaml.SequenceNode {
				return nil, fmt.Errorf("'load_tests' must be a list (line %d)", key.Line)
			}
			return value, nil
		}
	}

	return nil, fmt.Errorf("'load_tests' key not found")
}

func stringChain(str string, functions ...func(string) string) string {
	res := str
	for _, fn := range functions {
		res = fn(res)
	}

	return res
}
