package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// filepath: c:\Users\Ismenia Guevara\Documents\Repositorios GIT\PracticasComputacionNube\entrega3\AlgoritmoDisparador.go
// Función para obtener la ruta del archivo
func obtenerRutaArchivo() (string, error) {
	directorioActual, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error al obtener el directorio actual: %w", err)
	}
	return filepath.Join(directorioActual, "archivo.txt"), nil
}

// Función para leer el estado del archivo
func obtenerestado(archivoParam string) (string, error) {
	contenido, err := ioutil.ReadFile(archivoParam)
	if err != nil {
		return "", fmt.Errorf("error al leer el archivo: %w", err)
	}

	estado := string(contenido)

	// Elimina espacios en blanco y saltos de línea
	estadoLimpia := ""
	for _, r := range estado {
		if r != ' ' && r != '\n' && r != '\r' {
			estadoLimpia += string(r)
		}
	}
	estado = estadoLimpia

	return estado, nil
}

// Función para determinar el mensaje según la estado
func obtenerMensaje(estado string) string {
	switch estado {
	case "A":
		return "A eliminación de 1 servidor del balanceador"
	case "B":
		return "B lanzamiento de 1 nuevo servidor"
	case "C":
		return "C estado óptimo de trabajo del servidor"
	default:
		return "Estado desconocido" // Manejo de estado no válida
	}
}

func main1() {
	// Obtiene la ruta del archivo
	archivoParam, err := obtenerRutaArchivo()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Variable para almacenar la estado leída
	var estadoAlmacenado string

	// Bucle infinito para ejecución constante
	for {
		// Llama a la función para obtener la estado
		estado, err := obtenerestado(archivoParam)
		if err != nil {
			fmt.Println(err)
			// Espera un poco antes de volver a intentar
			time.Sleep(5 * time.Second)
			continue // Vuelve al inicio del bucle
		}

		// Valida que la estado sea A, B o C
		if estado != "A" && estado != "B" && estado != "C" {
			fmt.Println("En espera de un estado...")
			// Espera un poco antes de volver a intentar
			time.Sleep(5 * time.Second)
			continue // Vuelve al inicio del bucle
		}

		// Define el mensaje según la estado
		mensaje := obtenerMensaje(estado)

		// Almacena la estado en la variable
		estadoAlmacenado = estado

		fmt.Println("estado almacenado:", estadoAlmacenado)
		fmt.Println(mensaje)

		// Espera un tiempo antes de la siguiente iteración
		time.Sleep(5 * time.Second)
	}
}
