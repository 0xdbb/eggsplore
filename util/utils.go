package util

import (
	cr "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	PRECISION              = 12
	EPSILON                = 1.0e-12
	Pi                     = math.Pi
	Tau                    = 2.0 * Pi
	Sqrt2                  = math.Sqrt2
	E                      = math.E
	MaxFloat64             = math.MaxFloat64
	SmallestNonzeroFloat64 = math.SmallestNonzeroFloat64
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// OTPExpirationDuration defines how long an OTP is valid (5 minutes).
var OTPExpirationDuration = 5 * time.Minute

// GenerateSecureToken generates a secure random token of specified length
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := cr.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func FloatEqual(a, b float64) bool {
	return a == b || Abs(a-b) < EPSILON
}

func ExcludeValue(values []float64, exclude float64) []float64 {
	res := make([]float64, 0, len(values))
	for _, v := range values {
		if !FloatEqual(v, exclude) {
			res = append(res, v)
		}
	}
	return res
}

func ExcludeNaN(values []float64) []float64 {
	res := make([]float64, 0, len(values))
	for _, v := range values {
		if !math.IsNaN(v) {
			res = append(res, v)
		}
	}
	return res
}

func MinValue(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}
	min := values[0]
	for _, v := range values {
		min = math.Min(min, v)
	}
	return min
}

func AsF64(value string) float64 {
	v := math.NaN()
	var err error
	if len(value) > 0 {
		v, err = strconv.ParseFloat(value, 64)
		CheckError(err)
	}
	return v
}

func Decode64(s string) []byte {
	dat, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return dat
}

func Encode64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// RoundFloor rounds a float to the nearest whole number float
func RoundFloor(f float64) float64 {
	return math.Trunc(f + math.Copysign(0.5, f))
}

// Round rounds a number to the nearest decimal place
func Round(x float64, digits ...int) float64 {
	d := 0
	if len(digits) > d {
		d = digits[d]
	}
	m := math.Pow(10.0, float64(d))
	return RoundFloor(x*m) / m
}

func ConvertToByteArrayStr(b []byte) string {
	tokens := make([]string, 0, len(b))
	for _, o := range b {
		tokens = append(tokens, strconv.Itoa(int(o)))
	}
	return fmt.Sprintf("[%v]", strings.Join(tokens, ","))
}

func ImageToBase64(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	var base64Encoding string
	mimeType := http.DetectContentType(bytes)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += Encode64(bytes)

	return base64Encoding, nil
}
