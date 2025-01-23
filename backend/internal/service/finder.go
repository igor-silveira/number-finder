package service

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Finder struct {
	numbers []int
}

type Result struct {
	Index         int  `json:"index"`
	Number        int  `json:"number"`
	IsApproximate bool `json:"is_approximate"`
}

type FinderService interface {
	Find(target int, thresholdPercentage float64) (*Result, error)
}

func loadNumbers(filepath string) ([]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var numbers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for _, field := range fields {
			num, err := strconv.Atoi(field)
			if err != nil {
				return nil, fmt.Errorf("invalid number in file: %w", err)
			}
			numbers = append(numbers, num)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return numbers, nil
}

func (f *Finder) findExact(target int, left, right *int) *Result {
	for *left <= *right {
		mid := (*left + *right) / 2
		if f.numbers[mid] == target {
			return &Result{
				Index:         mid,
				Number:        target,
				IsApproximate: false,
			}
		}
		if f.numbers[mid] < target {
			*left = mid + 1
		} else {
			*right = mid - 1
		}
	}
	return nil
}

func (f *Finder) findAdjacentWithinThreshold(target int, thresholdPercentage float64, left, right int) *Result {
	threshold := float64(target) * thresholdPercentage

	if left < len(f.numbers) {
		if diff := math.Abs(float64(f.numbers[left] - target)); diff < threshold {
			return &Result{
				Index:         left,
				Number:        f.numbers[left],
				IsApproximate: true,
			}
		}
	}

	if right >= 0 {
		if diff := math.Abs(float64(f.numbers[right] - target)); diff < threshold {
			return &Result{
				Index:         right,
				Number:        f.numbers[right],
				IsApproximate: true,
			}
		}
	}

	return nil
}

func (f *Finder) Find(target int, thresholdPercentage float64) (*Result, error) {
	left, right := 0, len(f.numbers)-1

	if exactMatch := f.findExact(target, &left, &right); exactMatch != nil {
		return exactMatch, nil
	}

	if thresholdPercentage == 0 {
		return nil, errors.New("number not found")
	}

	if result := f.findAdjacentWithinThreshold(target, thresholdPercentage, left, right); result != nil {
		return result, nil
	}

	return nil, errors.New("number not found within acceptable threshold")
}

func NewFinder(filepath string) (FinderService, error) {
	numbers, err := loadNumbers(filepath)

	if err != nil {
		return nil, fmt.Errorf("failed to load numbers: %w", err)
	}

	return &Finder{numbers: numbers}, nil
}
