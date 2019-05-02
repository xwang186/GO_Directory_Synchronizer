package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"log"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		// host   = "192.168.1.5"	
		port = "9090"
		// remote = host + ":" + port

		remote = ":" + port
	)

	fmt.Println("Server initializing...")

	lis, err := net.Listen("tcp", remote)
	defer lis.Close()

	if err != nil {
		fmt.Println("Error occurs listening to the port! ", remote)
		os.Exit(-1)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("An Error occurs on client: ", err.Error())
			// os.Exit(0)
			continue
		}

		go receiveFile(conn)
	}
}

func receiveFile(con net.Conn) {
	var (
		res          string
		tempFileName string                    
		data         = make([]byte, 1024*1024)
		by           []byte
		databuf      = bytes.NewBuffer(by) 
		fileNum      int                  
	)
	defer con.Close()

	fmt.Println("New Connection established ", con.RemoteAddr())
	j := 0 
	for {
		length, err := con.Read(data)
		if err != nil {
			da := databuf.Bytes()
			fmt.Printf("client %v has disconnnected %2d %d \n", con.RemoteAddr(), fileNum, len(da))
			return
		}
		if 0 == j {
			res = string(data[0:8])
			if "fileover" == res { 
				xienum := int(data[8])
				mergeFileName := string(data[9:length])
				go mainMergeFile(xienum, mergeFileName) 
				res = "Received " + mergeFileName
				con.Write([]byte(res))
				fmt.Println(mergeFileName, "Received")
				return
			}	else if "crtfoler" == res{
				folderpath := string(data[9:length])
				os.MkdirAll(folderpath, os.ModePerm)
			}   else if "updtDict" == res{
				cont:=LogAllFile("./vcv")
				//cont:=sendBackDicInfo()
				con.Write([]byte(cont))
				fmt.Println("Sent Server's folder information to client")
			}	else if "ask4file" == res{
				msgs := string(data[9:length])
				files:=strings.Split(msgs,"\n")
				for x:=range files{
					if(len(files[x])>0){
						sendFile(files[x])
					}
				}
				fmt.Print("Files",files)
			}	else{ 
				fileNum = int(data[0])
				tempFileName = string(data[1:length]) + strconv.Itoa(fileNum)
				fmt.Println("create temp file:", tempFileName)
				fout, err := os.Create(tempFileName)
				if err != nil {
					fmt.Println("Cannot create temp file", tempFileName)
					return
				}
				fout.Close()
			}
		} else {
			// databuf.Write(data[0:length])
			writeTempFileEnd(tempFileName, data[0:length])
		}

		res = strconv.Itoa(fileNum) + "received"
		con.Write([]byte(res))
		j++
	}

}
func sendBackDicInfo()string{
	file,err := os.Open("./fileInfo")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	fileinfo,err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := fileinfo.Size()
	buffer := make([]byte,fileSize)

	file.Read(buffer)
	return string(buffer)
}
func writeTempFileEnd(fileName string, data []byte) {
	tempFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		// panic(err)
		fmt.Println("Open file failed", err)
		return
	}
	defer tempFile.Close()
	tempFile.Write(data)
}

func mainMergeFile(connumber int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error create file", err)
		return
	}
	defer file.Close()

	for i := 0; i < connumber; i++ {
		mergeFile(filename+strconv.Itoa(i), file)
	}

	for i := 0; i < connumber; i++ {
		os.Remove(filename + strconv.Itoa(i))
	}

}

func mergeFile(rfilename string, wfile *os.File) {

	// fmt.Println(rfilename, wfilename)
	rfile, err := os.OpenFile(rfilename, os.O_RDWR, 0666)
	defer rfile.Close()
	if err != nil {
		fmt.Println("Error occura when merging files", rfilename)
		return
	}

	stat, err := rfile.Stat()
	if err != nil {
		panic(err)
	}

	num := stat.Size()

	buf := make([]byte, 1024*1024)
	for i := 0; int64(i) < num; {
		length, err := rfile.Read(buf)
		if err != nil {
			fmt.Println("File Cannot Read.")
		}
		i += length

		wfile.Write(buf[:length])
	}

}