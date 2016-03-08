package pbar

import (
	"fmt"
)

type ProgressBar struct {
	count float32
	size  float32
}

func New(size float32) *ProgressBar {
	return &ProgressBar{0, size}
}

func (pb *ProgressBar) Advance(count float32) {
	if pb.count+count > pb.size {
		pb.count = pb.size
	} else {
		pb.count += count
	}
}

func (pb *ProgressBar) Percent() float32 {
	return (pb.count / pb.size) * 100
}

func (pb *ProgressBar) Print(msg string) {
	fmt.Print("\r[")
	for i := 0; i < 50; i++ {
		percent := (i + 1) * 2
		if pb.Percent() >= float32(percent) {
			fmt.Print("#")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("] [%.2f%%] %s", pb.Percent(), msg)
}
