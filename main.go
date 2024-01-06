package main

import (
	features "cron/Features"
	populate "cron/Populate"
	"fmt"
	"os"
	"sync"
)

func loadChoices() {
	fmt.Println("Choose one from the following options:")
	fmt.Println("1. Add a student")
	fmt.Println("2. View all students")
	fmt.Println("3. Update a student")
	fmt.Println("4. Delete a student")
	fmt.Println("5. Get Stats")
	fmt.Println("6. Exit")
	fmt.Println()
}

func main() {
	populate.PopulateFile()
	data := make(chan *features.Class)
	var wg sync.WaitGroup
	
	go features.NewClass(data, &wg)
	wg.Wait()
	class := <-data
	close(data)
	fmt.Println("Welcome to the Student DB!")
	fmt.Println()
	for {

		loadChoices()

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			class.AddStudent()
		case 2:
			class.ShowStudents()
		case 3:
			class.UpdateStudent()
		case 4:
			class.DeleteStudent()
		case 5:
			class.GetStat()
		case 6:
			// show with exit status 1 as well
			os.Exit(0)
		default:
			fmt.Println("Invalid choice!")
			fmt.Println()
		}
	}

}
