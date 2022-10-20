# Audio converter to Mp3
Basic package to convert an audio files to mp3.


# Supported formats
Currently only these formats are supported:
* Wav (PCM16)
* Ogg
* Mp3

# Dependencies
This package use differents go-packages (thanks for them):
* [lame-go](https://github.com/viert/go-lame)
* [oggVorbis](https://github.com/jfreymuth/oggvorbis)
* [go-audio/wav](https://github.com/go-audio/wav)
* [minimp3](https://github.com/tosone/minimp3)

# Example of implementation

```go
func main() {
    //open file
	file, err := os.Open("inputFile.ext")
	defer file.Close()
	if err != nil {
		fmt.Println("error opening input file: ", err)
		return
	}

	// Create converter
	converter, err := NewConverter(file)
	if err != nil {
		fmt.Println("error creating converter: ", err)
		return
	}

	//Convert data
	_, err := converter.ConvertToMp3()
	if err != nil {
		fmt.Println("error in converting: ", err)
		return
	}

    //write into a new file
	err = converter.WriteFile("myOutputFile.mp3")
	if err != nil {
		fmt.Println("error in writing file: ", err)
		return
	}

	fmt.Println("Job is done!")
	return
}
```

# TODO
* Testings
* Support other bit depth than 16 for wav