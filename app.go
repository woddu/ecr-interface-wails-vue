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
	} else if index >= 0 && index < len(tracks) {
		a.track = tracks[index]
	}

	go func(track string) {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(CT)Failed to open: %v", err))
			return
		}
		defer f.Close()
		rows, err := f.GetRows(sheetsList[0]) // "INPUT DATA"
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to get rows: %v", err))
			return
		}
		if len(rows) < 7 {
			runtime.EventsEmit(a.ctx, "excel:error", "Not enough rows in INPUT DATA sheet")
			return
		}

		row := rows[7] // 8 row (index 7)

		// write to row[colNameToNumber("AE")]
		row[colNameToNumber("AE")] = a.track

		if len(row) >= 3 {
			track := row[colNameToNumber("AE")]
			a.track = track
			a.weightedScores = weightedScores[indexOf(tracks[:], a.track)]
		}

		if err := f.Save(); err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to save: %v", err))
			return
		}
		runtime.EventsEmit(a.ctx, "excel:track_changed", a.weightedScores)
	}(a.track)
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
		runtime.EventsEmit(a.ctx, "excel:choose_cancelled")
		return nil
	}

	a.filePath = filePath
	a.fileName = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	go func() {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(OFD)Failed to open: %v", err))
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

		for row := 13; row <= 42; row++ {
			cell, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
			if err != nil {
				continue // skip errors silently
			}
			if cell != "" {
				male = append(male, cell)
			}
		}
		runtime.EventsEmit(a.ctx, "excel:students_male", male)

		for row := 64; row <= 93; row++ {
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
	}()

	return nil
}

func (a *App) ChangeSem(firstSem bool) {
	a.firstSem = firstSem
	fmt.Println("Changed semester to", firstSem)
	go func() {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(OFD)Failed to open: %v", err))
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

		var index int

		if a.firstSem {
			index = 1
		} else {
			index = 2
		}

		rows, err = f.GetRows(sheetsList[index]) // "1ST"
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

		for row := 13; row <= 42; row++ {
			cell, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
			if err != nil {
				continue // skip errors silently
			}
			if cell != "" {
				male = append(male, cell)
			}
		}
		runtime.EventsEmit(a.ctx, "excel:students_male", male)

		for row := 64; row <= 93; row++ {
			cell, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
			if err != nil {
				continue // skip errors silently
			}
			if cell != "" {
				female = append(female, cell)
			}
		}
		runtime.EventsEmit(a.ctx, "excel:students_female", female)
		fmt.Println("Done reading after changing sem")
		runtime.EventsEmit(a.ctx, "excel:done_reading")
	}()
}

type ScoresResult struct {
	WwHighestScores  [10]float32 `json:"wwHighestScores"`
	PtHighestScores  [10]float32 `json:"ptHighestScores"`
	ExamHighestScore float32     `json:"examHighestScore"`
	WeightedScores   [3]float32  `json:"weightedScores"`
}

func (a *App) Scores() ScoresResult {
	return ScoresResult{
		WwHighestScores:  a.wwHighestScores,
		PtHighestScores:  a.ptHighestScores,
		ExamHighestScore: a.examHighestScore,
		WeightedScores:   a.weightedScores,
	}
}

func (a *App) EditHighestScores(scores [10]float32, writtenWorks bool) {
	if writtenWorks {
		a.wwHighestScores = scores
	} else {
		a.ptHighestScores = scores
	}

	go func(writtenWorks bool) {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(EHC)Failed to open: %v", err))
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
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to get rows: %v", err))
			return
		}
		rowIndex := 10
		if rowIndex < len(rows) {
			row := rows[rowIndex]
			if len(row) > colNameToNumber("AF") {
				if writtenWorks {
					for i := range 10 {
						row[colNameToNumber("F")+i] = fmt.Sprintf("%.0f", a.wwHighestScores[i])
					}
				} else {
					for i := range 10 {
						row[colNameToNumber("S")+i] = fmt.Sprintf("%.0f", a.ptHighestScores[i])
					}
				}
				if err := f.Save(); err != nil {
					runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to save: %v", err))
					return
				}
				if writtenWorks {
					runtime.EventsEmit(a.ctx, "excel:done_editing_highest_scores", writtenWorks, a.wwHighestScores)
				} else {
					runtime.EventsEmit(a.ctx, "excel:done_editing_highest_scores", writtenWorks, a.ptHighestScores)
				}
			}
		}
	}(writtenWorks)

}

func (a *App) EditExamHighestScore(score float32) {
	a.examHighestScore = score
	go func() {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(EEHS)Failed to open: %v", err))
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
		if len(rows) < 7 {
			return
		}
		row := rows[7]

		if len(row) > 3 {
			row[colNameToNumber("AF")] = fmt.Sprintf("%.0f", a.examHighestScore)

			if err := f.Save(); err != nil {
				runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to save: %v", err))
				return
			}
		}
		runtime.EventsEmit(a.ctx, "excel:done_editing_exam_highest_score", a.examHighestScore)
	}()
}

