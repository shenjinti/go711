package go711

const (
	BIAS       = 0x84
	SEG_SHIFT  = 4
	SEG_MASK   = 0x70
	QUANT_MASK = 0x0F
	SIGN_BIT   = 0x80
)

var segEnd = []int{0xFF, 0x1FF, 0x3FF, 0x7FF, 0xFFF, 0x1FFF, 0x3FFF, 0x7FFF}

func search(val int, table []int, size int) int {
	for i := 0; i < size; i++ {
		if val <= table[i] {
			return i
		}
	}
	return size
}

func linear2alaw(pcmVal int) byte {
	var mask int
	var seg int
	var aval int

	if pcmVal >= 0 {
		mask = 0xD5 /* sign (7th) bit = 1 */
	} else if pcmVal < -8 {
		mask = 0x55 /* sign bit = 0 */
		pcmVal = -pcmVal - 8
	} else {
		// For all values in range [-7; -1] end result will be 0^0xD5,
		// so we may just return it from here.
		// NOTE: This is not just optimization! For values in that range
		//       expression (-pcm_val - 8) will return negative result,
		//       messing all following calculations. This is old bug,
		//       coming from original code from Sun Microsystems, Inc.
		// Btw, seems there is no code for 0, value after decoding will be
		// either 8 (if we return 0xD5 here) or -8 (if we return 0x55 here).
		return 0xD5
	}
	/* Convert the scaled magnitude to segment number. */
	seg = search(pcmVal, segEnd, 8)

	/* Combine the sign, segment, and quantization bits. */

	if seg >= 8 { /* out of range, return maximum value. */
		return byte(0x7F ^ mask)
	} else {
		aval = seg << SEG_SHIFT
		if seg < 2 {
			aval |= (pcmVal >> 4) & QUANT_MASK
		} else {
			aval |= (pcmVal >> (seg + 3)) & QUANT_MASK
		}
		return byte(aval ^ mask)
	}
}

func alaw2linear(aVal byte) int16 {
	aVal ^= 0x55
	t := int16(aVal&QUANT_MASK) << 4
	seg := int16((aVal & SEG_MASK) >> SEG_SHIFT)
	switch seg {
	case 0:
		t += 8
	case 1:
		t += 0x108
	default:
		t += 0x108
		t <<= (seg - 1)
	}
	if aVal&SIGN_BIT != 0 {
		return t
	}
	return -t
}

func linear2ulaw(pcmVal int) byte {
	var mask int
	var seg int
	//unsigned char aval;
	var uval int

	/* Get the sign and the magnitude of the value. */
	if pcmVal < 0 {
		pcmVal = BIAS - pcmVal
		mask = 0x7F
	} else {
		pcmVal += BIAS
		mask = 0xFF
	}

	/* Convert the scaled magnitude to segment number. */
	seg = search(pcmVal, segEnd, 8)

	/*
	 * Combine the sign, segment, quantization bits;
	 * and complement the code word.
	 */
	if seg >= 8 { /* out of range, return maximum value. */
		return byte(0x7F ^ mask)
	} else {
		uval = (seg << 4) | ((pcmVal >> (seg + 3)) & 0xF)
		return byte(uval ^ mask)
	}
}

func ulaw2linear(uVal byte) int {
	var t int
	/* Complement to obtain normal u-law value. */
	uVal = ^uVal

	/*
	 * Extract and bias the quantization bits. Then
	 * shift up by the segment number and subtract out the bias.
	 */
	t = int((uVal&QUANT_MASK)<<3) + BIAS
	t <<= (uint8(uVal) & SEG_MASK) >> SEG_SHIFT

	if uVal&SIGN_BIT != 0 {
		return BIAS - t
	}
	return t - BIAS
}

// DecodePCMA decodes a-law encoded data to PCM data
func DecodePCMA(dataBytes []byte) ([]byte, error) {
	data := make([]byte, len(dataBytes)*2)
	j := 0
	for i := 0; i < len(dataBytes); i++ {
		pcmInt := alaw2linear(dataBytes[i])
		data[j] = byte(pcmInt)
		data[j+1] = byte(pcmInt >> 8)
		j += 2
	}
	return data, nil
}

// EncodePCMA encodes PCM data to a-law encoded data
func EncodePCMA(dataBytes []byte) ([]byte, error) {
	data := make([]byte, len(dataBytes)/2)
	j := 0
	for i := 0; i < len(dataBytes); i += 2 {
		pcmInt := int16(dataBytes[i+1])<<8 | int16(dataBytes[i])
		data[j] = linear2alaw(int(pcmInt))
		j++
	}
	return data, nil
}

// DecodePCMU decodes u-law encoded data to PCM data
func DecodePCMU(dataBytes []byte) ([]byte, error) {
	data := make([]byte, len(dataBytes)*2)
	j := 0
	for i := 0; i < len(dataBytes); i++ {
		pcmInt := ulaw2linear(dataBytes[i])
		data[j] = byte(pcmInt)
		data[j+1] = byte(pcmInt >> 8)
		j += 2
	}
	return data, nil
}

// EncodePCMU encodes PCM data to u-law encoded data
func EncodePCMU(dataBytes []byte) ([]byte, error) {
	data := make([]byte, len(dataBytes)/2)
	j := 0
	for i := 0; i < len(dataBytes); i += 2 {
		pcmInt := int16(dataBytes[i+1])<<8 | int16(dataBytes[i])
		data[j] = linear2ulaw(int(pcmInt))
		j++
	}
	return data, nil
}
