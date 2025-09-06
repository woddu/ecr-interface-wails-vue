package main

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

var tracks = [4]string{
	"Core Subject (All Tracks)",
	"Academic Track (except Immersion)",
	"Work Immersion/ Culminating Activity (for Academic Track)",
	"TVL/ Sports/ Arts and Design Track",
}

var sheetsList = [4]string{
	"INPUT DATA",
	"1ST",
	"2ND",
	"Final Semestral Grade",
}

var weightedScores = [4][3]float32{
	{0.25, 0.50, 0.25}, // Core Subject (All Tracks)
	{0.25, 0.45, 0.30}, // Academic Track (except Immersion)
	{0.35, 0.40, 0.25}, // Work Immersion/ Culminating Activity (for Academic Track)
	{0.20, 0.60, 0.20}, // TVL/ Sports/ Arts and Design Track
}

// App struct
type App struct {
	ctx              context.Context
	filePath         string
	fileName         string
	track            string
	firstSem         bool
	wwHighestScores  [10]float32
	ptHighestScores  [10]float32
	examHighestScore float32
	weightedScores   [3]float32
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.firstSem = true
}

func (a *App) Tracks() []string {
	return tracks[:]
}

func (a *App) FileName() string {
	return a.fileName
}

func (a *App) Track() string {
	return a.track
}

func (a *App) ChangeTrack(newTrack string, index int) {
	// check if newTrack is in tracks
	if !slices.Contains(tracks[:], newTrack) && !(index >= 0 && index < len(tracks)) {
		return
	} else if slices.Contains(tracks[:], newTrack) {
		a.track = newTrack
		a.weightedScores = weightedScores[indexOf(tracks[:], a.track)]
	} else if index >= 0 && index < len(tracks) {
		a.track = tracks[index]
		a.weightedScores = weightedScores[indexOf(tracks[:], a.track)]
	}
}

func (a *App) OpenFileDialog() error {

	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select an Excel file",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files", Pattern: "*.xlsx;*.xls"},
			{DisplayName: "All Files", Pattern: "*"},
		},
	})

	if err != nil {
		return err
	}

	if filePath == "" {
		runtime.EventsEmit(a.ctx, "excel:done")
		return nil
	}

	go func(path string) {
		f, err := excelize.OpenFile(path)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to open: %v", err))
			return
		}
		defer f.Close()

		allSheets := f.GetSheetList()

		for _, required := range sheetsList {
			if !slices.Contains(allSheets, required) {
				runtime.EventsEmit(a.ctx, "excel:is_ecr", false)
				return
			}
		}

		runtime.EventsEmit(a.ctx, "excel:is_ecr", true)
		a.filePath = path
		a.fileName = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

		rows, err := f.GetRows(sheetsList[0]) // "INPUT DATA"
		if err != nil {
			return
		}
		if len(rows) < 7 {
			return
		}
		row := rows[7] // 8 row (index 7)

		if len(row) >= 3 {
			track := row[colNameToNumber("AE")]
			a.track = track
			a.weightedScores = weightedScores[indexOf(tracks[:], track)]
		}

		rows, err = f.GetRows(sheetsList[1]) // "1ST"
		if err != nil {
			return
		}

		rowIndex := 10
		if rowIndex < len(rows) {
			row := rows[rowIndex]
			// Columns F (index 5) to O (index 14)
			// Columns S (index 18) to AB (index 27)
			// Column AF (index 31)
			// Make sure the row has enough columns
			if len(row) > colNameToNumber("AF") {
				var score float32

				values := row[colNameToNumber("F") : colNameToNumber("O")+1] // slice is end-exclusive
				for i, v := range values {
					fmt.Sscanf(v, "%f", &score)
					if v == "" {
						score = 0
					}
					a.wwHighestScores[i] = score
				}

				values = row[colNameToNumber("S") : colNameToNumber("AB")+1]
				for i, v := range values {
					fmt.Sscanf(v, "%f", &score)
					if v == "" {
						score = 0
					}
					a.ptHighestScores[i] = score
				}

				fmt.Sscanf(row[colNameToNumber("AF")], "%f", &score)
				a.examHighestScore = score
			}
		}
		sheet := sheetsList[0] // e.g. "INPUT DATA"
		var male []string
		var female []string

		for row := 13; row <= 37; row++ {
			cell, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
			if err != nil {
				continue // skip errors silently
			}
			if cell != "" {
				male = append(male, cell)
			}
		}
		runtime.EventsEmit(a.ctx, "excel:students_male", male)

		for row := 64; row <= 88; row++ {
			cell, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
			if err != nil {
				continue // skip errors silently
			}
			if cell != "" {
				female = append(female, cell)
			}
		}
		runtime.EventsEmit(a.ctx, "excel:students_female", female)

		runtime.EventsEmit(a.ctx, "excel:done_reading")
	}(filePath)

	fmt.Println(a.weightedScores)
	fmt.Println(a.wwHighestScores)
	fmt.Println(a.ptHighestScores)
	fmt.Println(a.examHighestScore)

	return nil
}

