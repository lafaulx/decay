package dcy

import (
	"bytes"
	"encoding/base64"
	"github.com/lafaulx/decay/misc"
	"github.com/lafaulx/decay/model"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
)

var imageCache map[int]string = make(map[int]string)
var imageCallsCache map[int]int = make(map[int]int)

type Image struct{}

func (k *Image) getFile(t *model.Dcy) []byte {
	var file []byte
	cachedFileStr, exists := imageCache[t.Id]

	if exists {
		file, _ = base64.StdEncoding.DecodeString(cachedFileStr)
	} else {
		file, _ = ioutil.ReadFile("images/" + t.Content)
	}

	return file
}

func (k *Image) getFileBase64(t *model.Dcy) string {
	var file []byte
	cachedFileStr, exists := imageCache[t.Id]

	if exists {
		return cachedFileStr
	} else {
		file, _ = ioutil.ReadFile("images/" + t.Content)
		return k.convertFileToBase64(file)
	}
}

func (k *Image) processFile(file []byte) []byte {
	fileReader := bytes.NewReader(file)

	predamagedImage, _ := jpeg.Decode(fileReader)
	x := predamagedImage.Bounds().Size().X
	y := predamagedImage.Bounds().Size().Y

	damagedImage := image.NewRGBA(predamagedImage.Bounds())
	draw.Draw(damagedImage, damagedImage.Bounds(), predamagedImage, image.Point{0, 0}, draw.Src)

	xAt := misc.RandIntFromInterval(0, x)
	yAt := misc.RandIntFromInterval(0, y)
	xPos := misc.RandIntFromInterval(0, x)
	yPos := misc.RandIntFromInterval(0, y)
	xLim := misc.GetDataCountToDamageImage(x)
	yLim := misc.GetDataCountToDamageImage(y)

	sr := image.Rectangle{image.Point{xPos, yPos}, image.Point{xPos + xLim, yPos + yLim}}
	dp := image.Point{xAt, yAt}

	r := image.Rectangle{dp, dp.Add(sr.Size())}
	draw.Draw(damagedImage, r, damagedImage, sr.Min, draw.Src)

	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, damagedImage, &jpeg.Options{misc.RandIntFromInterval(0, 100)})
	return buffer.Bytes()
}

func (k *Image) convertFileToBase64(file []byte) string {
	return base64.StdEncoding.EncodeToString(file)
}

func (k *Image) cleanup(t *model.Dcy, damagedFile []byte, damagedFileBase64 string) {
	imageCache[t.Id] = damagedFileBase64
	_, exists := imageCallsCache[t.Id]

	if !exists {
		imageCallsCache[t.Id] = 0
	}

	imageCallsCache[t.Id] += 1

	if imageCallsCache[t.Id] > 3 {
		imageCallsCache[t.Id] = 0
		ioutil.WriteFile("images/"+t.Content, damagedFile, 0644)
	}
}

func (k *Image) Damage(t *model.Dcy) *model.Dcy {
	if misc.FortuneWheel(t.CallCount) {
		damagedFile := k.processFile(k.getFile(t))

		damagedFileBase64 := k.convertFileToBase64(damagedFile)
		k.cleanup(t, damagedFile, damagedFileBase64)

		t.Content = "data:image/jpeg;base64," + damagedFileBase64
	} else {
		fileBase64 := k.getFileBase64(t)
		imageCache[t.Id] = fileBase64
		t.Content = "data:image/jpeg;base64," + fileBase64
	}

	return t
}
