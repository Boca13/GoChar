// GoChar project http.go
package main

import (
	"fmt"
	"image"
	"net/http"
	"container/list"
	"image/color"
	"log"
	"math"
	"net"
	"strconv"

	"github.com/cenkalti/rpc2"
	"github.com/nfnt/resize"
)

import _ "image/png"
import _ "image/gif"
import _ "image/jpeg"

// Nodo
type Nodo struct {
	idTrabajo int
	idNodo    int
	cliente   *rpc2.Client
	resultado byte
}

// Globales
var cuentaTrabajos int = 0
var cuentaNodos int = 0
var nodos *list.List
var indexRobin int = 0

// Tipos para IO de funciones RPC
type Args_Conexiones int
type Reply_Conexiones int

type Args_RecibeRespuesta struct {
	Id        int
	Resultado byte
}
type Reply_RecibeRespuesta bool
type Args_RecibeImagen struct {
	Imagen *image.Gray
}

func AceptaConexiones(client *rpc2.Client, args *Args_Conexiones, reply *Reply_Conexiones) error {
	//Añado un nuevo nodo a la lista
	//Creo un cliente apuntando al servidor del nodo
	n := Nodo{-1, cuentaNodos, client, 0} //Creo el nodo y lo inicializo.
	*reply = Reply_Conexiones(cuentaNodos)
	nodos.PushBack(&n)
	log.Println("Conectado cliente con id: ", cuentaNodos)
	cuentaNodos++
	return nil
}

func CierraConexiones(client *rpc2.Client, args *Args_Conexiones, reply *Reply_Conexiones) error {
	//Busco nodo en la lista/map y hago un .Remove sobre el.
	var nodo *Nodo
	for e := nodos.Front(); e != nil; e = e.Next() {
		nodo = e.Value.(*Nodo)
		if nodo.idNodo == int(*args) {
			*reply = -2
			//nodo.cliente.Close()
			nodos.Remove(e)			
			log.Println("Desconectado nodo ", nodo.idNodo)
			
		}
	}
	return nil
}

func RecibeRespuesta(client *rpc2.Client, args *Args_RecibeRespuesta, reply *Reply_RecibeRespuesta) error {
	var nodo *Nodo
	for e := nodos.Front(); e != nil; e = e.Next() {
		nodo = e.Value.(*Nodo)
		if nodo.idNodo == args.Id {
			nodo.resultado = args.Resultado
		}
	}
	return nil
}

// -------------
// Servidor HTTP

func (n *Nodo) AsignarTrabajo(id int, imagen *image.Gray) bool {
	if n.idTrabajo == -1 {
		// Call RPC
		var res Reply_RecibeRespuesta
		log.Println("Enviando imagen al nodo por RPC...")
		err := n.cliente.Call("RecibeImagen", Args_RecibeImagen{imagen}, &res)
		log.Println("Respuesta: ", res)
		if err == nil {
			n.idTrabajo = cuentaTrabajos
			log.Println("Asignando id de trabajo ", n.idTrabajo)
			return true
		} else {
			log.Println("Error al pasar la imagen a esclavo", err)
			return false
		}
	} else {
		return false
	}
}

func handler_subir(w http.ResponseWriter, r *http.Request) {
	// Recibir archivo
	cuentaTrabajos++
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	reader, err := r.MultipartReader()
	log.Println("Recibiendo archivo...")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error en la consulta:")
		fmt.Fprintln(w, err)
		return
	}

	parte, err := reader.NextPart()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error del servidor:")
		fmt.Fprintln(w, err)
		return
	}
	log.Println("Decodificando imagen...")
	img, _, err := image.Decode(parte)
	// Redimensionar
	img = resize.Resize(28, 28, img, resize.Bicubic)

	// Create a new grayscale image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			oldColor := img.At(x, y)
			r, g, b, _ := oldColor.RGBA()
			avg := 0.2125*float64(r) + 0.7154*float64(g) + 0.0721*float64(b)
			grayColor := color.Gray{uint8(math.Ceil(avg))}
			gray.Set(x, y, grayColor)
		}
	}
	// ------

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error del servidor:")
		fmt.Fprintln(w, err)
		return
	}
	log.Println("Imagen decodificada.")
	// Pasarle la image a otro método
	// Decidir a quien. Round-robin
	e := nodos.Front()
	for i := 0; (e != nil) && (i < indexRobin) && (e.Value.(*Nodo).idTrabajo != -1); e = e.Next() {
	}
	if (e == nil) || (e.Value.(*Nodo).idTrabajo != -1) {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Todos los nodos están ocupados")
		indexRobin = 0
		return
	}
	var nodo *Nodo = e.Value.(*Nodo)

	log.Println("Trabajo asignado al nodo ", nodo.idNodo)

	//nodo.idTrabajo = -1
	if nodo.AsignarTrabajo(cuentaTrabajos, gray) == true {
		log.Println("Trabajo número ", cuentaTrabajos, " está asignado al nodo", nodo.idNodo)
		indexRobin++
		fmt.Fprint(w, cuentaTrabajos)
	} else {
		fmt.Fprint(w, "-1")
	}
}

func handler_estado(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	identificador, err := strconv.Atoi(r.URL.Query()["id"][0])
	log.Println("Petición de estado del trabajo ", identificador)
	if err == nil {
		var nodo *Nodo
		for e := nodos.Front(); e != nil; e = e.Next() {
			nodo = e.Value.(*Nodo)
			if nodo.idTrabajo == identificador {
				log.Println(" Trabajo encargado al nodo ", nodo.idNodo)
				if nodo.resultado == 0 {
					log.Println("  -> El trabajo no se ha terminado")
					fmt.Fprint(w, "0")
				} else {
					log.Println("  -> El nodo ha devuelto el valor ", nodo.resultado)
					fmt.Fprintf(w, "%c", nodo.resultado)
					nodo.idTrabajo = -1
					nodo.resultado = 0
				}
			}
		}
	} else {

	}
}

func main() {
	/*	URLs:
	*	- Pedir foto: /subir
	*	- Preguntar por foto: /estado
	*/

	nodos = list.New() //Lista enlazada de objetos nodo

	server := rpc2.NewServer()
	server.Handle("AceptaConexiones", AceptaConexiones)
	server.Handle("CierraConexiones", CierraConexiones)
	server.Handle("RecibeRespuesta", RecibeRespuesta)

	listener, err := net.Listen("tcp", "0.0.0.0:12345")
	go server.Accept(listener)
	if err != nil {
		log.Fatalf("No puedo arrancar el servidor: [%s]", err)
	}
	// Servidor RPC arrancado en este punto

	// Ahora servidor HTTP
	http.HandleFunc("/subir", handler_subir)
	http.HandleFunc("/estado", handler_estado)
	http.ListenAndServe(":80", nil)
}
