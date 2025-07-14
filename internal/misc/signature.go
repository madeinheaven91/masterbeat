package internal

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TimeSignature struct {
	Top    int
	Bottom int
}

func NewTimeSignature(numerator int, denominator int) (*TimeSignature, error) {
	if math.Log2(float64(denominator)) != float64(int(math.Log2(float64(denominator)))) {
		return nil, fmt.Errorf("denominator should be a power of 2")
	}
	return &TimeSignature{Top: numerator, Bottom: denominator}, nil
}

func ParseTimeSignature(input string) (*TimeSignature, error) {
	tmp := strings.Split(input, "/")
	if len(tmp) != 2 {
		return nil, fmt.Errorf("invalid time signature")
	}
	numerator, err := strconv.Atoi(tmp[0])
	if err != nil {
		return nil, fmt.Errorf("invalid time signature")
	}
	denominator, err := strconv.Atoi(tmp[1])
	if err != nil {
		return nil, fmt.Errorf("invalid time signature")
	}
	sig, err := NewTimeSignature(numerator, denominator)
	return sig, err
}

func (t TimeSignature) String() string {
	return fmt.Sprintf("%d/%d", t.Top, t.Bottom)
}
