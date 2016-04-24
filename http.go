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
	"log"
	"strconv"

	"github.com/valyala/gorpc"
)

import _ "image/png"
import _ "image/gif"
import _ "image/jpeg"

// Nodo
type Nodo struct {
	idtrabajo   int
	ip          string
	cliente     *gorpc.Client
	clienteDisp *gorpc.DispatcherClient //*?
	resultado   byte
}

// Globales
var cuentaTrabajos int = 0

var nodos *list.List
var indexRobin int = 0

// Servidor RPC
type exportaServer struct{}

//ServerAddr es la direccion del servidor que levanta el esclavo
func (s *exportaServer) AceptaConexiones(clientAddr string, sinUso int) {
	//Añado un nuevo nodo a la lista
	//Creo un cliente apuntando al servidor del nodo

	d := gorpc.NewDispatcher()
	c := gorpc.NewTCPClient(clientAddr)
	c.Start() //->Deberia arrancar el cliente? puede que para mas adelante?
	dc := d.NewServiceClient("servicioRPCEsclavo", c)

	n := Nodo{-1, clientAddr, c, dc, 0} //Creo el nodo y lo inicializo.

	nodos.PushBack(n)

}

func (s *exportaServer) CierraConexiones(clientAddr string, sinUso int) {
	//Busco nodo en la lista/map y hago un .Remove sobre el.
	var nodo Nodo
	for e := nodos.Front(); e != nil; e = e.Next() {
		nodo = e.Value.(Nodo)
		if nodo.ip == clientAddr {
			nodo.cliente.Stop()
			nodos.Remove(e)
		}
	}
}

func (s *exportaServer) RecibeRespuesta(clientAddr string, caracter_final byte) {

	//No se que carajo haces con la respuesta y el http handler
	var nodo Nodo
	for e := nodos.Front(); e != nil; e = e.Next() {
		nodo = e.Value.(Nodo)
		if nodo.ip == clientAddr {
			nodo.resultado = caracter_final
		}
	}

}

// -------------
// Servidor HTTP

func (n *Nodo) AsignarTrabajo(id int, imagen *image.YCbCr) bool {
	if n.idtrabajo == -1 {
		// Call RPC
		res, err := n.clienteDisp.Call("RecibeImagen", imagen)
		fmt.Print(res)
		if err != nil {
			n.idtrabajo = cuentaTrabajos
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

func handler_subir(w http.ResponseWriter, r *http.Request) {
	// Recibir archivo
	cuentaTrabajos++

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
	if nodo.AsignarTrabajo(cuentaTrabajos, img) == true {
		indexRobin++
		fmt.Fprint(w, cuentaTrabajos)
	} else {
		fmt.Fprint(w, "-1")
	}
}

func handler_estado(w http.ResponseWriter, r *http.Request) {
	identificador, err := strconv.Atoi(r.URL.Query()["id"][0])
	if err == nil {
		var nodo Nodo
		for e := nodos.Front(); e != nil; e = e.Next() {
			nodo = e.Value.(Nodo)
			if nodo.idtrabajo == identificador {
				if nodo.resultado == 0 {
					fmt.Fprint(w, "0")
				} else {
					fmt.Fprintf(w, "%c", nodo.resultado)

				}
			}
		}
	} else {

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

	addr := "DIRECCION SERVIDOR RPC"
	nodos = list.New() //Lista enlazada de objetos nodo

	// Registrar tipo *image.YCbCr para RPC
	var img *image.YCbCr
	gorpc.RegisterType(img)

	serverDispatcher := gorpc.NewDispatcher()
	service := &exportaServer{}
	serverDispatcher.AddService("goChar", service)
	rpcServer := gorpc.NewTCPServer(addr, serverDispatcher.NewHandlerFunc())

	if err := rpcServer.Start(); err != nil {
		log.Fatalf("No puedo arrancar el servidor: [%s]", err)
	}
	defer rpcServer.Stop()

	//Servidor RPC arrancado en este punto

	http.HandleFunc("/subir", handler_subir)
	http.HandleFunc("/estado", handler_estado)
	http.HandleFunc("/estadisticas", handler_estadisticas)
	http.ListenAndServe(":80", nil)
}
