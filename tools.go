package GoAudioConverterToMp3

import (
	"io"
	"math"
	"unsafe"

	"github.com/jfreymuth/oggvorbis"
)

func GetByteSlice(r io.Reader) ([]byte, *oggvorbis.Format, error) {
	oggAudio, format, err := oggvorbis.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	//oggvorbis give []float32 slice, we need to convert it to []in16
	audioBytes := convertSliceToInt16Slice(oggAudio)

	return *(*[]byte)(unsafe.Pointer(&audioBytes)), format, nil
}

func convertSliceToInt16Slice(mySlice []float32) []int16 {
	retval := make([]int16, 0)
	for _, v := range mySlice {
		retval = append(retval, float32toint16(v))
	}
	return retval
}

func float32toint16(num float32) int16 {

	return int16(math.Max(1-math.Pow(2, 15), (math.Min(math.Pow(2, 15)-1, float64(num)*math.Pow(2, 16)))))
}
