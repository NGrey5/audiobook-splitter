package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// ffmpeg -i .\input.mp3 -ss {startTimeSecs} -to {endTimeSecs} -c copy {output-file-name.mp3}

type Args struct {
	InputFile string
	LabelFile string
	OutputDirectory string
}

func main() {

	args := getArgs()

	inputFile := args.InputFile
	inputFileExt := path.Ext(inputFile) // Get the file extension
	outputDir := args.OutputDirectory

	fmt.Println("ffmpeg Audiobook Splitter")
	fmt.Println("=============================================================")
	fmt.Printf("Input File: %s\n", inputFile);
	fmt.Printf("Labels File: %s\n", args.LabelFile);
	fmt.Printf("Output Directory: %s\n", outputDir);
	fmt.Println("=============================================================")

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		log.Fatal("Input file could not be found!")
	}

	// Check if label file exists
	if _, err := os.Stat(args.LabelFile); os.IsNotExist(err) {
		log.Fatal("Label file could not be found!")
	}

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating output directory \"%s\"\n%e", outputDir, err)
		}
	}

	labels := parseLabels(args.LabelFile)

	// Loop through each label and split the input file
	for i, label := range labels {
		
		fileName := fmt.Sprintf("%s%s", label.Name, inputFileExt) // The new file name with the same extension as the input
		filePath := path.Join(outputDir, fileName) // The file path within the output directory
		trackNum := fmt.Sprintf("track=%d", i+1)

		fmt.Printf("%d) Start: %s  End: %s  Name: %s  File: %s\n", i+1, label.StartTime, label.EndTime, label.Name, fileName)

		cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", label.StartTime, "-to", label.EndTime, "-c", "copy", "-metadata", trackNum, filePath, "-y")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
    	fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    	return
		}
	}

}

func getArgs() Args {

	var inputFile string
	flag.StringVar(&inputFile, "i", "./input.mp3", "the input file to parse")
	var labelsFile string
	flag.StringVar(&labelsFile, "l", "./labels.txt", "the label file to parse chapter points")
	var outputPath string
	flag.StringVar(&outputPath, "o", "./output", "the output path to create the split files")

	flag.Parse()

  return Args{
		InputFile: inputFile,
		OutputDirectory: outputPath,
		LabelFile: labelsFile,
	}
}

type Label struct {
	StartTime string
	EndTime   string
	Name      string
}

func parseLabels(labelFile string) []Label {
	// Read the label file
	bytes, err := os.ReadFile(labelFile)
	if err != nil {
		log.Fatalf("Failed to read label file \"%s\"", labelFile)
	}

	// Split each line into an array
	lines := strings.Split(string(bytes), "\n")

	var labels []Label

	for _, line := range lines {
		line = strings.Trim(line, "\n") // Remove any new line characters
		line = strings.Trim(line, "\r") // Remove an carriage return characters

		// If the end of the file, break out of the loop
		if strings.TrimSpace(line) == "" {
			break;
		}

		// Split each part START_TIME, END_TIME, CHAPTER_NAME by the tab character
		parts := strings.Split(line, "\t")
		
		// Create a label and append it to the array
		labels = append(labels, Label{StartTime: parts[0], EndTime: parts[1], Name: parts[2]})
	}
	return labels
}