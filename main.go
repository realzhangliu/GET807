package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

const (
	//DICTVOICE source link
	DICTVOICE = "http://dict.youdao.com/dictvoice?audio=%s&type=1"
)

//DownloadMP3 get mp3 files
func DownloadMP3(word string) {
	client := http.Client{}
	targetURL := fmt.Sprintf(DICTVOICE, word)
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile(path.Join("data", word+".mp3"), body, 0755)
}

//LaunchDownload  s
func LaunchDownload() {
	os.MkdirAll("data", os.ModePerm)
	f, err := os.Open("807.txt")
	if err != nil {
		fmt.Println(err)
	}
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		DownloadMP3(scan.Text())
	}
}

func main() {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		fmt.Println("Starting download audio data.")
		LaunchDownload()
	}
	var words []string
	fs, _ := ioutil.ReadDir("data")
	fmt.Println(len(fs))
	for i := 0; i < len(fs); i++ {
		words = append(words, fs[i].Name())
	}

	//listening & type correct word.
	input := bufio.NewScanner(os.Stdin)

	rd := rand.Intn(808)

	PlayWord(path.Join("data", words[rd]))

	for input.Scan() {
		rd = rand.Intn(808)

		PlayWord(path.Join("data", words[rd]))

		if input.Text() == words[rd] {
			fmt.Printf("correct :%s\n", words[rd])
		}

		if input.Text() == "q" {
			fmt.Printf("WORD:%s\n", words[rd])
		}

	}

}
func PlayWord(fpath string) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
