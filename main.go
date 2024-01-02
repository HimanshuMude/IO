package main

import (
	features "cron/Features"
)

func main() {
	class:=features.NewClass()
	// class.AddStudent("22620005","Suyash","150")
	// class.AddStudent("22620003","Shivam","50")
	// class.AddStudent("22620002","Soham","60")
	// class.AddStudent("22620010","Saurabh","75")
	// class.AddStudent("22620012","sanket","99")
	// class.AddStudent("22620009","shardul","10")
	// class.ShowStudents()
	// class.DeleteStudent("22620005")
	// class.ShowStudents()
	// class.UpdateStudent("22620003","Speed","80")
	// class.ShowStudents()
	class.GetStat()
	class.GetStatG()

}
