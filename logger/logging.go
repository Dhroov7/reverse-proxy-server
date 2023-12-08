package logger

import (
	"bufio"
	"fmt"
	"os"
)

func Log(text string) error {
	return appendToFile("./log.txt", text)
}

func appendToFile(filename, text string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a writer for the file
	writer := bufio.NewWriter(file)

	// Write the text to the file
	_, err = fmt.Fprintln(writer, text)
	if err != nil {
		return err
	}

	// Flush the writer to ensure the data is written to the file
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
