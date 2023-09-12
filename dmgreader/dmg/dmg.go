package dmg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"howett.net/plist"
)

var (
	ErrNoPropertyList = errors.New("dmg: no XML property list")
	ErrNoResourceFork = errors.New("dmg: no resource fork")
)

// DMG is a structure representing a DMG file and its related metadata.
type DMG struct {
	Koly      *KolyBlock
	Resources *Resources
	Data      []byte
}

// DMGFile is a structure representing a DMG file on a filesystem.
type DMGFile struct {
	name string
	file *os.File
}

// OpenFile returns a new instance to interact with a DMG file. This function
// might return an error if the file does not exist or if there is a problem
// opening it. When you create an instance of DMG, you must close it.
func OpenFile(name string) (*DMGFile, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("dmg: %w", err)
	}

	return &DMGFile{name: name, file: file}, nil
}

func (d *DMGFile) Parse() (*DMG, error) {
	return ParseDMG(d.file)
}

// Close closes the DMG file.
func (d *DMGFile) Close() error {
	if d.file != nil {
		return d.file.Close()
	}

	return nil
}

// TODO: probably should break this down a bit more for unit testing purposes
func ParseDMG(input ReaderSeeker) (*DMG, error) {
	dmg := new(DMG)

	block, err := parseKolyBlock(input)
	if err != nil {
		return dmg, fmt.Errorf("dmg: %w", err)
	}

	if block.XMLLength == 0 {
		return dmg, ErrNoPropertyList
	}

	buf := make([]byte, block.XMLLength)
	if _, err := input.ReadAt(buf, int64(block.XMLOffset)); err != nil {
		return dmg, fmt.Errorf("dmg: %w", err)
	}

	var data map[string]interface{}
	if err := plist.NewDecoder(bytes.NewReader(buf)).Decode(&data); err != nil {
		return dmg, fmt.Errorf("dmg: %w", err)
	}

	fork, ok := data["resource-fork"].(map[string]interface{})

	if !ok {
		return dmg, ErrNoResourceFork
	}

	resources, err := parseResources(fork)
	if err != nil {
		return dmg, fmt.Errorf("dmg: %w", err)
	}

	input.Seek(0, io.SeekStart)
	dmgData, err := io.ReadAll(input)
	if err != nil {
		return dmg, fmt.Errorf("dmg: %w", err)
	}

	dmg.Koly = block
	dmg.Resources = resources
	dmg.Data = dmgData

	return dmg, nil
}