type ScoresResult struct {
	WwHighestScores  [10]float32 `json:"wwHighestScores"`
	PtHighestScores  [10]float32 `json:"ptHighestScores"`
	ExamHighestScore float32     `json:"examHighestScore"`
	WeightedScores   [3]float32  `json:"weightedScores"`
}

func (a *App) Scores() ScoresResult {
	fmt.Println(a.weightedScores)
	fmt.Println(a.wwHighestScores)
	fmt.Println(a.ptHighestScores)
	fmt.Println(a.examHighestScore)
	return ScoresResult{
		WwHighestScores:  a.wwHighestScores,
		PtHighestScores:  a.ptHighestScores,
		ExamHighestScore: a.examHighestScore,
		WeightedScores:   a.weightedScores,
	}
}

func (a *App) EditHighestScores(scores [10]float32, writtenWorks bool) [10]float32 {
	if writtenWorks {
		a.wwHighestScores = scores
		return a.wwHighestScores
	} else {
		a.ptHighestScores = scores
		return a.ptHighestScores
	}
}

func (a *App) EditExamHighestScore(score float32) float32 {
	a.examHighestScore = score
	return a.examHighestScore
}

func (a *App) GetStudent(row int) {
	go func(row int) {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to open: %v", err))
			return
		}
		defer f.Close()

		var index int
		if a.firstSem {
			index = 1
		} else {
			index = 2
		}
		rows, err := f.GetRows(sheetsList[index])
		if err != nil {
			return
		}
		if len(rows) < row {
			return
		}
		studentRow := rows[row+12] // row is 1-based index
		if len(studentRow) > colNameToNumber("AF") {
			var score float32

			values := studentRow[colNameToNumber("F") : colNameToNumber("O")+1] // slice is end-exclusive
			for i, v := range values {
				fmt.Sscanf(v, "%f", &score)
				if v == "" {
					score = 0
				}
				// a.wwHighestScores[i] = score
			}
			values = studentRow[colNameToNumber("S") : colNameToNumber("AB")+1]
			for i, v := range values {
				fmt.Sscanf(v, "%f", &score)
				if v == "" {
					score = 0
				}
				// a.ptHighestScores[i] = score
			}
			fmt.Sscanf(studentRow[colNameToNumber("AF")], "%f", &score)
			// a.examHighestScore = score
		}
	}(row)
}

func colNameToNumber(col string) int {
	col = strings.ToUpper(col)
	result := 0
	for i := 0; i < len(col); i++ {
		ch := col[i]
		// Convert ASCII letter to number: 'A' → 1, 'B' → 2, ...
		result = result*26 + int(ch-'A'+1)
	}
	return result - 1
}

func indexOf(slice []string, target string) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}
