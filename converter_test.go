package GoAudioConverterToMp3

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/tcolgate/mp3"
)

func getDuration(conv *Converter) float64 {
	file, err := os.Create(conv.inputFile.Name() + "tempfile.mp3")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//defer os.Remove(file.Name())
	_, err = file.Write(conv.GetDataConverted())
	if err != nil {
		fmt.Println(err)
	}
	file.Sync()
	file.Close()

	file, err = os.Open("tempfile.mp3")
	if err != nil {
		fmt.Println(err)
	}

	t := 0.0
	d := mp3.NewDecoder(file)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return 0
		}

		t = t + f.Duration().Seconds()
	}
	return t
}

func TestNewConverter(t *testing.T) {
	// init

	file, err := os.Open("assets_for_testing/file_example_OOG_2MG.ogg")
	if err != nil {
		t.Fatalf(`error while opening file_example_OOG_2MG.ogg`)
	}
	defer file.Close()

	// proceed
	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	// test default values
	if conv.toDiscard != 0 {
		t.Fatalf(`error default value toDiscard. Is %d, want 0`, conv.toDiscard)
	}
	if conv.outputQuality != 4 {
		t.Fatalf(`error default value outputQuality. Is %d, want 4`, conv.outputQuality)
	}
	if conv.outputMode != -1 {
		t.Fatalf(`error default value outputMode. Is %d, want -1`, conv.outputMode)
	}
	if conv.outputSampleRate != 44100 {
		t.Fatalf(`error default value outputSampleRate. Is %d, want 44100`, conv.outputSampleRate)
	}

	// test stats
	if conv.inputFileStat.Name() != "file_example_OOG_2MG.ogg" {
		t.Fatalf(`error name of file input. Is %s, want file_example_OOG_2MG.ogg`, conv.inputFileStat.Name())
	}
	if conv.inputFileStat.Size() != 2076666 {
		t.Fatalf(`error size of file input. Is %d, want 2 076 666`, conv.inputFileStat.Size())
	}
}

func TestSetParameters(t *testing.T) {
	// init

	file, err := os.Open("assets_for_testing/file_example_OOG_2MG.ogg")
	if err != nil {
		t.Fatalf(`error while opening file_example_OOG_2MG.ogg`)
	}
	defer file.Close()

	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	// proceed
	conv.SetQuality(3)
	conv.SetNumsChannels(8)
	err1 := conv.SetOutSampleRate(55665)
	err2 := conv.SetOutSampleRate(11025)

	// test

	if conv.outputQuality != 3 {
		t.Fatalf(`error SetQuality(). Is %d, want 3`, conv.outputQuality)
	}
	if conv.outputMode != 8 {
		t.Fatalf(`error SetNumsChannels(). Is %d, want 8`, conv.outputMode)
	}
	if err1 == nil {
		t.Fatalf(`error SetOutSampleRate(). err is nil, want different from nil`)
	}
	if err2 != nil {
		t.Fatalf(`error SetOutSampleRate(). err is %s, want nil`, err1.Error())
	}
	if conv.outputSampleRate != 11025 {
		t.Fatalf(`error SetOutSampleRate(). outputSampleRate is %d, want 11025`, conv.outputSampleRate)
	}
}

func TestSetConverterFormat(t *testing.T) {
	// init

	file, err := os.Open("assets_for_testing/file_example_OOG_2MG.ogg")
	if err != nil {
		t.Fatalf(`error while opening file_example_OOG_2MG.ogg`)
	}
	defer file.Close()

	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	conv.outputQuality = 3
	conv.outputMode = 3
	conv.outputSampleRate = 12000
	conv.inputNumChan = 1
	conv.inputSampleRate = 8000

	// proceed
	err = conv.SetConverterFormat()

	//test
	if err != nil {
		t.Fatalf(`error SetConverterFormat(), err: %s`, err.Error())
	}
	if conv.encoder.Quality() != 3 {
		t.Fatalf(`error encoder.Quality(), is %d, want 3`, conv.encoder.Quality())
	}
	if conv.encoder.InSamplerate() != 8000 {
		t.Fatalf(`error encoder.InSamplerate(), is %d, want 8000`, conv.encoder.InSamplerate())
	}
	if conv.encoder.NumChannels() != 1 {
		t.Fatalf(`error encoder.NumChannels() (input), is %d, want 1`, conv.encoder.NumChannels())
	}
}

