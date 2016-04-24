package main

import (

	"github.com/valyala/gorpc"
	"log"
	"container/list"
	"image"
	 //"io"	
    
     "golang.org/x/image/bmp"
    //"strings"
    "os"
     _ "image/gif"     
     _ "image/jpeg"
    

)



func main(){
	

	writer, err := os.Create("final.png")
    if err != nil{
        log.Fatal(err)
    }
    defer writer.Close()

    var img *image.YCbCr
   	gorpc.RegisterType(img)
    

	log.Print("Servidor RPC iniciado en direccion: 127.0.0.1:12345")

  		//Creo la lista enlazada donde guardo las direcciones
	l := list.New() //Deberia crearla fuera del main.
	//Registro el tipo IMG para trabajar con el.
	


	s := &gorpc.Server{
    		// Accept clients on this TCP address.
   	 		Addr: "0.0.0.0:12345",
   	 		
    		// Manejador de las peticiones del servidor
   	 		Handler: func(clientAddr string, request interface{}) interface{}{
   	 		imagen := request.(image.Image)
      	  	log.Printf("Obtained request from the client %s\n", clientAddr)

       	 	
			//log.Print("Respuesta '" , request ,"' enviada al cliente ", clientAddr)

			l.PushBack(clientAddr)
			log.Printf("IP guardada correctamente en la lista.")
			log.Print("Lo que me envia el cliente es....")
			
			//log.Print(request)
				
			//La imagen me llega dentro del request

			bmp.Encode(writer,imagen)
			     return request
   	 },



	}
	if err := s.Serve(); err != nil {
    	log.Fatalf("Cannot start rpc server: %s", err)
	}

	
}