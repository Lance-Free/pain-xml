package schematest

import (
	"encoding/xml"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

const schemaFile = "pain.001.001.03.xsd"

// Validate fails the test if document does not match the pain.001.001.03 schema. It skips the test if xmllint is not
// installed.
func Validate(t *testing.T, document any) {
	t.Helper()

	xmllint, err := exec.LookPath("xmllint")
	if err != nil {
		t.Skip("xmllint is not installed")
	}

	encoded, err := xml.Marshal(document)
	if err != nil {
		t.Fatalf("marshalling document: %v", err)
	}

	path := filepath.Join(t.TempDir(), "document.xml")
	if err := os.WriteFile(path, append([]byte(xml.Header), encoded...), 0o600); err != nil {
		t.Fatalf("writing document: %v", err)
	}

	output, err := exec.Command(xmllint, "--noout", "--schema", schemaPath(), path).CombinedOutput()
	if err != nil {
		t.Fatalf("document does not match the schema: %v\n%s\n%s", err, output, encoded)
	}
}

func schemaPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", schemaFile)
}
