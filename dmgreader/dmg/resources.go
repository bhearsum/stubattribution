package dmg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type ResourceData struct {
	Attributes string
	CFName     string
	Data       []uint8
	ID         string
	Name       string
}

type Resources struct {
	Entries map[string][]ResourceData
}

var (
	ErrResourceNotFound = errors.New("resources: named resource not found")
)

func parseResources(unparsed map[string]interface{}) (*Resources, error) {
	res := new(Resources)
	err := mapstructure.Decode(unparsed, &res.Entries)
	if err != nil {
		return res, fmt.Errorf("resources: %w", err)
	}

	// If present, the plst resource's Name field is treated specially, as it may
	// contain attribution information, which is really a base64 encoded binary object
	// printed as text. mapstructure treats the Name field as a string, so we must
	// do some post-processing to strip any whitespace that may be present.
	if plst, ok := res.Entries["plst"]; ok {
		for i, _ := range plst {
			if res.Entries["plst"][i].Name != "" {
				res.Entries["plst"][i].Name = strings.ReplaceAll(strings.ReplaceAll(res.Entries["plst"][i].Name, "\t", ""), "\n", "")
			}
		}
	}

	return res, nil
}

func (r *Resources) GetResourceDataByName(name string) ([]ResourceData, error) {
	for k, v := range r.Entries {
		if k == name {
			return v, nil
		}
	}

	return []ResourceData{}, ErrResourceNotFound
}
