package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// processLogLine procesa una línea del archivo de log y escribe en el CSV.
func processLogLine(line string, csvWriter *csv.Writer) error {
	fmt.Println("Procesando línea:", line)

	re := regexp.MustCompile(`nameHost=(?P<host>[^,]+),nameIP=(?P<ip>[^,]+),dateFormat=(?P<date>[^,]+).*`)
	match := re.FindStringSubmatch(line)

	if len(match) > 0 {
		host := match[re.SubexpIndex("host")]
		ip := match[re.SubexpIndex("ip")]

		fmt.Println("Coincidencia host/ip:", host, ip)

		record := []string{host, ip, "OK"}
		err := csvWriter.Write(record)
		if err != nil {
			return fmt.Errorf("error writing record to CSV: %w", err)
		}
		return nil
	}

	stateRe := regexp.MustCompile(`(?P<STATE>STATE_\w+),.*`)
	stateMatch := stateRe.FindStringSubmatch(line)

	if len(stateMatch) > 0 {
		state := stateMatch[stateRe.SubexpIndex("STATE")]
		fmt.Println("Coincidencia de estado:", state)
		record := []string{"-", "-", state}
		err := csvWriter.Write(record)
		if err != nil {
			return fmt.Errorf("error writing state to CSV: %w", err)
		}
		return nil
	}

	fmt.Println("No hubo coincidencia")
	return nil
}

// processFile lee un archivo de log línea por línea y procesa cada línea.
func processFile(filePath string, csvWriter *csv.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading line from file %s: %w", filePath, err)
		}
		line = strings.TrimSpace(line)
		if line != "" {
			err = processLogLine(line, csvWriter)
			if err != nil {
				return fmt.Errorf("error processing line in file %s: %w", filePath, err)
			}
			csvWriter.Flush()
		}
	}
	return nil
}

// getLastMinuteFiles encuentra y procesa los archivos de log modificados en el último minuto.
func getLastMinuteFiles(dirPath string, csvFilePath string) error {
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %w", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	header := []string{"Host", "IP", "Estado"}
	err = csvWriter.Write(header)
	if err != nil {
		return fmt.Errorf("error writing header to CSV: %w", err)
	}
	csvWriter.Flush()
	if err != nil {
		return fmt.Errorf("error flushing csv writer: %w", err)
	}

	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".log.old") && info.ModTime().After(oneMinuteAgo) {
			fmt.Printf("Procesando archivo: %s\n", path)
			err := processFile(path, csvWriter)
			if err != nil {
				return fmt.Errorf("error processing file %s: %w", path, err)
			}
		}

		return nil
	})
}

func main() {
	logDir := "logs.log.old" // Reemplaza con la ruta a tu directorio de logs
	csvFile := "output.csv"

	// Crear el directorio si no existe
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			fmt.Printf("Error al crear el directorio: %v\n", err)
			return
		}
	}

	err := getLastMinuteFiles(logDir, csvFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Proceso completado. El archivo CSV se ha creado exitosamente.")
	}
}
