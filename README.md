# GoChar
Web application that classifies a character from a picture using CNN.

Users access a frontend page to upload an image to the master backend server. When it receives a request, it chooses a working node from a cluster and sends a request using a RPC (Remote Procedure Call). The nodes execute a python script that identifies the character on the picture and returns the result in a callback.
