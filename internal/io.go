package waitter

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"time"
)



// GetFileCheckSumMD5 retorna o checksum md5 de um arquivo
func GetFileCheckSumMD5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("- ERRO em GetFileCheckSumMD5 ao acessar arquivo para calcular MD5. Detalhes: ", err)
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

// DirectoryWatchList Lista diretórios com uma regex e retorna uma lista de arquivos
func DirectoryWatchList(directoryPath string, regex string) []FileList {
	var fileListReturn []FileList
	d, err := os.Open(directoryPath)
	if err != nil {
		log.Println("- ERRO em DirectoryWatchList ao testar existência do diretório informado. Detalhes: ", err)
		os.Exit(1)
	}
	defer d.Close()
	fileInfo, err := d.Readdir(-1)
	if err != nil {
		log.Println("- ERRO em DirectoryWatchList ao fazer a leitura do diretório informado. Detalhes: ", err)
		os.Exit(1)
	}
	for _, fileInfo := range fileInfo {
		if fileInfo.Mode().IsRegular() {
			if !fileInfo.IsDir() {
				matched, err := regexp.MatchString(regex, fileInfo.Name())
				if err != nil {
					log.Println("- ERRO em DirectoryWatchList ao validar RegEx com nome de arquivo. Detalhes: ", err)
					continue
				}
				if matched {
					ftl := FileList{}
					ftl.fullPath = path.Join(path.Dir(directoryPath), fileInfo.Name())
					ftl.path = directoryPath
					ftl.name = fileInfo.Name()
					ftl.size = fileInfo.Size()
					ftl.lastWriteTime = fileInfo.ModTime()
					md5checksum, err := GetFileCheckSumMD5(ftl.fullPath)
					if err != nil {
						log.Println("- ERRO em DirectoryWatchList ao obter o checksum de um arquivo. Detalhes: ", err)
						fmt.Println(err)
						continue
					}
					ftl.md5 = md5checksum
					fileListReturn = append(fileListReturn, ftl)
				}
			}
		}
	}
	return fileListReturn
}

// Interface used to sort structs.
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

// FileList represents a file in filesystem.
type FileList struct {
	path          string    // Directory
	name          string    // FileName
	fullPath      string    // FullPath = Directory + FileName
	size          int64     // FileSize
	lastWriteTime time.Time // Last Write Time
	md5           string    // MD5 checksum
}

// ByFileSize Represents a list of FileList and is used to sort by file Size
type ByFileSize []FileList

func (a ByFileSize) Len() int           { return len(a) }
func (a ByFileSize) Less(i, j int) bool { return a[i].size < a[j].size }
func (a ByFileSize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
