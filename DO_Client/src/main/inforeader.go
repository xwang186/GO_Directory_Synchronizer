package main

import (
	"os"
	"fmt"
)
func showFileInfo(filepath string){
	fileInfo, err := os.Stat(filepath)
	if(os.IsNotExist(err)) {
	    fmt.Println("The target file ",filepath," does not exits")
	}
    fmt.Println("File name:", fileInfo.Name())
    fmt.Println("Size in bytes:", fileInfo.Size())
    fmt.Println("Permissions:", fileInfo.Mode())
    fmt.Println("Last modified:", fileInfo.ModTime())
    fmt.Println("Is Directory: ", fileInfo.IsDir())
    fmt.Printf("System interface type: %T\n", fileInfo.Sys())
}