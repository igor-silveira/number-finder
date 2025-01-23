package service

import (
	"os"
	"testing"
)

type testFile struct {
	Name     string
	Content  string
	WantErr  bool
	Filepath string
}

func createTestFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}

func createTestFiles(t *testing.T, files []testFile) []testFile {
	t.Helper()
	for i := range files {
		if files[i].Name != "non-existent file" {
			files[i].Filepath = createTestFile(t, files[i].Content)
		} else {
			files[i].Filepath = "nonexistent.txt"
		}
	}
	return files
}

func cleanupTestFiles(t *testing.T, files []testFile) {
	t.Helper()
	for _, file := range files {
		if file.Name != "non-existent file" {
			if err := os.Remove(file.Filepath); err != nil {
				t.Errorf("failed to remove test file: %v", err)
			}
		}
	}
}

func TestNewFinder(t *testing.T) {
	testFiles := []testFile{
		{
			Name:    "valid file with sorted numbers",
			Content: "1\n2\n3\n4\n5\n6\n7\n8\n9\n10",
			WantErr: false,
		},
		{
			Name:    "invalid number in file",
			Content: "1\n2\n3\ninvalid\n5",
			WantErr: true,
		},
		{
			Name:    "non-existent file",
			WantErr: true,
		},
	}

	testFiles = createTestFiles(t, testFiles)
	defer cleanupTestFiles(t, testFiles)

	for _, tt := range testFiles {
		t.Run(tt.Name, func(t *testing.T) {
			finder, err := NewFinder(tt.Filepath)
			if (err != nil) != tt.WantErr {
				t.Errorf("NewFinder() error = %v, wantErr %v", err, tt.WantErr)
				return
			}
			if !tt.WantErr && finder == nil {
				t.Error("NewFinder() returned nil but wanted valid Finder")
			}
		})
	}
}

func TestFinder_Find(t *testing.T) {
	content := []byte("1\n3\n5\n7\n9\n11\n13\n15")
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	finder, err := NewFinder(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to create Finder: %v", err)
	}

	tests := []struct {
		name             string
		target           int
		thresholdPercent float64
		wantResult       Result
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name:             "exact match",
			target:           7,
			thresholdPercent: 0.1,
			wantResult: Result{
				Index:         3,
				Number:        7,
				IsApproximate: false,
			},
			wantErr: false,
		},
		{
			name:             "approximate match within threshold",
			target:           8,
			thresholdPercent: 0.2,
			wantResult: Result{
				Index:         4,
				Number:        9,
				IsApproximate: true,
			},
			wantErr: false,
		},
		{
			name:             "no match with zero threshold",
			target:           8,
			thresholdPercent: 0,
			wantResult:       Result{},
			wantErr:          true,
			wantErrMsg:       "number not found",
		},
		{
			name:             "no match within threshold",
			target:           100,
			thresholdPercent: 0.1,
			wantResult:       Result{},
			wantErr:          true,
			wantErrMsg:       "number not found within acceptable threshold",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := finder.Find(tt.target, tt.thresholdPercent)

			if tt.wantErr {
				if err == nil {
					t.Error("Find() expected error but got none")
					return
				}
				if err.Error() != tt.wantErrMsg {
					t.Errorf("Find() error message = %v, want %v", err.Error(), tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("Find() unexpected error: %v", err)
				return
			}

			if result.Index != tt.wantResult.Index {
				t.Errorf("Find() Index = %v, want %v", result.Index, tt.wantResult.Index)
			}
			if result.Number != tt.wantResult.Number {
				t.Errorf("Find() Number = %v, want %v", result.Number, tt.wantResult.Number)
			}
			if result.IsApproximate != tt.wantResult.IsApproximate {
				t.Errorf("Find() IsApproximate = %v, want %v", result.IsApproximate, tt.wantResult.IsApproximate)
			}
		})
	}
}

func TestFinder_FindClosestWithinThreshold(t *testing.T) {
	numbers := []int{1, 3, 5, 7, 9, 11, 13, 15}
	finder := &Finder{numbers: numbers}

	tests := []struct {
		name             string
		target           int
		thresholdPercent float64
		want             *Result
		checkApproximate bool
	}{
		{
			name:             "exact match",
			target:           7,
			thresholdPercent: 0.1,
			want: &Result{
				Index:  3,
				Number: 7,
			},
			checkApproximate: false,
		},
		{
			name:             "approximate match",
			target:           8,
			thresholdPercent: 0.2,
			want: &Result{
				Index:         4,
				Number:        9,
				IsApproximate: true,
			},
			checkApproximate: true,
		},
		{
			name:             "no match within threshold",
			target:           100,
			thresholdPercent: 0.1,
			want:             nil,
			checkApproximate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := finder.Find(tt.target, tt.thresholdPercent)

			if got != nil && err != nil {
				t.Errorf("Find() unexpected error: %v", err)
			}

			if tt.want == nil {
				if got != nil {
					t.Errorf("findClosestWithinThreshold() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("findClosestWithinThreshold() = nil, want non-nil")
			}

			if got.Index != tt.want.Index {
				t.Errorf("findClosestWithinThreshold() Index = %v, want %v", got.Index, tt.want.Index)
			}

			if got.Number != tt.want.Number {
				t.Errorf("findClosestWithinThreshold() Number = %v, want %v", got.Number, tt.want.Number)
			}

			if tt.checkApproximate && got.IsApproximate != tt.want.IsApproximate {
				t.Errorf("findClosestWithinThreshold() IsApproximate = %v, want %v", got.IsApproximate, tt.want.IsApproximate)
			}
		})
	}
}
