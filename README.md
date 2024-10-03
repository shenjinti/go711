# go711

This is a simple G711 codec implementation in Go. 
It supports both A-law (PCMA )and μ-law (PCMU) encoding and decoding.

## Installation
```bash
go get github.com/shenjinti/go711
```

## Usage Example
```go
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
```