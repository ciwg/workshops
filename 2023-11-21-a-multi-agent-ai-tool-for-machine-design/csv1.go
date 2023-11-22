package main

import (
        "encoding/csv"
        "fmt"
        "log"
        "os"
        "strconv"
)

func main() {
        file, err := os.Open("file.csv")
        if err != nil {
                log.Fatal(err)
        }
        defer file.Close()

        reader := csv.NewReader(file)
        reader.Comma = ','
        reader.FieldsPerRecord = -1

        csvData, err := reader.ReadAll()
        if err != nil {
                log.Fatal(err)
        }

        var total int
        for _, row := range csvData {
                qty, err := strconv.Atoi(row[1]) // assuming 'qty' is at index 1
                if err != nil {
                        log.Fatal(err)
                }
                total += qty
        }

        fmt.Printf("Total quantity: %d", total)
}
