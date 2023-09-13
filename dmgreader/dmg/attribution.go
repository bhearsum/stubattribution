package dmg

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// TODO: this should be attr, but only after we start flipping the signature properly in libdmg-hfsplus
	attrBlockSignature = "rtta"
	attrBlockVersion   = 1
	attrBlockSize      = 76
)

var (
	ErrBadAttrBase64     = errors.New("attr: couldn't decode base64 data")
	ErrBadAttrLength     = errors.New("attr: bad object length")
	ErrBadAttrBinaryData = errors.New("attr: couldn't parse binary attribution data")
	ErrBadAttrSignature  = errors.New("attr: invalid attribution signature")
	ErrBadAttrVersion    = errors.New("attr: invalid attribution resource version")
)

type AttributionResource struct {
	Signature                  [4]byte
	Version                    uint32
	BeforeCompressedChecksum   uint32
	BeforeCompressedLength     uint64
	BeforeUncompressedChecksum uint32
	BeforeUncompressedLength   uint64
	RawPos                     uint64
	RawLength                  uint64
	RawChecksum                uint32
	AfterCompressedChecksum    uint32
	AfterCompressedLength      uint64
	AfterUncompressedChecksum  uint32
	AfterUncompressedLength    uint64
}

func ParseAttribution(raw string) (*AttributionResource, error) {
	attr := new(AttributionResource)

	if raw == "" {
		return attr, nil
	}

	buf, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return attr, ErrBadAttrBase64
	}

	if len(buf) != attrBlockSize {
		return attr, ErrBadAttrLength
	}

	if err := binary.Read(bytes.NewReader(buf), binary.LittleEndian, attr); err != nil {
		return attr, fmt.Errorf("attr: %w", err)
	}

	if !bytes.Equal(attr.Signature[:], []byte(attrBlockSignature)) {
		return attr, ErrBadAttrSignature
	}

	if attr.Version != attrBlockVersion {
		return attr, ErrBadAttrVersion
	}

	return attr, nil
}
