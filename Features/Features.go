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
	Marks string
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
		fields := strings.Split(line, " ")

		student := Student{
			PRN:   fields[0],
			Name:  fields[1],
			Marks: fields[2],
		}

		studentsData = append(studentsData, student)
	}
	return &Class{students: studentsData}
}

func (c *Class) WriteToFile() {
	file, err := os.Create("db.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, student := range c.students {
		line := fmt.Sprintf("%s %s %s\n", student.PRN, student.Name, student.Marks)
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (c *Class) AddStudent(PRN, Name, Marks string) {
	newStudent := Student{
		PRN:   PRN,
		Name:  Name,
		Marks: Marks,
	}

	c.students = append(c.students, newStudent)
	c.WriteToFile()
}

func (c Class) ShowStudents() {
	fmt.Println(c.students)
}

func (c *Class) UpdateStudent(PRN, Name, Marks string) {
	tempData := c.students

	for i, v := range tempData {
		if v.PRN == PRN {
			tempData[i].Name = Name
			tempData[i].Marks = Marks
			break
		}
	}
	c.students = tempData
	c.WriteToFile()

}

func (c *Class) DeleteStudent(PRN string) {
	var index int

	tempData := c.students

	for i, v := range tempData {
		if v.PRN == PRN {
			index = i
			break
		}
	}

	tempData = append(tempData[:index], tempData[index+1:]...)
	c.students = tempData
	c.WriteToFile()

}

func (c Class) GetStatG() {
	t:=time.Now()
	var sum, avg, low, high int
	// find Sum
	for i:=0;i<20;i++{
		sum =sum+ c.findSum()
	}
	for i:=0;i<20;i++{
		avg = c.findAverage()
	}
	for i:=0;i<20;i++{
		low = c.findLowest()
	}
	for i:=0;i<20;i++{
		high = c.findHighest()
	}

	fmt.Println(time.Since(t))
	fmt.Println(sum, avg, low, high)
}

func (c Class) GetStat() {
	t := time.Now()
	var wg sync.WaitGroup
	var sum, avg, low, high int
	// find Sum
	wg.Add(4)
	go func() {
		for i:=0;i<20;i++{
			sum = c.findSum()
		}
		wg.Done()
	}()

	go func() {
		for i:=0;i<20;i++{
			avg = c.findAverage()
		}
		wg.Done()
	}()

	go func() {
		for i:=0;i<20;i++{
			low = c.findLowest()
		}
		wg.Done()
	}()

	go func() {
		for i:=0;i<20;i++{
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
		val, err := strconv.Atoi(v.Marks)
		if err != nil {
			fmt.Println(err)
		}
		sum += val
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
		val, err := strconv.Atoi(v.Marks)
		if err != nil {
			fmt.Println(err)
		}
		lowest = min(lowest, val)
	}
	return lowest
}

func (c Class) findHighest() int {
	highest := 0

	tempData := c.students

	for _, v := range tempData {
		val, err := strconv.Atoi(v.Marks)
		if err != nil {
			fmt.Println(err)
		}
		highest = max(highest, val)
	}
	return highest
}
