package go711

import (
	"testing"
)

func TestPCMA(t *testing.T) {
	pcmBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	// Encode PCM bytes to A-law bytes
	alawBytes, err := EncodePCMA(pcmBytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(alawBytes) != 4 {
		t.Fatalf("a-law bytes length is not 4: %d", len(alawBytes))
	}
	// Decode A-law bytes to PCM bytes
	pcmBytes, err = DecodePCMA(alawBytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(pcmBytes) != 8 {
		t.Fatal("PCM bytes length is not 8", len(pcmBytes))
	}
}
func TestPCMU(t *testing.T) {
	pcmBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	// Encode PCM bytes to U-law bytes
	ulawBytes, err := EncodePCMU(pcmBytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(ulawBytes) != 4 {
		t.Fatalf("u-law bytes length is not 4: %d", len(ulawBytes))
	}
	// Decode U-law bytes to PCM bytes
	pcmBytes, err = DecodePCMU(ulawBytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(pcmBytes) != 8 {
		t.Fatal("PCM bytes length is not 8", len(pcmBytes))
	}
}
