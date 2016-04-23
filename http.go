// GoChar project http.go
package main

import (
	"fmt"
	"image"
	//"io"
	//"mime/multipart"
	"net/http"
	//"os"
	"strconv"
)

import _ "image/png"
import _ "image/gif"
import _ "image/jpeg"

// Nodo
type Nodo struct {
	ip        string
	idtrabajo int
}

func (n *Nodo) AsginarTrabajo(id int, imagen image.Image) bool {
	if n.idtrabajo == -1 {
		n.idtrabajo = id
		// Call RPC

		return true
	} else {
		return false
	}
}

var id int = 0

//nodos list

func handler_subir(w http.ResponseWriter, r *http.Request) {
	// Recibir archivo
	id++

	//archivo, header, err := r.FormFile("archivo")
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error del servidor:")
		fmt.Fprintln(w, err)
		return
	}

	// Pasarle la image a otro método
	// Quien?
	var nodo Nodo
	if nodo.AsginarTrabajo(id, imagen) == true {
		fmt.Fprint(w, id)
	} else {
		fmt.Fprint(w, "-1")
	}
}

func handler_estado(w http.ResponseWriter, r *http.Request) {
	identificador, err := strconv.Atoi(r.URL.Query()["id"][0])
	if err == nil {
		fmt.Fprintf(w, "%d", identificador)
		// fmt.Fprintf(w, "%d", mapatrabajos[identificador])
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

	//sd
}

func main() {
	/*	URLs:
	*	- Pedir foto: /subir
	*	- Preguntar por foto: /estado
	*	- Estadísticas: /estadisticas
	 */
	http.HandleFunc("/subir", handler_subir)
	http.HandleFunc("/estado", handler_estado)
	http.HandleFunc("/estadisticas", handler_estadisticas)
	http.ListenAndServe(":80", nil)
}
