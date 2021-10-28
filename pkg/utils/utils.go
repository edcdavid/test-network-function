package utils

import (
	"path"
	"strings"
	"time"

	"github.com/onsi/gomega"
	"github.com/test-network-function/test-network-function/pkg/tnf"
	"github.com/test-network-function/test-network-function/pkg/tnf/handlers/generic"
	"github.com/test-network-function/test-network-function/pkg/tnf/interactive"
)

const (
	ocCommandTimeOut = time.Second * 10
)

var (
	// TestFile is the file location of the command.json test case relative to the project root.
	TestFile = path.Join("pkg", "tnf", "handlers", "command", "command.json")
	// RelativeSchemaPath is the relative path to the generic-test.schema.json JSON schema.
	relativeSchemaPath = path.Join(pathRelativeToRoot, schemaPath)
	// PathRelativeToRoot is used to calculate relative filepaths for the `test-network-function` executable entrypoint.
	pathRelativeToRoot = path.Join("..")
	// pathToTestFile is the relative path to the command.json test case.
	pathToTestFile = path.Join(pathRelativeToRoot, TestFile)
	// schemaPath is the path to the generic-test.schema.json JSON schema relative to the project root.
	schemaPath = path.Join("schemas", "generic-test.schema.json")
)

// ArgListToMap takes a list of strings of the form "key=value" and translate it into a map
// of the form {key: value}
func ArgListToMap(lst []string) map[string]string {
	retval := make(map[string]string)
	for _, arg := range lst {
		splitArgs := strings.Split(arg, "=")
		if len(splitArgs) == 1 {
			retval[splitArgs[0]] = ""
		} else {
			retval[splitArgs[0]] = splitArgs[1]
		}
	}
	return retval
}

// FilterArray takes a list and a predicate and returns a list of all elements for whom the predicate returns true
func FilterArray(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
func ExecuteCommand(command string) (string, error) {
	values := make(map[string]interface{})
	values["COMMAND"] = command
	values["TIMEOUT"] = ocCommandTimeOut.Nanoseconds()
	context := interactive.GetContext()
	tester, handler, result, err := generic.NewGenericFromMap(pathToTestFile, relativeSchemaPath, values)

	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(result).ToNot(gomega.BeNil())
	gomega.Expect(result.Valid()).To(gomega.BeTrue())
	gomega.Expect(handler).ToNot(gomega.BeNil())
	gomega.Expect(tester).ToNot(gomega.BeNil())

	test, err := tnf.NewTest(context.GetExpecter(), *tester, handler, context.GetErrorChannel())
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(tester).ToNot(gomega.BeNil())
	if err != nil {
		return "", err
	}
	test.RunAndValidate()

	genericTest := (*tester).(*generic.Generic)
	gomega.Expect(genericTest).ToNot(gomega.BeNil())

	matches := genericTest.Matches
	gomega.Expect(len(matches)).To(gomega.Equal(1))
	match := genericTest.GetMatches()[0]
	return match.Match, nil
}
