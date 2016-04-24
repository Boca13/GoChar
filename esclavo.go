package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"

	"github.com/valyala/gorpc"
	"golang.org/x/image/bmp"
)

var dc *gorpc.DispatcherClient
var direccionMaestro string = "192.168.1.10:12345" //direccion server maestro RPC
var addr string = "0.0.0.0:12345"                  //direccion server esclavo RPC

type dispatcherEsclavo struct{}

//ServerAddr es la direccion del servidor que levanta el esclavo
func (s *dispatcherEsclavo) RecibeImagen(clientAddr string, imagen *image.YCbCr) {
	log.Printf("Llega una petición:")
	writer, err := os.Create("char.bmp")
	if err != nil {
		log.Fatal(err)
	}
	bmp.Encode(writer, imagen)
	writer.Close()

	// Llamar al python
	cmd := exec.Command("python", "reconocer.py", "char.bmp")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	log.Printf("Resultado: %v", err)
	resultado, err := strconv.Atoi(err.Error())
	if err != nil {
		log.Fatal(err)
	}
	res, err := dc.Call("RecibeRespuesta", resultado) //Envio el caracter con el resultado al maestro.
}

func main() {
	// Registrar tipo *image.YCbCr para RPC
	var img *image.YCbCr
	gorpc.RegisterType(img)

	d := gorpc.NewDispatcher()

	service := &dispatcherEsclavo{}

	d.AddService("servicioRPCEsclavo", service)

	//Arranco el servidor rpc

	s := gorpc.NewTCPServer(addr, d.NewHandlerFunc())
	if err := s.Start(); err != nil {
		log.Fatalf("Cannot start rpc server: [%s]", err)
	}

	// Aqui debo poner la direccion del servidor RPC maestro
	c := gorpc.NewTCPClient(direccionMaestro)
	c.Start()

	dc = d.NewServiceClient("goChar", c)

	res, err := dc.Call("AceptaConexiones", 3) //Mando el int sin uso

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt)
	go func() {
		for sig := range cc {
			// sig is a ^C, handle it
			res, err := dc.Call("CierraConexiones", 3) //Mando el int sin uso
			c.Stop()
			s.Stop()

		}
	}()

	fmt.Scanln()

}