type StudentInfo struct {
	Row          int         `json:"row"`
	Name         string      `json:"name"`
	WrittenWorks [10]float32 `json:"writtenWorks"`
	Performance  [10]float32 `json:"performance"`
	Exam         float32     `json:"exam"`
}

func (a *App) GetStudent(row int) {
	go func(row int) {
		var studentInfo StudentInfo
		studentInfo.Row = row + 12
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(GS)Failed to open: %v", err))
			return
		}
		defer f.Close()

		var index int

		if a.firstSem {
			index = 1
		} else {
			index = 2
		}

		rows, err := f.GetRows(sheetsList[0]) // "INPUT DATA"
		if err != nil {
			return
		}

		if len(rows) < row {
			return
		}

		studentInfo.Name = rows[studentInfo.Row][colNameToNumber("B")] // row is 1-based index

		rows, err = f.GetRows(sheetsList[index])

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
				studentInfo.WrittenWorks[i] = score
				// a.wwHighestScores[i] = score
			}
			values = studentRow[colNameToNumber("S") : colNameToNumber("AB")+1]
			for i, v := range values {
				fmt.Sscanf(v, "%f", &score)
				if v == "" {
					score = 0
				}
				studentInfo.Performance[i] = score
				// a.ptHighestScores[i] = score
			}
			fmt.Sscanf(studentRow[colNameToNumber("AF")], "%f", &score)
			studentInfo.Exam = score
			runtime.EventsEmit(a.ctx, "excel:done_getting_student", studentInfo)
		}
	}(row)
}

func (a *App) EditStudentScores(studentRow int, scores [10]float32, writtenWorks bool) {

	go func(studentRow int, scores [10]float32, writtenWorks bool) {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(ESS)Failed to open: %v", err))
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
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to get rows: %v", err))
			return
		}

		if studentRow < len(rows) {
			row := rows[studentRow]
			if len(row) > colNameToNumber("AF") {
				if writtenWorks {
					for i := range 10 {
						row[colNameToNumber("F")+i] = fmt.Sprintf("%.0f", scores[i])
					}
				} else {
					for i := range 10 {
						row[colNameToNumber("S")+i] = fmt.Sprintf("%.0f", scores[i])
					}
				}

				if err := f.Save(); err != nil {
					runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to save: %v", err))
					return
				}

				if writtenWorks {
					runtime.EventsEmit(a.ctx, "excel:done_editing_student_scores", writtenWorks, scores)
				} else {
					runtime.EventsEmit(a.ctx, "excel:done_editing_student_scores", writtenWorks, scores)
				}
			}
		}
	}(studentRow, scores, writtenWorks)

}

func (a *App) EditStudentExamScore(studentRow int, score float32) {
	go func(studentRow int, score float32) {
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
		if len(rows) < studentRow {
			return
		}
		row := rows[studentRow]

		if len(row) > 3 {
			row[colNameToNumber("AF")] = fmt.Sprintf("%.0f", score)

			if err := f.Save(); err != nil {
				runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to save: %v", err))
				return
			}
		}
		runtime.EventsEmit(a.ctx, "excel:done_editing_student_exam_score", score)
	}(studentRow, score)
}

func (a *App) AddStudent(name string, isMale bool) {
	go func(name string, isMale bool) {
		f, err := excelize.OpenFile(a.filePath)
		if err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("(AS)Failed to open: %v", err))
			return
		}
		defer f.Close()

		sheet := sheetsList[0] // e.g. "INPUT DATA"
		var names []string
		if isMale {
			for row := 13; row <= 42; row++ {
				cell, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
				if err != nil {
					continue // skip errors silently
				}
				if cell != "" {
					names = append(names, cell)
				}
			}
		} else {
			for row := 64; row <= 93; row++ {
				cell, err := f.GetCellValue(sheet, fmt.Sprintf("B%d", row))
				if err != nil {
					continue // skip errors silently
				}
				if cell != "" {
					names = append(names, cell)
				}
			}
		}

		names = append(names, strings.ToUpper(name))
		slices.Sort(names)

		if isMale {
			for i, v := range names {
				cell := fmt.Sprintf("B%d", i+13) // B13 onwards
				if err := f.SetCellValue(sheet, cell, v); err != nil {
					runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to write to names: %v", err))
				}
			}
		} else {
			for i, v := range names {
				cell := fmt.Sprintf("B%d", i+64) // B64 onwards
				if err := f.SetCellValue(sheet, cell, v); err != nil {
					runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to write to names: %v", err))
				}
			}
		}

		if err := f.Save(); err != nil {
			runtime.EventsEmit(a.ctx, "excel:error", fmt.Sprintf("Failed to save: %v", err))
			return
		}
		if isMale {
			runtime.EventsEmit(a.ctx, "excel:students_male", names)
		} else {
			runtime.EventsEmit(a.ctx, "excel:students_female", names)
		}

	}(name, isMale)
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
