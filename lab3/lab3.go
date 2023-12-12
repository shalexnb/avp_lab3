package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)


type Question struct {
	ID         int
	Text       string
	Options    []sql.NullString
	CorrectIdx int
}



type Quiz struct {
    Questions []Question
    Score     int
}

func displayQuestion(q Question) {
    fmt.Println(q.Text)
    for i, option := range q.Options {
        if option.Valid {
            fmt.Printf("%d. %s\n", i+1, option.String)
        } else {
            //fmt.Printf("%d. [NULL]\n", i+1)
        }
    }
    fmt.Print("Выберите номер вашего ответа: ")
}


func checkAnswer(userAnswer int, correctAnswer int) bool {
    return userAnswer == correctAnswer
}

func main() {
    db, err := sql.Open("sqlite3", "./quiz.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT * FROM questions")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    var quiz Quiz

    for rows.Next() {
        var question Question
        var option1, option2, option3, option4 sql.NullString
        err := rows.Scan(&question.ID, &question.Text, &option1, &option2, &option3, &option4, &question.CorrectIdx)
        if err != nil {
            panic(err)
        }
    
        question.Options = []sql.NullString{option1, option2, option3, option4}
    
        quiz.Questions = append(quiz.Questions, question)
    }
    

    for _, question := range quiz.Questions {
        displayQuestion(question)
        var userAnswer int
        fmt.Scanln(&userAnswer)

        if checkAnswer(userAnswer, question.CorrectIdx+1) {
            fmt.Println("Правильно!\n")
            quiz.Score++
        } else {
            correctOption := question.Options[question.CorrectIdx-1]
            if correctOption.Valid {
                fmt.Printf("Неправильно. Правильный ответ: %s\n\n", correctOption.String)
            } else {
                fmt.Printf("Неправильно. Правильный ответ: [NULL]\n\n")
            }
        }
        
    }

    fmt.Printf("Ваш итоговый балл: %d из %d\n", quiz.Score, len(quiz.Questions))
}
