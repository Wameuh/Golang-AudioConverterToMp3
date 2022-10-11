# Golang-OggToMp3
Basic package to decode Ogg and code it to Mp3.

It uses [lame-go](https://github.com/viert/go-lame) and [oggVorbis](https://github.com/jfreymuth/oggvorbis).

# Example of implementation

```go
package OggToMp3_example

import (
	OggToMp3 "OggToMp3"
	"fmt"
	"os"
)

func main() {
    //open file
	file, err := os.Open("yourOggFile")
	defer file.Close()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

    //decode file and get a byteslice (int16)
	audioSlice, format, err := OggToMp3.GetByteSlice(file)
	if err != nil {
		fmt.Println("error in GetByteSlice: ", err)
		return
	}

    //path the byte slice to lame
	mp3AudioSlice := OggToMp3.EncodeMP3Slice(audioSlice, format)

    //write into a new file
	err = os.WriteFile("FileName", mp3AudioSlice, 0777)
	if err != nil {
		fmt.Println("error in WriteFile: ", err)
		return
	}

}```
