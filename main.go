package main

import (
	"database/sql"
	"log"
	"fmt"

	"strings"
	_ "github.com/mattn/go-sqlite3"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"
)

type Question struct {
	ID        int
	Text      string
	Options   []string
	CorrectID int
}

type QuizApp struct {
	questions []Question
	current   int
	score     int
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Тесты, тесты, тесты")

	content := container.NewMax()

	
	db, err := sql.Open("sqlite3", "./quiz.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	quiz := loadQuestionsFromDB(db)

	quizList := widget.NewList(
		func() int {
			return len(quiz.questions)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Question")
			return container.NewVBox(label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			text := obj.(*fyne.Container).Objects[0].(*widget.Label)
			text.SetText(quiz.questions[id].Text)
		})
		totalQuestions := len(quiz.questions)

	quizList.OnSelected = func(id widget.ListItemID) {
		if quiz.current >= len(quiz.questions) {
		
			showResult(content, quiz, totalQuestions, func() {
				quiz.current = 0
				quiz.score = 0
				quizList.Select(0)
			})
			return
		}

		question := quiz.questions[id]
		options := widget.NewRadioGroup(question.Options, func(selected string) {
			selectedID := getSelectedOptionID(selected, question.Options)
			if selectedID == question.CorrectID {
				quiz.score++
			}
			quiz.current++
			if quiz.current < len(quiz.questions) {
				quizList.Select(quiz.current)
			} else {
				showResult(content, quiz, totalQuestions, func() {
					quiz.current = 0
					quiz.score = 0
					quizList.Select(0)
				})
			}
		})

		content.Objects = []fyne.CanvasObject{
			container.NewVBox(widget.NewLabel(wrapText(question.Text, 50)), options),
		}
	}

	split := container.NewHSplit(quizList, content)
	split.Offset = 0.2
	myWindow.SetContent(split)
	myWindow.Resize(fyne.NewSize(480, 360))
	myWindow.ShowAndRun()
}


func wrapText(text string, lineLength int) string {
	wrapped := wrap(text, lineLength)
	return wrapped
}

func wrap(s string, lineLength int) string {
	words := strings.Fields(s)
	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+len(word) < lineLength {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	lines = append(lines, currentLine)
	return strings.Join(lines, "\n")
}

// работаем с базой данных вопросов
func loadQuestionsFromDB(db *sql.DB) QuizApp {
	rows, err := db.Query("SELECT text, option1, option2, option3, option4, correct_id FROM questions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var quiz QuizApp

	for rows.Next() {
		var qText, option1, option2, option3, option4 string
		var correctID int
		err := rows.Scan(&qText, &option1, &option2, &option3, &option4, &correctID)
		if err != nil {
			log.Fatal(err)
		}

		question := Question{
			Text:      qText,
			Options:   []string{option1, option2, option3, option4},
			CorrectID: correctID,
		}

		quiz.questions = append(quiz.questions, question)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return quiz
}

func getSelectedOptionID(selected string, options []string) int {
	for i, option := range options {
		if option == selected {
			return i
		}
	}
	return -1
}



func showResult(content *fyne.Container, quiz QuizApp, totalQuestions int, restartFunc func()) {
	resultLabel := widget.NewLabel(fmt.Sprintf("Тест завершен! Ваша оценка: %d/%d", quiz.score, totalQuestions))
	restartButton := widget.NewButton("Перепройти тест", func() {
		restartFunc()
	})

	content.Objects = []fyne.CanvasObject{
		container.NewVBox(resultLabel, restartButton),
	}
}
