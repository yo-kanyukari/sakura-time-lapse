package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func unpackTar() {
	var file *os.File
	var err error

	//tarのopen
	if file, err = os.Open("test/200506-125.tar"); err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	//tar fileにまとめられているfileの先頭のfileを受け取る
	reader := tar.NewReader(file)

	//一つずつfileを作っていく
	for i := 0; true; i++ {
		_, err = reader.Next()
		if err == io.EOF {
			// ファイルの最後
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, reader); err != nil {
			log.Fatalln(err)
		}

		if err = ioutil.WriteFile(fmt.Sprintf("test/source%04d.jpg", i), buf.Bytes(), 0755); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("unpack tar")
}

func makeTimeLapse() {
	err := exec.Command("ffmpeg", "-f", "image2", "-r", "3", "-i", "test/source%04d.jpg", "-r", "15", "-an", "-vcodec", "libx264", "-pix_fmt", "yuv420p", "test/video.mp4","-y").Run()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("make time-lapse")
}

func main() {
	unpackTar()
	makeTimeLapse()
}
