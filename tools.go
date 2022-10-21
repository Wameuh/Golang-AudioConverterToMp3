package GoAudioConverterToMp3

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/jfreymuth/oggvorbis"
)

func GetByteSlice(r io.Reader) ([]byte, *oggvorbis.Format, error) {
	oggAudio, format, err := oggvorbis.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	return convertSliceToInt16Slice(oggAudio), format, nil
}

func convertSliceToInt16Slice(mySlice []float32) []byte {
	retval := make([]byte, len(mySlice)*2)
	for i, v := range mySlice {
		binary.LittleEndian.PutUint16(retval[i*2:], uint16(float32toint16(v)))
	}
	return retval
}

func float32toint16(num float32) int16 {

	return int16(math.Max(1-(1<<15), (math.Min((1<<15)-1, float64(num)*(1<<16)))))
}
