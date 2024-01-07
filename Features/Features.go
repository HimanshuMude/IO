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

func NewClass(data chan *Class, wg *sync.WaitGroup) {
	// get info from db.txt file
	wg.Add(1)
	var myWG sync.WaitGroup
	var mu sync.Mutex
	file, err := os.Open("db.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	var studentsData []Student

	scanner := bufio.NewScanner(file)


	/************************************
			Implementing Goroutines
	************************************/

	for scanner.Scan() {
		line := scanner.Text()
		myWG.Add(1)
		go func(line string){
				
				// line= "2262005 Utkarsh 100"
				fields := strings.Split(line, " ")

				// fields=["22620005","Utkarsh","100"]

				marks, _ := strconv.Atoi(fields[2])

				student := Student{
					PRN:   fields[0],
					Name:  fields[1],
					Marks: marks,
				}
				mu.Lock()
				studentsData = append(studentsData, student)
				mu.Unlock()
				myWG.Done()

		}(line)
		
	}
	myWG.Wait()
	data <- &Class{students: studentsData}
	wg.Done()
}

// Using go routine
func (c *Class) WriteToFile(wg *sync.WaitGroup) {
	t := time.Now()
	var lines []string

	for _, student := range c.students {
		sMarks := strconv.Itoa(student.Marks)
		line := fmt.Sprintf("%s %s %s\n", student.PRN, student.Name, sMarks)
		lines = append(lines, line)
	}

	err := os.WriteFile("db.txt", []byte(strings.Join(lines, "")), 0644)
	if err != nil {
		fmt.Println(err)
	}

	defer wg.Done()
	fmt.Println(time.Since(t))
}

// Without Go routine
func (c *Class) WriteToFileSerial() {
	t := time.Now()
	var lines []string

	for _, student := range c.students {
		sMarks := strconv.Itoa(student.Marks)
		line := fmt.Sprintf("%s %s %s\n", student.PRN, student.Name, sMarks)
		lines = append(lines, line)
	}

	err := os.WriteFile("db.txt", []byte(strings.Join(lines, "")), 0644)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(time.Since(t))
}

func (c *Class) AddStudent() {
	// var wg sync.WaitGroup
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


	// IMPORTANT CHANGE

	// 1: sequential
	c.WriteToFileSerial()
	fmt.Println("Doing something else")

	// 2: Async
	// wg.Add(1)
	// go func(){
	// 	c.WriteToFile(&wg)
	// }()
	// fmt.Println("Doing something else")
	// wg.Wait()

	// NOTE!!!!
	// 	the first version executes WriteToFileSerial() synchronously, blocking the main thread until the file writing is complete
	// the second version launches a goroutine to handle file writing asynchronously, allowing the main thread to continue with other tasks while the file writing happens in the background.
	// the main thread doesn't idle in either case. it's either actively executing the file writing task or managing other tasks while the goroutine runs.

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
	var wg sync.WaitGroup
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

			wg.Add(1)
			go c.WriteToFile(&wg)
			wg.Wait()
			return
		}
	}

	fmt.Println("Student not found")
}

func (c *Class) DeleteStudent() {

	var wg sync.WaitGroup

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

	// have to create a slide for this
	dataBase = append(dataBase[:index], dataBase[index+1:]...)
	c.students = dataBase

	fmt.Println("Student deleted successfully!")
	fmt.Println()
	wg.Add(1)
	go c.WriteToFile(&wg)
	wg.Wait()
}

// we are running this for 40 times just to show the difference go routines make when used in heavy data processing scenarios
// in normal scenarios code without go routines works better
// ex: try running the GetStat & GetStatG function without those loops
// the function without concurrency will surely perform better because for small size of data it is efficient as compared to go routines

// without go routines
func (c Class) GetStat() {
	var sum, avg, low, high int

	sum = sum + c.findSum()

	avg = c.findAverage()

	low = c.findLowest()

	high = c.findHighest()

	fmt.Println(sum, avg, low, high)
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
