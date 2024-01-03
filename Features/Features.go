package features

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Student struct {
	PRN   string
	Name  string
	Marks int
}

type Class struct {
	students []Student
}

func NewClass() *Class {
	// get info from db.txt file
	file, err := os.Open("db.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	var studentsData []Student

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		// line= "2262005 Utkarsh 100"
		fields := strings.Split(line, " ")

		// fields=["22620005","Utkarsh","100"]

		marks, _ := strconv.Atoi(fields[2])

		student := Student{
			PRN:   fields[0],
			Name:  fields[1],
			Marks: marks,
		}

		studentsData = append(studentsData, student)
	}
	return &Class{students: studentsData}
}

func (c *Class) WriteToFile() {
	var lines []string

	for _, student := range c.students {
		sMarks:=strconv.Itoa(student.Marks)
		line := fmt.Sprintf("%s %s %s\n", student.PRN, student.Name, sMarks)
		lines = append(lines, line)
	}

	err := os.WriteFile("db.txt", []byte(strings.Join(lines, "")), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Class) AddStudent() {
	fmt.Println("Enter the first name of the student:")
	var fname string
	fmt.Scanln(&fname)

	fmt.Println("Enter the PRN of the student:")
	var prn string
	fmt.Scanln(&prn)

	fmt.Println("Enter the marks of the student:")
	var marks int
	fmt.Scanln(&marks)

	newStudent := Student{
		PRN:   prn,
		Name:  fname,
		Marks: marks,
	}

	c.students = append(c.students, newStudent)
	c.WriteToFile()
}

func (c Class) ShowStudents() {
	for _, student := range c.students {
		fmt.Println("PRN:", student.PRN)
		fmt.Println("Name:", student.Name)
		fmt.Println("Marks:", student.Marks)
		fmt.Println()
	}
	fmt.Println()
}

func (c *Class) UpdateStudent() {
	fmt.Println("Enter the PRN of the student to be updated:")
	var prn string
	fmt.Scanln(&prn)

	for i, student := range c.students {
		if student.PRN == prn {
			fmt.Println("Enter the first name of the student:")
			var fname string
			fmt.Scanln(&fname)

			fmt.Println("Enter the marks of the student:")
			var marks int
			fmt.Scanln(&marks)

			c.students[i] = Student{PRN: prn, Name: fname, Marks: marks}

			fmt.Println("Student updated successfully!")
			fmt.Println()

			c.WriteToFile()
			return
		}
	}

	fmt.Println("Student not found")
}

func (c *Class) DeleteStudent() {

	dataBase := c.students

	fmt.Println("Enter the PRN of the student to be deleted:")
	var prn string
	fmt.Scanln(&prn)

	var index int

	for i, student := range c.students {
		if student.PRN == prn {
			index = i
		}
	}

	dataBase = append(dataBase[:index], dataBase[index+1:]...)
	c.students = dataBase

	fmt.Println("Student deleted successfully!")
	fmt.Println()

	c.WriteToFile()
}

// we are running this for 40 times just to show the difference go routines make when used in heavy data processing scenarios
// in normal scenarios code without go routines works better 
// ex: try running the GetStat & GetStatG function without those loops
// the function without concurrency will surely perform better because for small size of data it is efficient as compared to go routines

// without go routines
func (c Class) GetStatG() {
	t := time.Now()
	var sum, avg, low, high int
	// find Sum
	for i := 0; i < 40; i++ {
		sum = sum + c.findSum()
	}
	for i := 0; i < 40; i++ {
		avg = c.findAverage()
	}
	for i := 0; i < 40; i++ {
		low = c.findLowest()
	}
	for i := 0; i < 40; i++ {
		high = c.findHighest()
	}

	fmt.Println(time.Since(t))
	fmt.Println(sum, avg, low, high)
}

// with goroutine
func (c Class) GetStat() {
	t := time.Now()
	var wg sync.WaitGroup
	var sum, avg, low, high int
	// find Sum
	wg.Add(4)
	go func() {
		for i := 0; i < 40; i++ {
			sum = c.findSum()
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 40; i++ {
			avg = c.findAverage()
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 40; i++ {
			low = c.findLowest()
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 40; i++ {
			high = c.findHighest()
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println(sum, avg, low, high)
	fmt.Println(time.Since(t))
}

func (c Class) findSum() int {
	sum := 0

	tempData := c.students

	for _, v := range tempData {
		sum += v.Marks
	}
	return sum
}

func (c Class) findAverage() int {
	return c.findSum() / len(c.students)
}

func (c Class) findLowest() int {
	lowest := 200

	tempData := c.students

	for _, v := range tempData {
		lowest = min(lowest, v.Marks)
	}
	return lowest
}

func (c Class) findHighest() int {
	highest := 0

	tempData := c.students

	for _, v := range tempData {
		highest = max(highest, v.Marks)
	}
	return highest
}
