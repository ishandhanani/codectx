package main

import (
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    var path string
    var filetypes string

    flag.StringVar(&path, "path", "", "Specify the path to the codebase directory.")
    flag.StringVar(&filetypes, "filetype", "", "Specify file extensions to include, separated by commas (e.g., .py,.js,.html). Leave empty to include all files.")
    flag.Parse()

    if path == "" {
        fmt.Println("Error: Path is required.")
        flag.Usage()
        os.Exit(1)
    }

    if _, err := os.Stat(path); os.IsNotExist(err) {
        fmt.Println("Error: Provided path does not exist.")
        os.Exit(1)
    }

    outputFile, err := os.Create("combined_code.txt")
    if err != nil {
        fmt.Println("Error creating output file:", err)
        os.Exit(1)
    }
    defer outputFile.Close()

    extensions := strings.Split(filetypes, ",")
    extMap := make(map[string]bool)
    for _, ext := range extensions {
        if ext != "" {
            extMap[ext] = true
        }
    }

    err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        // Skip hidden directories
        if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
            return filepath.SkipDir
        }
        if !info.IsDir() && (len(extMap) == 0 || extMap[filepath.Ext(filePath)]) {
            return appendToFile(outputFile, filePath)
        }
        return nil
    })

    if err != nil {
        fmt.Println("Error walking the path:", err)
        os.Exit(1)
    }

    fmt.Println("All files have been combined into combined_code.txt")
}

func appendToFile(outputFile *os.File, filePath string) error {
    inputFile, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer inputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    return err
}

