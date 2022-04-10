package main

import "fmt"

type CommonAttributes interface {
	GetMyName() string
	GetMyAge() int
}

type Person struct {
	Name string
	Age  int
}

func NewPerson(name string, age int) *Person {
	return &Person{
		Name: name,
		Age:  age,
	}
}

func (p Person) GetMyName() string {
	return p.Name
}

func (p Person) GetMyAge() int {
	return p.Age
}

type Student struct {
	CommonAttributes
	Class  string
	RollNo string
}

type StudentInfo struct {
	Person
	Class  string
	RollNo string
}

func NewStudent(name string, age int) *Student {
	return &Student{
		Class:  "Computer",
		RollNo: "COMP001",
		CommonAttributes: Person{
			Name: name,
			Age:  age,
		},
	}
}

func NewStudentInfo(name string, age int) *StudentInfo {
	return &StudentInfo{
		Class:  "Computer",
		RollNo: "COMP001",
		// Info: Person{
		// 	Name: name,
		// 	Age:  age,
		// },
	}
}

func main() {
	s := NewStudent("Aaron", 5)
	fmt.Println("Age is %v", s.GetMyAge())
	fmt.Println("Age is ----- %v", s.CommonAttributes.GetMyAge())

	fmt.Println("-----------------------------")
	s1 := NewStudentInfo("Aaron", 50)
	fmt.Println("Age is %v", s1.GetMyAge())
}
