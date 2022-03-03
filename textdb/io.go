package textdb

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
)

func dirExist(name string) bool {
	file, err := os.Open(name)
	if os.IsNotExist(err) {
		return false
	}
	defer file.Close()
	return true
}

func readContent(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logError("Error reading content file", filename, err)
	}
	return string(content)
}

func readMetadata(filename string) Metadata {
	reader, err := os.Open(filename)
	if err != nil {
		logError("Error reading metadata file", filename, err)
	}
	defer reader.Close()

	// Read the bytes and unmarshall into our metadata struct
	byteValue, _ := ioutil.ReadAll(reader)
	var metadata Metadata
	xml.Unmarshal(byteValue, &metadata)
	return metadata
}

func saveContent(entry TextEntry) error {
	filename := filepath.Join(entry.Path, "content.md")
	return ioutil.WriteFile(filename, []byte(entry.Content), 0644)
}

func saveMetadata(entry TextEntry) error {
	// Convert our Metadata struct to an XML string...
	xmlDeclaration := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n"
	buffer := bytes.NewBufferString(xmlDeclaration)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("  ", "    ")
	err := encoder.Encode(entry.Metadata)
	if err != nil {
		return err
	}
	// ... and save it.
	filename := filepath.Join(entry.Path, "metadata.xml")
	return ioutil.WriteFile(filename, buffer.Bytes(), 0644)
}
