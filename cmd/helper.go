package cmd

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

// TODO: float?
func getFileSizeBytes(path string) (int64, error) {
	fi, err := os.Stat(filePath)
	check(err)

	return fi.Size(), nil
}

// From https://pkg.go.dev/internal/bytealg#Count
func countByteOccurence(b []byte, c byte) (n int64) {
	for _, x := range b {
		if x == c {
			n++
		}
	}
	return n
}

// Based on https://stackoverflow.com/a/24563853/9283726
func lineCounter(r io.Reader) (count int64, err error) {
	buf := make([]byte, bufio.MaxScanTokenSize)

	lineBreak := '\n'

	for {
		c, err := r.Read(buf)
		count += countByteOccurence(buf[:c], byte(lineBreak))

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}

}

func wordCounter(r *bufio.Reader) (count int64, err error) {
	prevRune := ' '

	for {
		r, _, err := r.ReadRune()
		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}

		if !unicode.IsSpace(r) && unicode.IsSpace(prevRune) {
			count++
		}

		prevRune = r

	}
}
