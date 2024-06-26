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

    flag.Usage = customUsage
    if flag.NFlag() == 0 {
        flag.Usage()
        os.Exit(1)
    }

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

    totalTokens := 0
    err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            if verbose {
                fmt.Println("Error accessing file:", filePath, err)
            }
            return err
        }

        absFilePath, err := filepath.Abs(filePath)
        if err != nil {
            if verbose {
                fmt.Println("Error determining absolute path for file:", filePath, err)
            }
            return err
        }

        if info.IsDir() {
            if strings.HasPrefix(filepath.Base(filePath), ".") || filepath.Base(filePath) == "node_modules" {
                if verbose {
                    fmt.Println("Skipping directory:", filePath)
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

            fileName := filepath.Base(filePath)
            if fileName == "go.mod" || fileName == "go.sum" {
                if verbose {
                    fmt.Println("Skipping file:", filePath)
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
                tokens, err := appendToFile(outputFile, filePath)
                if err != nil {
                    if verbose {
                        fmt.Println("Error writing file to output:", filePath, err)
                    }
                    return err
                }
                totalTokens += tokens
            }
        }
        return nil
    })

    if err != nil {
        fmt.Println("Error walking the path:", err)
        os.Exit(1)
    }

    fmt.Println("Total Tokens:", totalTokens)
    fmt.Println("All files have been combined into", outputFile.Name())
}

func customUsage() {
    fmt.Println("Usage: codectx")
    fmt.Println("\nFlags:")
    flag.PrintDefaults()
}

func appendToFile(outputFile *os.File, filePath string) (int, error) {
    inputFile, err := os.Open(filePath)
    if err != nil {
        return 0, err
    }
    defer inputFile.Close()

    content, err := io.ReadAll(inputFile)
    if err != nil {
        return 0, err
    }

    tokenCount := countTokens(string(content))

    _, err = outputFile.WriteString(fmt.Sprintf("File: %s\n", filepath.Base(filePath)))
    if err != nil {
        return tokenCount, err
    }

    _, err = outputFile.Write(content)
    if err != nil {
        return tokenCount, err
    }

    _, err = outputFile.WriteString("\n\n")
    return tokenCount, err
}

func countTokens(text string) int {
    tokens := strings.Fields(text)
    return len(tokens)
}
