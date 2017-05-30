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
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	IsDir bool `json:"isDir"`
	Children []*FileMap `json:"children"`
}

func getChildren(fullFilePath string, relativeFilePath string, fileMap *FileMap) {
	files, err := ioutil.ReadDir(fullFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < len(files); i++ {
		childFullFilePath := fmt.Sprintf("%s%s%s",fullFilePath, string(os.PathSeparator), files[i].Name())
		childRelativeFilePath := fmt.Sprintf("%s%s%s",relativeFilePath, string(os.PathSeparator), files[i].Name())
		childFileMap := &FileMap{FileName: files[i].Name(), FilePath: childRelativeFilePath, IsDir: files[i].IsDir(), Children:nil}
		fileMap.Children = append(fileMap.Children, childFileMap)
		if files[i].IsDir() {
			getChildren(childFullFilePath, childRelativeFilePath, childFileMap)
		}
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	baseDir := os.Getenv("EXUP_DIR")
	filePath := r.URL.Query()["fp"][0]
	if len(filePath) > 0 && filePath[0:1] == string(os.PathSeparator) {
		filePath = filePath[1:]
	}
	filePath = fmt.Sprintf("%s%s%s", baseDir, string(os.PathSeparator), filePath)
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
	baseDir := os.Getenv("EXUP_DIR")
	relativeDir := string(os.PathSeparator)
	fileMap := &FileMap{FileName: relativeDir, FilePath: relativeDir, IsDir: true, Children: nil}
	getChildren(baseDir, relativeDir, fileMap)
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
