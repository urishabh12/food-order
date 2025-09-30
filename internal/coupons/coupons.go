package coupons

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Store struct {
	filesCount map[string]int
}

func Load(dir string, files ...string) (*Store, error) {
	if len(files) == 0 {
		files = []string{"couponbase1", "couponbase2", "couponbase3"}
	}
	counts := map[string]int{}

	for _, name := range files {
		path := filepath.Join(dir, name)
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", path, err)
		}

		defer f.Close()

		seenInThisFile := map[string]struct{}{}
		sc := bufio.NewScanner(f)

		for sc.Scan() {
			code := strings.ToUpper(strings.TrimSpace(sc.Text()))
			if code == "" {
				continue
			}

			if _, ok := seenInThisFile[code]; ok {
				continue
			}
			seenInThisFile[code] = struct{}{}
		}

		for code := range seenInThisFile {
			counts[code]++
			if counts[code] > 1 {
				fmt.Printf("%s Valid coupon code", code)
			}
		}
	}
	return &Store{filesCount: counts}, nil
}

func (s *Store) Valid(code string) (bool, string) {
	c := strings.ToUpper(strings.TrimSpace(code))
	if l := len(c); l < 8 || l > 10 {
		return false, "promo code must be 8–10 characters"
	}
	if s.filesCount[c] < 2 {
		return false, "promo code not found in at least two sources"
	}
	return true, ""
}
