package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type FileMap struct {
	FileName string
	IsDir bool
	Children []*FileMap
}

func getChildren(filePath string, fileMap *FileMap) {
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < len(files); i++ {
		childFileMap := &FileMap{FileName: files[i].Name(), IsDir: files[i].IsDir(), Children:nil}
		fileMap.Children = append(fileMap.Children, childFileMap)
		if files[i].IsDir() {
			getChildren(filePath + string(os.PathSeparator) + files[i].Name(), childFileMap)
		}
	}
}

//func printFileMap(fileMap *FileMap, depth int) {
//	for i := 0; i<len(fileMap.Children); i++ {
//		for j := 0; j < depth; j++ {
//			fmt.Print(" ")
//		}
//		fmt.Println(fileMap.Children[i].FileName)
//		if fileMap.Children[i].IsDir {
//			printFileMap(fileMap.Children[i], depth+1)
//		}
//	}
//}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	baseDir := os.Getenv("DC_DIR")
	filePath := r.URL.Query()["fp"][0]
	if len(filePath) > 0 && filePath[0:1] == "/" {
		filePath = filePath[1:]
	}
	filePath = baseDir + "/" + filePath
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func fileListHandler(w http.ResponseWriter, r *http.Request) {
	baseDir := os.Getenv("DC_DIR")
	fileMap := &FileMap{FileName: baseDir, IsDir: true, Children: nil}
	getChildren(baseDir, fileMap)
	//printFileMap(fileMap, 0)
	_, err := json.Marshal(fileMap)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileMap)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}
	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}
	staticFileHandler := http.FileServer(http.Dir("public"))
	http.HandleFunc("/api/files", fileListHandler)
	http.HandleFunc("/api/file", fileHandler)
	http.Handle("/", staticFileHandler)
	err := http.ListenAndServe(":"+os.Args[1], nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
