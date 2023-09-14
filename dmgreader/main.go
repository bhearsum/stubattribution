package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"hash/crc32"

	"github.com/willdurand/go-dmg-reader/dmg"
	"howett.net/plist"
)

func gf2MatrixTimes(mat []uint, vec uint) uint {
	var sum uint

	for vec != 0 {
		if vec&1 != 0 {
			sum ^= mat[0]
		}
		vec >>= 1
		mat = mat[1:]
	}
	return sum
}

func gf2MatrixSquare(square, mat []uint) {
	for n := 0; n < 32; n++ {
		square[n] = gf2MatrixTimes(mat, mat[n])
	}
}

// CRC32Combine returns the combined CRC-32 hash value of the two passed CRC-32
// hash values crc1 and crc2. poly represents the generator polynomial
// and len2 specifies the byte length that the crc2 hash covers.
func CRC32Combine(poly uint32, crc1, crc2 uint32, len2 int64) uint32 {
	// degenerate case (also disallow negative lengths)
	if len2 <= 0 {
		return crc1
	}

	even := make([]uint, 32) // even-power-of-two zeros operator
	odd := make([]uint, 32)  // odd-power-of-two zeros operator

	// put operator for one zero bit in odd
	odd[0] = uint(poly) // CRC-32 polynomial
	row := uint(1)
	for n := 1; n < 32; n++ {
		odd[n] = row
		row <<= 1
	}

	// put operator for two zero bits in even
	gf2MatrixSquare(even, odd)

	// put operator for four zero bits in odd
	gf2MatrixSquare(odd, even)

	// apply len2 zeros to crc1 (first square will put the operator for one
	// zero byte, eight zero bits, in even)
	crc1n := uint(crc1)
	for {
		// apply zeros operator for this bit of len2
		gf2MatrixSquare(even, odd)
		if len2&1 != 0 {
			crc1n = gf2MatrixTimes(even, crc1n)
		}
		len2 >>= 1

		// if no more bits set, then done
		if len2 == 0 {
			break
		}

		// another iteration of the loop with odd and even swapped
		gf2MatrixSquare(odd, even)
		if len2&1 != 0 {
			crc1n = gf2MatrixTimes(odd, crc1n)
		}
		len2 >>= 1

		// if no more bits set, then done
		if len2 == 0 {
			break
		}
	}

	// return combined crc
	crc1n ^= uint(crc2)
	return uint32(crc1n)
}

func main() {
	file, err := dmg.OpenFile("/home/bhearsum/repos/stubattribution/dmgreader/tests/fixtures/attributable.dmg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := file.Parse()
	if err != nil {
		panic(err)
	}

	blkxRes, _ := data.Resources.GetResourceDataByName("blkx")
	plstRes, _ := data.Resources.GetResourceDataByName("plst")

	attr, err := dmg.ParseAttribution(plstRes[0].Name)
	if err != nil {
		log.Printf("err: %s", err)
	}


	rawBlock := string(data.Data[attr.RawPos:attr.RawPos + attr.RawLength])
	attrOffset := strings.Index(rawBlock, "__MOZCUSTOM__abcdef")
	replacement := []byte("__MOZCUSTOM__bhearsum was here")
	fullAttrOffset := int(attr.RawPos) + attrOffset
	// TODO: need to zero out the area
	copy(data.Data[fullAttrOffset:fullAttrOffset+len(replacement)], replacement[:])

	// generic new crc value for raw block
	rawCrc := crc32.Checksum(data.Data[attr.RawPos:attr.RawPos + attr.RawLength], crc32.MakeTable(0xffffffff))

	// pull out blkx metadata
	blkx, err := dmg.ParseBlkxData(blkxRes[3].Data)
	if err != nil {
		log.Printf("err: %s", err)
	}

	// combine checksums
	blkx.Table.Checksum.Data[0] = CRC32Combine(0xffffffff, CRC32Combine(0xffffffff, attr.BeforeUncompressedChecksum, rawCrc, int64(attr.RawLength)), attr.AfterUncompressedChecksum, int64(attr.AfterUncompressedLength))
	// update resources in data.Data
	// TODO: this is not right -- we need to somehow include the new checksum we calculate here.
	// i think the new checksum is in `blkx`, but we're obviously not pulling that in...
	resources := new(bytes.Buffer)
	plist.NewEncoder(resources).Encode(data.Data[data.Koly.XMLOffset:data.Koly.XMLOffset+data.Koly.XMLLength])

	var newData map[string]interface{}
	newData["resource-fork"] = resources

	// update koly block checksum
	// update koly block in data

	os.Stdout.Write(data.Data[:])
}
