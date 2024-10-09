package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// Producto representa un artículo en la factura
type Producto struct {
	Nombre   string
	Cantidad int
	Precio   float64
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Solicitar el nombre del cliente
	fmt.Print("Ingrese el nombre del cliente: ")
	nombreCliente, _ := reader.ReadString('\n')
	nombreCliente = strings.TrimSpace(nombreCliente)

	// Obtener la fecha actual automáticamente
	fechaFactura := time.Now().Format("02/01/2006")

	// Lista de productos
	var productos []Producto
	for {
		// Ingresar datos del producto
		fmt.Print("Ingrese el nombre del producto (o 'fin' para terminar): ")
		nombreProducto, _ := reader.ReadString('\n')
		nombreProducto = strings.TrimSpace(nombreProducto)

		if strings.ToLower(nombreProducto) == "fin" {
			break
		}

		fmt.Print("Ingrese la cantidad: ")
		cantidadStr, _ := reader.ReadString('\n')
		cantidad, err := strconv.Atoi(strings.TrimSpace(cantidadStr))
		if err != nil {
			fmt.Println("Cantidad no válida.")
			continue
		}

		fmt.Print("Ingrese el precio unitario: ")
		precioStr, _ := reader.ReadString('\n')
		precio, err := strconv.ParseFloat(strings.TrimSpace(precioStr), 64)
		if err != nil {
			fmt.Println("Precio no válido.")
			continue
		}

		// Agregar producto a la lista
		productos = append(productos, Producto{
			Nombre:   nombreProducto,
			Cantidad: cantidad,
			Precio:   precio,
		})
	}

	// Crear el PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Título
	pdf.Cell(40, 10, "Factura")

	// Información del cliente y la fecha
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Cliente: %s", nombreCliente))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Fecha: %s", fechaFactura))
	pdf.Ln(12)

	// Encabezados de la tabla con líneas de separación
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(70, 10, "Producto", "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 10, "Cantidad", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Precio Unitario", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Total", "1", 0, "C", false, 0, "")
	pdf.Ln(10)

	// Contenido de la tabla con líneas de separación
	pdf.SetFont("Arial", "", 12)
	var totalGeneral float64
	for _, producto := range productos {
		totalProducto := float64(producto.Cantidad) * producto.Precio
		totalGeneral += totalProducto

		pdf.CellFormat(70, 10, producto.Nombre, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", producto.Cantidad), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("$%.2f", producto.Precio), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("$%.2f", totalProducto), "1", 0, "C", false, 0, "")
		pdf.Ln(10)
	}

	// Total general con líneas de separación
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(10)
	pdf.CellFormat(70, 10, "", "", 0, "", false, 0, "")
	pdf.CellFormat(30, 10, "", "", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, "Total General:", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("$%.2f", totalGeneral), "1", 0, "C", false, 0, "")

	// Guardar el archivo PDF
	err := pdf.OutputFileAndClose("factura.pdf")
	if err != nil {
		fmt.Println("Error al generar el PDF:", err)
	} else {
		fmt.Println("Factura generada exitosamente: factura.pdf")
	}
}
