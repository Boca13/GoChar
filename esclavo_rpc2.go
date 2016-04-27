package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	

	"github.com/cenkalti/rpc2"
)

//Globales
var conn, _ = net.Dial("tcp", "127.0.0.1:12345")
var clt *rpc2.Client
var identificador_nodo int
var identif_desconex int

type Args_RecibeImagen struct {
	Imagen *image.Gray
}
type Args_Conexiones int
type Reply_RecibeImagen bool


type Args_RecibeRespuesta struct {
	Id        int
	Resultado byte
}

func RecibeImagen(client *rpc2.Client, args *Args_RecibeImagen, reply *Reply_RecibeImagen) error {
	log.Printf("Me ha llegado una imagen.")

	writer, err := os.Create("char.png")
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(writer, args.Imagen)
	writer.Close()
	log.Printf("Imagen convertida a archivo .png correctamente")

	log.Printf("Comenzando deteccion de caracter...")

	//LLAMADA A PYTHON

	// Llamar al python

	cmd := exec.Command("python", "4c_identificar.py", "char.png")
	resultado, err := cmd.Output()
	
	
	log.Printf("Resultado: %s", resultado)
	//err = errors.New("97")

	//FIN LLAMADA A PYTHON
	
	if err != nil {
		log.Printf("El error que suelta python es  != de nil")
		log.Fatal(err)
	}
	*reply = true

	log.Printf("Deteccion de caracter finalizada. Enviando respuesta...")
	//_, err = dc.Call("RecibeRespuesta", resultado) //Envio el caracter con el resultado al maestro.

	var final = Args_RecibeRespuesta{identificador_nodo, byte(resultado[0])}

	err = clt.Call("RecibeRespuesta", &final, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Respuesta de caracter enviada.")

	return nil
}

func main() {
	var conexiones Args_Conexiones = 1 //Bandera que envio al servidor
	

	//Funcion que recibe las imagenes.
	log.Print("Esclavo de goRpc iniciado.")
	log.Print("Conexion a servidor RPC....")
	clt = rpc2.NewClient(conn)

	clt.Handle("RecibeImagen", RecibeImagen)

	log.Printf("Conectado correctamente a servidorRPC.")
	//Cuando ejecuto el cliente, llamo al servidor por aqui.

	go clt.Run()                                                         //Se crea en otro hilo
	err := clt.Call("AceptaConexiones", conexiones, &identificador_nodo) //Me registro en servidor.
		//fmt.Println("Identificador nodo (conexion): ", identificador_nodo) DEBUG
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Registrado en servidorRPC correctamente.")
	//Programar ctrl-c para desconectar.
	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt)

	fmt.Scanln()

	log.Printf("Cerrando programa, desconexion del servidor....")
	//Comprobar esta llamada.
	err = clt.Call("CierraConexiones", identificador_nodo, &identif_desconex)
		//fmt.Println("Identificador nodo(desconexion): ",identif_desconex) DEBUG
	
	if identif_desconex == -2{
		fmt.Println("Desconexion correcta del servidor.")
		clt.Close()
	}

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)

}
