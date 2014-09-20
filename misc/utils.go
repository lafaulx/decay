package misc

import (
	crand "crypto/rand"
	"io"
	"math"
	"math/rand"
	"os"
)

func RandIntFromInterval(min int, max int) int {
	return int(math.Floor(rand.Float64()*float64((max-min+1)) + float64(min)))
}

func RandomCharFromString(s string) string {
	charPosition := rand.Intn(len(s))

	return s[charPosition : charPosition+1]
}

func ReplaceCharInString(s string, char string, index int) string {
	result := s[0:index] + char + s[index+1:]

	return result
}

func FortuneWheel(cnt int64) bool {
	return rand.Float64() < (math.Log10(float64(cnt+1))/math.Log10(100))/10.0
}

func GetDataCountToDamage(length int) int {
	length64 := float64(length)
	return int(math.Max((math.Log10(length64) + math.Sqrt(length64)/5), 1))
}

func GetDataCountToDamageImage(length int) int {
	return RandIntFromInterval(0, int(float32(length)*0.2))
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func GenerateID() int {
	n, _ := crand.Prime(crand.Reader, 30)

	return int(n.Uint64())
}
