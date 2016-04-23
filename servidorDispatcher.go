package main

import (
    "errors"
    "fmt"
    "log"
    "github.com/valyala/gorpc"
)

type Nodo struct {
	//ip        string
	idtrabajo int
	cliente 		  gorpc.Client
	clienteDispatcher gorpc.DispatcherClients
	resultado byte
}



type exportaServer struct{

}

	//ServerAddr es la direccion del servidor que lvanta el esclavo

func (s *exportaServer) AceptaConexiones(serverAddr string){
	//Añado un nuevo nodo a la lista/map
	//Creo un cliente apuntando al servidor del nodo

	c := gorpc.NewTCPClient(serverAddr)
	c.Start() //->Deberia arrancar el cliente? puede que para mas adelante?
	dc := d.NewServiceClient("servicioRPCEsclavo",c)

	//dc.Call(.............)
	

	//defer c.Stop() Aqui no se hace, eso sera en el CierraConexiones.


	n := Nodo{-1, c, dc, 0} //Creo el nodo y lo inicializo.

	lista.PushBack(n)

}

func (s *exportaServer) CierraConexiones(serverAddr string){
	//Busco nodo en la lista/map y hago un .Remove sobre el.

	//Saco el nodo de la lista
	//Saco el cliente.
	//Paro el cliente
	//Elimino el objeto Nodo de la lista de objetos Nodo

	//Fin
}



