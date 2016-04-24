package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/valyala/gorpc"
)

type dispatcherEsclavo struct{}

func main() {

	d := gorpc.NewDispatcher()

	service := &dispatcherEsclavo{}

	d.AddService("servicioRPCEsclavo", service)

	//Arranco el servidor rpc
	addr := "0.0.0.0:12345"
	s := gorpc.NewTCPServer(addr, d.NewHandlerFunc())
	if err := s.Start(); err != nil {
		log.Fatalf("Cannot start rpc server: [%s]", err)
	}
	defer s.Stop()

	// Aqui debo poner la direccion del servidor RPC maestro
	c := gorpc.NewTCPClient("192.168.1.10:12345")
	c.Start()
	defer c.Stop()

	dc := d.NewServiceClient("goChar", c)

	res, err := dc.Call("AceptaConexiones", addr) //Le mando la direccion de mi servidor RPC

	fmt.Scanln()

}