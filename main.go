package main

import (
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

var verbose bool

func main() {
    var path string
    var filetypes string
    var outputFileName string

    flag.StringVar(&path, "path", "", "Specify the path to the codebase directory.")
    flag.StringVar(&filetypes, "filetype", "", "Specify file extensions to include, separated by commas (e.g., .py,.js,.html). Leave empty to include all files.")
    flag.StringVar(&outputFileName, "output", "combined_code", "Specify the base name of the output file (without extension).")
    flag.BoolVar(&verbose, "verbose", false, "Enable verbose output for debugging.")
    flag.Parse()

    outputFile, err := os.Create(outputFileName + ".txt")
    if err != nil {
        fmt.Println("Error creating output file:", err)
        os.Exit(1)
    }
    defer outputFile.Close()

    outputFilePath, err := filepath.Abs(outputFile.Name())
    if err != nil {
        fmt.Println("Error getting absolute path of output file:", err)
        os.Exit(1)
    }

    if verbose {
        fmt.Println("Output file path:", outputFilePath)
    }

    extensions := strings.Split(filetypes, ",")
    extMap := make(map[string]bool)
    for _, ext := range extensions {
        if ext != "" {
            extMap[ext] = true
        }
    }

    if verbose {
        fmt.Println("Starting to walk through the directory...")
    }

    err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            if verbose {
                fmt.Println("Error accessing file:", filePath, err)
            }
            return err
        }
        if verbose {
            fmt.Println("Processing:", filePath)
        }

        absFilePath, err := filepath.Abs(filePath)
        if err != nil {
            if verbose {
                fmt.Println("Error determining absolute path for file:", filePath, err)
            }
            return err
        }

        if info.IsDir() {
            if strings.HasPrefix(filepath.Base(filePath), ".") {
                if verbose {
                    fmt.Println("Skipping hidden directory:", filePath)
                }
                return filepath.SkipDir
            }
        } else {
            if absFilePath == outputFilePath {
                if verbose {
                    fmt.Println("Skipping output file:", filePath)
                }
                return nil
            }
            if info.Mode().Perm()&0111 != 0 {
                if verbose {
                    fmt.Println("Skipping executable file:", filePath)
                }
                return nil
            }
            if len(extMap) == 0 || extMap[filepath.Ext(filePath)] {
                if err := appendToFile(outputFile, filePath); err != nil {
                    if verbose {
                        fmt.Println("Error writing file to output:", filePath, err)
                    }
                    return err
                }
            }
        }
        return nil
    })

    if err != nil {
        fmt.Println("Error walking the path:", err)
        os.Exit(1)
    }

    fmt.Println("All files have been combined into", outputFile.Name())
}

func appendToFile(outputFile *os.File, filePath string) error {
    inputFile, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    if verbose {
        fmt.Println("Appending file to output:", filePath)
    }
    _, err = io.Copy(outputFile, inputFile)
    return err
}

