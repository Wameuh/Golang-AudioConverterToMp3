package OggToMp3

import (
	"errors"
	"fmt"
	"io"
	"math"
	"unsafe"

	"bytes"

	"github.com/jfreymuth/oggvorbis"
	"github.com/viert/go-lame"
)

func float32toint16(num float32) int16 {

	return int16(math.Max(1-math.Pow(2, 15), (math.Min(math.Pow(2, 15)-1, float64(num)*math.Pow(2, 16)))))
}

func convertSliceToInt16Slice(mySlice []float32) []int16 {
	retval := make([]int16, 0)
	for _, v := range mySlice {
		retval = append(retval, float32toint16(v))
	}
	return retval
}

func GetByteSlice(r io.Reader) ([]byte, *oggvorbis.Format, error) {
	var format *oggvorbis.Format
	format, err := oggvorbis.GetFormat(r)
	if err != nil {
		return nil, nil, err
	}

	oggAudio, format, err := oggvorbis.ReadAll(r)
	if err != nil {
		return nil, nil, errors.New("Could not decode ogg file?")
	}

	audioBytes := convertSliceToInt16Slice(oggAudio)

	return *(*[]byte)(unsafe.Pointer(&audioBytes)), format, nil
}

func EncodeMP3Slice(input []byte, format *oggvorbis.Format) []byte {
	output := new(bytes.Buffer)
	enc := lame.NewEncoder(output)

	defer enc.Close()

	if format.Channels == 1 {
		enc.SetNumChannels(1)
	}

	fmt.Println("sampleRate: ", format.SampleRate)

	enc.SetVBR(lame.VBRDefault)
	enc.SetInSamplerate(format.SampleRate)
	enc.SetVBRQuality(4)
	enc.SetQuality(4)
	enc.SetMode(lame.MpegMono)
	enc.SetWriteID3TagAutomatic(false)
	enc.Write(input)
	enc.Flush()
	fmt.Println("lameSlice length: ", len(output.Bytes()))

	return output.Bytes()
}
