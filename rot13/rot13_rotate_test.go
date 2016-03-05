package main

import "testing"

func TestRotateDowncaseCh(t *testing.T) {
	for ch := 'a'; ch <= 'z'; ch++ {
		charIndex := ch - 'a'
		rCh := (charIndex+13)%26 + 'a'
		haveCh := rotateCh(ch)
		if haveCh != rCh {
			t.Errorf("%d should be converted to %d, but had %d", ch, rCh, haveCh)
		}
	}
}

func TestRotateUpcaseCh(t *testing.T) {
	for ch := 'A'; ch <= 'Z'; ch++ {
		charIndex := ch - 'A'
		rCh := (charIndex+13)%26 + 'A'
		haveCh := rotateCh(ch)
		if haveCh != rCh {
			t.Errorf("%c should be converted to %c, but had %c", ch, rCh, haveCh)
		}
	}
}

func TestRotateTwice(t *testing.T) {
	for ch := 'a'; ch <= 'Z'; ch++ {
		haveCh := rotateCh(rotateCh(ch))
		if haveCh != ch {
			t.Errorf("%c should not change after been rotated twice")
		}
	}
}
