package GoAudioConverterToMp3

import (
	"math"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

func TestFloat32toint16(t *testing.T) {
	var num float32
	num = -2.0
	want := int16(1 - math.Pow(2, 15))
	result := float32toint16(num)
	if want != result {
		t.Fatalf(`float32toint16(-2.0) = %d, want match for %d,`, result, want)
	}

	num = 0.0
	want = int16(0)
	result = float32toint16(num)
	if want != result {
		t.Fatalf(`float32toint16(0.0) = %d, want match for %d,`, result, want)
	}

	num = 1.0
	want = int16(math.Pow(2, 15) - 1)
	result = float32toint16(num)
	if want != result {
		t.Fatalf(`float32toint16(1.0) = %d, want match for %d,`, result, want)
	}
}

func TestConvertSliceToInt16Slice(t *testing.T) {
	// init
	lenOfSlice := rand.Intn(1000)
	floatSlice := make([]float32, lenOfSlice)
	for i := range floatSlice {
		floatSlice[i] = rand.Float32()*2 - 1
	}
	emptyFloatSlice := make([]float32, 0)

	// proceed
	intSlice := convertSliceToInt16Slice(floatSlice)
	emptyIntSlice := convertSliceToInt16Slice(emptyFloatSlice)

	// test
	if len(intSlice) != lenOfSlice*2 {
		t.Fatalf(`convertSliceToInt16Slice(), len of %d, want %d`, len(intSlice), lenOfSlice*2)
	}
	if len(emptyIntSlice) != 0 {
		t.Fatalf(`convertSliceToInt16Slice(), len of %d, want %d`, len(emptyIntSlice), 0)
	}
	for _, v := range intSlice {
		if reflect.TypeOf(v).String() != "uint8" {
			t.Fatalf(`convertSliceToInt16Slice(), %d type of %s, want type int16`, v, reflect.TypeOf(v).String())
		}
	}

}

func TestGetByteSlice(t *testing.T) {

	// init
	file, err := os.Open("assets_for_testing/file_example_OOG_2MG.ogg")
	if err != nil {
		t.Fatalf(`error while opening file_example_OOG_2MG.ogg`)
	}
	defer file.Close()

	expectedNumberOfChannel := 2
	expectedSampleRate := 44100
	durationInSecond := 74
	expectedSizeOfDataInBytes := expectedSampleRate * expectedNumberOfChannel * durationInSecond * 2 // 2 bytes per channel
	//as duration is not really accurate, in the test we should expect +/- 1 second

	// proceed
	data, format, err := GetByteSlice(file)

	// test
	if format.Channels != expectedNumberOfChannel {
		t.Fatalf(`GetByteSlice(). Format.Channels is %d, want %d`, format.Channels, expectedNumberOfChannel)
	}
	if format.SampleRate != expectedSampleRate {
		t.Fatalf(`GetByteSlice(). Format.SampleRate is %d, want %d`, format.SampleRate, expectedSampleRate)
	}
	if len(data) < expectedSizeOfDataInBytes-44100*2 {
		t.Fatalf(`GetByteSlice(). len of data is %d, want more than %d`, len(data), expectedSizeOfDataInBytes-44100*2)
	}
	if len(data) > expectedSizeOfDataInBytes+44100*2 {
		t.Fatalf(`GetByteSlice(). len of data is %d, want less than %d`, len(data), expectedSizeOfDataInBytes+44100*2)
	}
}
