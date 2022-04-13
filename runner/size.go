package runner

import (
	"fmt"
	"strconv"
)

type Size uint64

// Print Size in string by unit
func (s Size) String() string {
	t := uint64(s)
	switch {
	case t < 1<<10:
		return fmt.Sprintf("%d B", t)
	case t < 1<<20:
		return fmt.Sprintf("%.1f KiB", float64(t)/float64(1<<10))
	case t < 1<<30:
		return fmt.Sprintf("%.1f MiB", float64(t)/float64(1<<20))
	default:
		return fmt.Sprintf("%.1f GiB", float64(t)/float64(1<<30))
	}
}

// Set size by string
func (s *Size) Set(str string) error {
	switch str[len(str)-1] {
	case 'b', 'B':
		str = str[:len(str)-1]
	}

	factor := 0
	switch str[len(str)-1] {
	case 'k', 'K':
		factor = 10
		str = str[:len(str)-1]
	case 'm', 'M':
		factor = 20
		str = str[:len(str)-1]
	case 'g', 'G':
		factor = 30
		str = str[:len(str)-1]
	}

	t, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	*s = Size(t << factor)
	return nil
}

// Return B by Size
func (s Size) Byte() uint64 {
	return uint64(s)
}

// Return KB by Size
func (s Size) KiB() uint64 {
	return uint64(s) >> 10
}

// Return MB by Size
func (s Size) MiB() uint64 {
	return uint64(s) >> 20
}

// Return GB by Size
func (s Size) GiB() uint64 {
	return uint64(s) >> 30
}
