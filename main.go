package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"

	"bytes"

	"github.com/jfreymuth/oggvorbis"
	"github.com/viert/go-lame"
)

func main() {
	file, err := os.Open("assets/test.ogg")
	defer file.Close()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	audioSlice, format, err := GetByteSlice(file)
	if err != nil {
		fmt.Println("error in GetByteSlice: ", err)
		return
	}

	fmt.Println("sampleRate: ", format.SampleRate)

	mp3AudioSlice := EncodeMP3Slice(audioSlice, format)

	err = os.WriteFile("test.mp3", mp3AudioSlice, 0777)
	if err != nil {
		fmt.Println("error in WriteFile: ", err)
		return
	}

	fmt.Println("program end")

}

func float32tobyte(num float32) byte {

	return byte(math.Max(0, (math.Min(255, float64(num)*math.Pow(2, 7)+math.Pow(2, 7)))))
}

func convertSlice(mySlice []float32) []byte {
	retval := make([]byte, 0)
	for _, v := range mySlice {
		retval = append(retval, float32tobyte(v))
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
		fmt.Println("erreur reading oggfile", err)
		return nil, nil, errors.New("Could not decode ogg file?")
	}
	fmt.Println("oggAudio length: ", len(oggAudio))
	audioBytes := convertSlice(oggAudio)
	fmt.Println("audioBytes length: ", len(audioBytes))
	return audioBytes, format, nil
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
