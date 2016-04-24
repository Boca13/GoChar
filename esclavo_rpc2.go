package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"net"
	"os"
	//"os/exec"
	"os/signal"
	"strconv"

	"github.com/cenkalti/rpc2"
	"golang.org/x/image/bmp"
)

//Globales
var conn, _ = net.Dial("tcp", "192.168.1.10:12345")
var clt *rpc2.Client
var identificador_nodo int

type Args_RecibeImagen struct {
	Imagen *image.YCbCr
}
type Args_Conexiones int
type Reply_RecibeImagen bool

type Args_RecibeRespuesta struct {
	Id        int
	Resultado byte
}

func RecibeImagen(client *rpc2.Client, args *Args_RecibeImagen, reply *Reply_RecibeImagen) error {
	log.Printf("Me ha llegado una imagen.")

	writer, err := os.Create("char.bmp")
	if err != nil {
		log.Fatal(err)
	}
	bmp.Encode(writer, args.Imagen)
	writer.Close()
	log.Printf("Imagen convertida a archivo .bmp correctamente")

	log.Printf("Comenzando deteccion de caracter...")

					//LLAMADA A PYTHON

	/*
		// Llamar al python

			cmd := exec.Command("python", "reconocer.py", "char.bmp")
			err = cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Wait()
			log.Printf("Resultado: %v", err)
	*/
	err = errors.New("97")


					//FIN LLAMADA A PYTHON
	resultado, err := strconv.Atoi(err.Error())
	if err != nil {
		log.Fatal(err)
	}
	*reply = true

	log.Printf("Deteccion de caracter finalizada. Enviando respuesta...")
	//_, err = dc.Call("RecibeRespuesta", resultado) //Envio el caracter con el resultado al maestro.

	var final = Args_RecibeRespuesta{identificador_nodo, byte(resultado)}

	err = clt.Call("RecibeRespuesta", &final, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Respuesta de caracter enviada.")

	return nil
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			err := clt.Call("CierraConexiones", identificador_nodo, identificador_nodo)
			log.Printf("Mi id_nodo antes de salir es: %d", identificador_nodo)
			if err != nil {
				log.Fatal(err)
			}

			os.Exit(1)
		}
	}()

	//Funcion que recibe las imagenes.
	log.Print("Esclavo de goRpc iniciado.")
	log.Print("Conexion a servidor RPC....")
	clt = rpc2.NewClient(conn)

	clt.Handle("RecibeImagen", RecibeImagen)

	log.Printf("Conectado correctamente a servidorRPC.")
	//Cuando ejecuto el cliente, llamo al servidor por aqui.

	var conexion Args_Conexiones = 3
	go clt.Run()                                                       //Se crea en otro hilo
	err := clt.Call("AceptaConexiones", conexion, &identificador_nodo) //Me registro en servidor.
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Registrado en servidorRPC correctamente.")
	//Programar ctrl-c para desconectar.
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt)

	fmt.Scanln()

}
