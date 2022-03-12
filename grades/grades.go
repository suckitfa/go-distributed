package grades

import (
	"fmt"
	"sync"
)

// 学生结构体
type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}
	return result / float32(len(s.Grades))
}

// 学生集合
type Students []Student

// 声明一个全局的学生集合
// 加上互斥锁头，保证并发的安全
var (
	students      Students
	studentsMutex sync.Mutex
)

func (ss Students) GetById(id int) (*Student, error) {
	for _, student := range ss {
		if student.ID == id {
			return &student, nil
		}
	}
	return nil, fmt.Errorf("Student with ID %d not found", id)
}

// 考试类型
type GradeType string

const (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
