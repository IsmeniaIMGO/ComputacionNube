package main //indica que este es un paquete ejecutable

import (
	"encoding/base64" //proporciona funciones para codificar y decodificar datos en base64
	"fmt"             //biblioteca estándar de Go y proporciona funciones de formato de entrada y salida,
	"io/ioutil"       //proporciona funciones para leer y escribir archivos
	"math/rand"       //proporciona funciones para generar números aleatorios
	"os"              //proporciona funciones de sistema operativo independientes de la plataforma
	"strings"         //proporciona funciones para manipular cadenas de texto
	"time"            //proporciona funciones para medir y mostrar el tiempo

	"github.com/sqweek/dialog" //proporciona funciones para mostrar cuadros de diálogo de archivos y directorios
)

func main() {

	//Menu de opciones
	for {
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Saludar")
		fmt.Println("2. Ingresar un texto, un número entero y mostrarlos")
		fmt.Println("3. Mostrar el nombre del host")
		fmt.Println("4. Mostrar nombres de archivos en una carpeta")
		fmt.Println("5. Mostrar nombres de imágenes en una carpeta")
		fmt.Println("6. Contar archivos e imágenes en una carpeta")
		fmt.Println("7. Seleccionar una imagen aleatoria")
		fmt.Println("8. Codificar una imagen aleatoria en Base64")
		fmt.Println("9. Codificar cuatro imágenes aleatorias en Base64 y guardarlas en archivos de texto")
		fmt.Println("0. Salir")

		var opcion int
		fmt.Print("Ingrese el número de la opción: ")
		fmt.Scanln(&opcion)

		switch opcion {
		case 1:
			punto1()
		case 2:
			punto2()
		case 3:
			punto3()
		case 4:
			punto4()
		case 5:
			punto5()
		case 6:
			punto6()
		case 7:
			punto7()
		case 8:
			punto8()
		case 9:
			punto9()
		case 0:
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
	}

}

func seleccionarCarpeta() (string, []os.FileInfo) {
	rutaCarpeta, err := dialog.Directory().Title("Seleccione una carpeta").Browse()
	if err != nil {
		fmt.Println("Error seleccionando la carpeta:", err)
		return "", nil
	}

	files, err := ioutil.ReadDir(rutaCarpeta)
	if err != nil {
		fmt.Println("Error leyendo la carpeta:", err)
		return rutaCarpeta, nil
	}

	return rutaCarpeta, files
}

func filtrarImagenes(files []os.FileInfo) []string {
	var imagenes []string
	for _, file := range files {
		if !file.IsDir() {
			nombreArchivo := file.Name()
			if strings.HasSuffix(strings.ToLower(nombreArchivo), ".jpg") ||
				strings.HasSuffix(strings.ToLower(nombreArchivo), ".jpeg") ||
				strings.HasSuffix(strings.ToLower(nombreArchivo), ".png") {
				imagenes = append(imagenes, nombreArchivo)
			}
		}
	}
	return imagenes
}

func seleccionarImagenesAleatorias(cantidad int) []string {

	rutaCarpeta, files := seleccionarCarpeta()
	if rutaCarpeta == "" || files == nil {
		return nil
	}

	imagenes := filtrarImagenes(files)
	if len(imagenes) < cantidad {
		fmt.Printf("No hay suficientes imágenes en la carpeta. Se encontraron %d imágenes.\n", len(imagenes))
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	var imagenesSeleccionadas []string
	for i := 0; i < cantidad; i++ {
		indice := rand.Intn(len(imagenes))
		imagenesSeleccionadas = append(imagenesSeleccionadas, rutaCarpeta+"\\"+imagenes[indice])
		imagenes = append(imagenes[:indice], imagenes[indice+1:]...)
	}

	return imagenesSeleccionadas
}

func punto1() {

	fmt.Println("Holaaaaaa")

}

func punto2() {
	var texto string
	var numero int

	// Solicita y almacena un argumento de tipo string
	fmt.Print("Ingrese un texto: ")
	fmt.Scanln(&texto)

	// Solicita y valida que el segundo argumento sea un entero
	for {
		fmt.Print("Ingrese un número entero: ")
		_, err := fmt.Scanln(&numero)
		if err == nil {
			break
		}
		fmt.Println("Entrada no válida. Por favor, ingrese un número entero.")
	}

	// Muestra el texto y el número ingresados
	fmt.Printf("Texto ingresado: %s\n", texto)
	fmt.Printf("Número ingresado: %d\n", numero)

	/*
		%s: Formatea un valor como una cadena de texto.
		%d: Formatea un valor como un número entero.
		%f: Formatea un valor como un número de punto flotante.
		%t: Formatea un valor como un booleano (true o false).
		%v: Formatea un valor en su representación predeterminada.
	*/

}

func punto3() {

	// Obtiene y muestra el nombre del host
	hostname, verificacion2 := os.Hostname()
	if verificacion2 != nil {
		fmt.Println("Error obteniendo el nombre del host:", verificacion2)
		return
	}
	fmt.Printf("Hostname: %s\n", hostname)
}

func punto4() {
	rutaCarpeta, files := seleccionarCarpeta()
	if rutaCarpeta == "" || files == nil {
		return
	}

	fmt.Println("Archivos en la carpeta:", rutaCarpeta)
	for _, file := range files {
		fmt.Println(file.Name())
	}
}

func punto5() {
	rutaCarpeta, files := seleccionarCarpeta()
	if rutaCarpeta == "" || files == nil {
		return
	}

	imagenes := filtrarImagenes(files)
	fmt.Println("Imágenes en la carpeta:", rutaCarpeta)
	for _, imagen := range imagenes {
		fmt.Println(imagen)
	}
}

func punto6() {
	rutaCarpeta, files := seleccionarCarpeta()
	if rutaCarpeta == "" || files == nil {
		return
	}

	imagenes := filtrarImagenes(files)
	totalArchivos := len(files)

	fmt.Printf("Total de archivos en la carpeta: %d\n", totalArchivos)
	fmt.Printf("Total de imágenes en la carpeta: %d\n", len(imagenes))
	fmt.Println("Nombres de las imágenes:")
	for _, imagen := range imagenes {
		fmt.Println(imagen)
	}
}

func punto7() string {
	rutaCarpeta, files := seleccionarCarpeta()
	if rutaCarpeta == "" || files == nil {
		return ""
	}

	imagenes := filtrarImagenes(files)
	if len(imagenes) == 0 {
		fmt.Println("No se encontraron imágenes en la carpeta.")
		return ""
	}

	// Selecciona una imagen de manera aleatoria
	rand.Seed(time.Now().UnixNano())
	imagenAleatoria := imagenes[rand.Intn(len(imagenes))]

	fmt.Printf("Imagen seleccionada aleatoriamente: %s\n", imagenAleatoria)
	return rutaCarpeta + "\\" + imagenAleatoria
}

func punto8() {
	imagenSeleccionada := punto7()
	if imagenSeleccionada == "" {
		return
	}

	// Lee el contenido de la imagen seleccionada
	contenido, err := ioutil.ReadFile(imagenSeleccionada)
	if err != nil {
		fmt.Println("Error leyendo la imagen:", err)
		return
	}

	// Codifica el contenido en Base64
	codificado := base64.StdEncoding.EncodeToString(contenido)

	// Muestra el contenido codificado en la terminal
	fmt.Println("Contenido de la imagen en Base64:")
	fmt.Println(codificado)
}

func punto9() {
	imagenesSeleccionadas := seleccionarImagenesAleatorias(4)
	if imagenesSeleccionadas == nil {
		return
	}

	for _, imagenSeleccionada := range imagenesSeleccionadas {
		// Lee el contenido de la imagen seleccionada
		contenido, err := ioutil.ReadFile(imagenSeleccionada)
		if err != nil {
			fmt.Println("Error leyendo la imagen:", err)
			continue
		}

		// Codifica el contenido en Base64
		codificado := base64.StdEncoding.EncodeToString(contenido)

		// Guarda el contenido codificado en un archivo de texto
		nombreArchivoB64 := imagenSeleccionada + ".b64"
		err = ioutil.WriteFile(nombreArchivoB64, []byte(codificado), 0644)
		if err != nil {
			fmt.Println("Error guardando el archivo:", err)
			continue
		}

		fmt.Printf("Imagen %s codificada y guardada como %s\n", imagenSeleccionada, nombreArchivoB64)
	}
}
