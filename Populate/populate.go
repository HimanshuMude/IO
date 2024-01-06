package populate

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func PopulateFile() {
    file, err := os.OpenFile("db.txt", os.O_APPEND, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    for i := 1; i <= 1000; i++ {
			rand.NewSource(time.Now().UnixNano())
				marks := rand.Intn(101)
        _, err := file.WriteString(fmt.Sprintf("21610%03d Student%d %d\n", i, i,marks))
        if err != nil {
            fmt.Println(err)
            return
        }
    }

    fmt.Println("File populated successfully.")
}