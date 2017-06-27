package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var allowOrigin string

type FileMap struct {
	FileName string     `json:"fileName"`
	FilePath string     `json:"filePath"`
	IsDir    bool       `json:"isDir"`
	Children []*FileMap `json:"children"`
}

func getChildren(fullFilePath string, relativeFilePath string, fileMap *FileMap) {
	files, err := ioutil.ReadDir(fullFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := 0; i < len(files); i++ {
		childFullFilePath := fmt.Sprintf("%s%s%s", fullFilePath, string(os.PathSeparator), files[i].Name())
		childRelativeFilePath := fmt.Sprintf("%s%s%s", relativeFilePath, string(os.PathSeparator), files[i].Name())
		childFileMap := &FileMap{FileName: files[i].Name(), FilePath: childRelativeFilePath, IsDir: files[i].IsDir(), Children: nil}
		fileMap.Children = append(fileMap.Children, childFileMap)
		if files[i].IsDir() {
			getChildren(childFullFilePath, childRelativeFilePath, childFileMap)
		}
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		fileHandlerPut(w, r)
	default:
		fileHandlerGet(w, r)
	}
}

func fileHandlerPut(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	baseDir := getBaseDir(r, values)
	filePath := values.Get("fp")
	if len(filePath) > 0 && filePath[0:1] == string(os.PathSeparator) {
		filePath = filePath[1:]
	}
	filePath = fmt.Sprintf("%s%s%s", baseDir, string(os.PathSeparator), filePath)
	fileContents := values.Get("contents")
	// write file
	err = ioutil.WriteFile(filePath, []byte(fileContents), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// read file and return
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func fileHandlerGet(w http.ResponseWriter, r *http.Request) {
	baseDir := getBaseDir(r, r.URL.Query())
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
	baseDir := getBaseDir(r, r.URL.Query())
	log.Printf("BASE DIR IS %s", baseDir)
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

func getBaseDir(r *http.Request, values url.Values) string {
	baseDir := os.Getenv("MINIENV_DIR")
	srcDir := ""
	srcDirs := values["src"]
	if len(srcDirs) > 0 {
		srcDir = string(srcDirs[0])
	}
	if len(srcDir) > 1 {
		if srcDir[0:1] == string(os.PathSeparator) {
			baseDir += srcDir
		} else {
			baseDir += string(os.PathSeparator)
			baseDir += srcDir
		}
	}
	return baseDir
}

func addCorsAndCacheHeadersThenServe(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Add("Cache-Control", "no-store, must-revalidate")
		w.Header().Add("Expires", "0")
		h.ServeHTTP(w, r)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}
	if _, err := strconv.Atoi(os.Args[1]); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}
	allowOrigin = os.Getenv("MINIENV_ALLOW_ORIGIN")
	staticFileHandler := http.FileServer(http.Dir("public"))
	http.HandleFunc("/api/files", fileListHandler)
	http.HandleFunc("/api/file", fileHandler)
	http.Handle("/", addCorsAndCacheHeadersThenServe(staticFileHandler))
	err := http.ListenAndServe(":"+os.Args[1], nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
