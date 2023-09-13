package dmg

import (
	"errors"
	"testing"
)

// This is an encoded version of an AttributionResource with the following data:
var encodedAttributionData = "cnR0YQEAAAAtKrjTLgQAAAAAAABr57W5AAA0AAAAAAAuBAAAAAAAAAAACAAAAAAAl5IVp2O5R1smIeUIAAAAALU2DPUAAPQbAAAAAA=="
var expectedAttributionData = AttributionResource{
	Signature:                  [4]byte{114, 116, 116, 97}, // "rtta" - fixme
	Version:                    1,
	BeforeCompressedChecksum:   3552061997,
	BeforeCompressedLength:     1070,
	BeforeUncompressedChecksum: 3115706219,
	BeforeUncompressedLength:   3407872,
	RawPos:                     1070,
	RawLength:                  524288,
	RawChecksum:                2803208855,
	AfterCompressedChecksum:    1531427171,
	AfterCompressedLength:      149233958,
	AfterUncompressedChecksum:  4111218357,
	AfterUncompressedLength:    468975616,
}

func TestParseAttribution(t *testing.T) {
	res, err := ParseAttribution(encodedAttributionData)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if *res != expectedAttributionData {
		t.Errorf("attribution data not parsed correctly, expected %+v, got %+v", expectedAttributionData, res)
	}
}

func TestParseAttributionInvalid(t *testing.T) {
	for _, tc := range []struct {
		input       string
		expectedErr error
	}{
		{input: "abcdefgh123445678", expectedErr: ErrBadAttrBase64},
		{input: "cnR0YQIAAAAtKrjTLgQAAAAAAABr57W5AAA0AAAAAAAuBAAAAAAAAAAACg==", expectedErr: ErrBadAttrLength},
		// Same as encodedAttributionData, except signature is zzzz
		{input: "enp6egEAAAAtKrjTLgQAAAAAAABr57W5AAA0AAAAAAAuBAAAAAAAAAAACAAAAAAAl5IVp2O5R1smIeUIAAAAALU2DPUAAPQbAAAAAA==", expectedErr: ErrBadAttrSignature},
		// Same as encodedAttributionData, except version is set to 2
		{input: "cnR0YQIAAAAtKrjTLgQAAAAAAABr57W5AAA0AAAAAAAuBAAAAAAAAAAACAAAAAAAl5IVp2O5R1smIeUIAAAAALU2DPUAAPQbAAAAAA==", expectedErr: ErrBadAttrVersion},
	} {
		_, err := ParseAttribution(tc.input)
		if err == nil {
			t.Errorf("expected error")
		}

		if !errors.Is(err, tc.expectedErr) {
			t.Errorf("expected error: %s, got %s", tc.expectedErr, err)
		}
	}
}
