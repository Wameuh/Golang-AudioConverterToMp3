package GoAudioConverterToMp3

//#cgo LDFLAGS: -lmp3lame
//#include <lame/lame.h>
import "C"

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	"github.com/go-audio/wav"
	"github.com/tosone/minimp3"
	"github.com/viert/go-lame"
)

type Converter struct {
	encoder       *lame.Encoder
	inputFile     *os.File
	inputFileStat os.FileInfo
	output        *bytes.Buffer // the buffer where mp3 data are stored

	outputMode       lame.MpegMode // default: lame picks based on compression ration and input channels
	outputQuality    int           // used for SetQuality in lame. From 0 to 9 - default is 4
	outputSampleRate uint32        //default is 44100

	inputNumChan    int
	inputSampleRate int
	toDiscard       int // used for wav format to discard the header

}

func NewConverter(inputFile *os.File) (*Converter, error) {
	this := new(Converter)
	err := this.init(inputFile)
	return this, err
}

func (c *Converter) init(inputFile *os.File) error {
	c.output = new(bytes.Buffer)
	c.encoder = lame.NewEncoder(c.output)
	c.inputFile = inputFile

	var err error
	c.inputFileStat, err = inputFile.Stat() //get info about the input file
	if err != nil {
		return err
	}

	c.toDiscard = 0
	c.outputQuality = 4
	c.outputMode = -1
	c.outputSampleRate = 44100
	return nil
}

func (c *Converter) Close() {
	c.encoder.Close()
}

func (c *Converter) SetQuality(quality int) {
	c.outputQuality = quality
}

func (c *Converter) SetNumsChannels(mode lame.MpegMode) {
	c.outputMode = mode
}

func (c *Converter) SetOutSampleRate(sampleRate uint32) error {
	if sampleRate == 8000 || sampleRate == 11025 || sampleRate == 12000 || sampleRate == 16000 || sampleRate == 22050 || sampleRate == 24000 || sampleRate == 32000 || sampleRate == 44100 || sampleRate == 48000 {
		c.outputSampleRate = sampleRate
		return nil
	} else {
		return errors.New("Wrong samplerate format for the output")
	}
}

func (c *Converter) SetConverterFormat() error {
	err := c.encoder.SetQuality(c.outputQuality)
	if err != nil {
		return err
	}
	err = c.encoder.SetVBR(lame.VBROff) //To be in CBR Mode
	if err != nil {
		return err
	}
	err = c.encoder.SetInSamplerate(c.inputSampleRate)
	if err != nil {
		return err
	}
	err = c.encoder.SetNumChannels(c.inputNumChan)
	if err != nil {
		return err
	}

	res := int(C.lame_set_out_samplerate(C.lame_init(), C.int(c.outputSampleRate)))
	if res < 0 {
		return lame.Error(res)
	}

	if c.outputMode != -1 {
		err = c.encoder.SetMode(c.outputMode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Converter) wavGetFormat() error {
	decod := wav.NewDecoder(c.inputFile)
	decod.ReadInfo()
	c.inputNumChan = int(decod.NumChans)
	c.inputSampleRate = int(decod.SampleRate)

	if decod.BitDepth != 16 {
		return errors.New("Bit depth of the input file is different from 16. This format is not supported")
	}

	err := decod.FwdToPCM()
	if err != nil {
		return err
	}
	PCMSize := decod.PCMLen()
	c.toDiscard = (int(c.inputFileStat.Size() - PCMSize)) // The difference between the size of the file and the size of PCM data

	return nil
}

func (c *Converter) FromOggToMp3() (int64, error) { // use oggvorbis
	audioBytes, format, err := GetByteSlice(c.inputFile)
	if err != nil {
		return 0, err
	}

	c.inputNumChan = format.Channels
	c.inputSampleRate = format.SampleRate

	err = c.SetConverterFormat()
	if err != nil {
		return 0, err
	}

	n, err := c.encoder.Write(audioBytes)
	c.encoder.Flush()

	return int64(n), err

}

func (c *Converter) FromWavToMp3() (n int64, err error) { //use go-audio/wav
	r := bufio.NewReader(c.inputFile)

	err = c.wavGetFormat()
	if err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	//Discard the header of the Wav
	_, err = r.Discard(c.toDiscard)
	if err != nil {
		return 0, err
	}
	n, err = r.WriteTo(c.encoder)
	c.encoder.Flush()

	return n, err
}

func (c *Converter) FromMp3ToMp3() (int64, error) { //use minimp3
	fileBytes, err := ioutil.ReadAll(c.inputFile)
	if err != nil {
		return 0, err
	}

	data, mp3Audio, err := minimp3.DecodeFull(fileBytes)
	if err != nil {
		return 0, err
	}
	c.inputNumChan = data.Channels
	c.inputSampleRate = data.SampleRate

	err = c.SetConverterFormat()
	if err != nil {
		return 0, err
	}

	if data.Channels == int(c.outputMode) && data.SampleRate == int(c.outputSampleRate) {
		//No conversion needed
		n, err := c.output.Write(fileBytes)
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	}

	n, err := c.encoder.Write(mp3Audio)
	if err != nil {
		return 0, err
	}
	c.encoder.Flush()

	return int64(n), nil

}

func (c *Converter) ConvertToMp3() (int64, error) {
	//check extension by getting last 3 char from the name of the file
	extension := c.inputFileStat.Name()[len(c.inputFileStat.Name())-3:]

	switch extension {
	case "mp3":
		return c.FromMp3ToMp3()
	case "ogg":
		return c.FromOggToMp3()
	case "wav":
		return c.FromWavToMp3()
	default:
		return 0, errors.New("Error, format not supported")
	}

}

func (c *Converter) WriteBufferToFile(filepath string) error {

	err := os.WriteFile(filepath, c.output.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}
