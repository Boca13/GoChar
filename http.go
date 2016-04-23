// GoChar project http.go
package main

import (
	"fmt"
	"image"
	//"io"
	//"mime/multipart"
	"net/http"
	//"os"
	"container/list"
	"strconv"

	"github.com/valyala/gorpc"
)

import _ "image/png"
import _ "image/gif"
import _ "image/jpeg"

// Nodo
type Nodo struct {
	//ip        string
	idtrabajo int
	cliente   *gorpc.DispatcherClient
	resultado byte
}

func (n *Nodo) AsignarTrabajo(id int, imagen *image.YCbCr) bool {
	if n.idtrabajo == -1 {
		// Call RPC
		res, err := n.cliente.Call("RecibeImagen", imagen)
		fmt.Print(res)
		if err != nil {
			n.idtrabajo = id
			n.resultado = 0
			//i := imagen.(image.YCbCr)
			fmt.Print(imagen)
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

var id int = 0

var nodos *list.List
var indexRobin int = 0

func handler_subir(w http.ResponseWriter, r *http.Request) {
	// Recibir archivo
	id++

	reader, err := r.MultipartReader()

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

	imagen, _, err := image.Decode(parte)
	img := imagen.(*image.YCbCr)
	otraimg := image.Image(img)
	fmt.Print(otraimg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error del servidor:")
		fmt.Fprintln(w, err)
		return
	}

	// Pasarle la image a otro método
	// Decidir a quien. Round-robin
	e := nodos.Front()
	for i := 0; (e != nil) && (i < indexRobin); e = e.Next() {
	}
	if e == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, "Todos los nodos están ocupados")
		indexRobin = 0
		return
	}
	var nodo Nodo = e.Value.(Nodo)

	nodo.idtrabajo = -1
	if nodo.AsignarTrabajo(id, img) == true {
		indexRobin++
		fmt.Fprint(w, id)
	} else {
		fmt.Fprint(w, "-1")
	}
}

func handler_estado(w http.ResponseWriter, r *http.Request) {
	identificador, err := strconv.Atoi(r.URL.Query()["id"][0])
	if err == nil {
		fmt.Fprintf(w, "%d", identificador)
	} else {
		fmt.Fprintf(w, "En marcha", r.URL.Path[1:])
	}
}

func handler_estadisticas(w http.ResponseWriter, r *http.Request) {
	identificador, err := strconv.Atoi(r.URL.Query()["id"][0])
	if err == nil {
		fmt.Fprintf(w, "%d", identificador)
	} else {
		fmt.Fprintf(w, "En marcha", r.URL.Path[1:])
	}

}

func main() {
	/*	URLs:
	*	- Pedir foto: /subir
	*	- Preguntar por foto: /estado
	*	- Estadísticas: /estadisticas
	 */

	nodos = list.New()

	http.HandleFunc("/subir", handler_subir)
	http.HandleFunc("/estado", handler_estado)
	http.HandleFunc("/estadisticas", handler_estadisticas)
	http.ListenAndServe(":80", nil)
}
