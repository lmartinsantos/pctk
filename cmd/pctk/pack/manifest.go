package pack

import (
	"fmt"
	"os"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

// ResourceType is the type of resource that the manifest describes.
type ResourceType string

const (
	// ManifestTypeScript is a script resource.
	ManifestTypeScript ResourceType = "script"
)

// Manifest is the description of a resource.
type Manifest struct {
	Type        ResourceType
	Compression pctk.ResourceCompression
	Data        any
}

// LoadManifestFromFile loads a manifest from a file.
func LoadManifestFromFile(path string) (*Manifest, error) {
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	var man Manifest
	if err := yaml.NewDecoder(data).Decode(&man); err != nil {
		return nil, err
	}
	return &man, nil
}

func (m *Manifest) UnmarshalYAML(n *yaml.Node) error {
	var header struct {
		Type        ResourceType
		Compression string
		Data        yaml.Node
	}
	if err := n.Decode(&header); err != nil {
		return err
	}
	m.Type = header.Type
	switch m.Type {
	case ManifestTypeScript:
		m.Data = new(ScriptData)
	default:
		return fmt.Errorf("unknown manifest type: %s", m.Type)
	}
	if err := header.Data.Decode(m.Data); err != nil {
		return err
	}

	switch header.Compression {
	case "", "none", "None", "NONE":
		m.Compression = pctk.CompressionNone
	case "gzip", "Gzip", "GZIP":
		m.Compression = pctk.CompressionGzip
	default:
		return fmt.Errorf("unknown compression type: %s", header.Compression)
	}

	return nil
}