func TestWavGetFormat(t *testing.T) {
	// init
	file, err := os.Open("assets_for_testing/soundForTest.wav")
	if err != nil {
		t.Fatalf(`error while opening soundForTest.wav`)
	}
	defer file.Close()

	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	// proceed
	err = conv.wavGetFormat()

	// test
	if err != nil {
		t.Fatalf(`error wavGetFormat. err is %s, want nil`, err.Error())
	}
	if conv.toDiscard == 0 {
		t.Fatalf(`error wavGetFormat. toDiscard is still 0, expected >0`)
	}
	if conv.inputNumChan != 1 {
		t.Fatalf(`error wavGetFormat. inputNumChan is %d, expected 1`, conv.inputNumChan)
	}
	if conv.inputSampleRate != 22050 {
		t.Fatalf(`error wavGetFormat. inputSampleRate is %d, expected 22050`, conv.inputSampleRate)
	}

}

func TestFromOggToMp3(t *testing.T) {
	// init
	file, err := os.Open("assets_for_testing/file_example_OOG_2MG.ogg")
	if err != nil {
		t.Fatalf(`error while opening file_example_OOG_2MG.ogg`)
	}
	defer file.Close()

	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	// proceed
	n, err := conv.FromOggToMp3()

	//
	if err != nil {
		t.Fatalf(`error FromOggToMp3(), %s`, err.Error())
	}
	if n < 74*44100*2*2 {
		t.Fatalf(`error FromOggToMp3(), amount of data is incorrect, is %d, expected more than %d`, n, 74*44100*2*2)
	}
	duration := getDuration(conv)
	if duration > 75 || duration < 74 {

		t.Fatalf(`error FromOggToMp3(), duration of mp3 is incorrect, is %f, expected around %f`, duration, 74.42)
	}
}

func TestFromWavToMp3(t *testing.T) {
	// init
	file, err := os.Open("assets_for_testing/soundForTest.wav")
	if err != nil {
		t.Fatalf(`error while opening soundForTest.wav`)
	}
	defer file.Close()

	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	// proceed
	n, err := conv.FromWavToMp3()
	duration := getDuration(conv)

	// test
	if err != nil {
		t.Fatalf(`error FromOggToMp3(), %s`, err.Error())
	}
	if n != conv.inputFileStat.Size()-int64(conv.toDiscard) {
		t.Fatalf(`error FromOggToMp3(), audio source size incorrect, is %d, expected %d`, n, conv.inputFileStat.Size()-int64(conv.toDiscard))
	}
	if duration > 61 || duration < 60 {

		t.Fatalf(`error FromOggToMp3(), duration of mp3 is incorrect, is %f, expected around %f`, duration, 60.0555)
	}
}

func TestFromMp3toMp3(t *testing.T) {
	// init
	file, err := os.Open("assets_for_testing/file_example_MP3_1MG_mod.mp3")
	if err != nil {
		t.Fatalf(`error while opening file_example_MP3_1MG_mod.mp3`)
	}
	defer file.Close()

	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	// proceed
	n, err := conv.FromMp3ToMp3()
	duration := getDuration(conv)

	// test
	if err != nil {
		t.Fatalf(`error FromMp3ToMp3(), %s`, err.Error())
	}
	if n < 26*16000*2 {
		t.Fatalf(`error FromMp3ToMp3(), n is %d, expected more than %d`, n, 26*16000*2)
	}
	if duration > 27 || duration < 26 {

		t.Fatalf(`error FromOggToMp3(), duration of mp3 is incorrect, is %f, expected around %f`, duration, 26.244000)
	}

}

func TestConvertToMp3(t *testing.T) {
	//init

	file, err := os.Open("assets_for_testing/file_example_MP3_1MG_mod.mp3")
	if err != nil {
		t.Fatalf(`error while opening file_example_MP3_1MG_mod.mp3`)
	}
	defer file.Close()

	conv, err := NewConverter(file)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv.Close()

	file2, err := os.Open("assets_for_testing/filewithoutextension")
	if err != nil {
		t.Fatalf(`error while opening filewithoutextension`)
	}
	defer file2.Close()

	conv2, err := NewConverter(file2)
	if err != nil {
		t.Fatalf(`error while creating new converter`)
	}
	defer conv2.Close()

	// proceed

	_, err = conv.ConvertToMp3()
	_, err2 := conv2.ConvertToMp3()

	// test

	if err != nil {
		t.Fatalf(`error ConvertToMp3(): %s`, err.Error())
	}

	if err2 == nil {
		t.Fatalf(`no error in ConvertToMp3(). Expected one.`)
	}

}
