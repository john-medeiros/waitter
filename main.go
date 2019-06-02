package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/natefinch/lumberjack"
)

// GlobalConfig representa as configurações do daemon/serviço
type GlobalConfig struct {
	Logging struct {
		Path       string `json:"path"`
		MaxSize    int    `json:"max_size"`
		MaxBackups int    `json:"max_backups"`
		MaxAge     int    `json:"max_age"`
		Compress   bool   `json:"compress"`
		Enabled    bool   `json:"enabled"`
	} `json:"logging"`
	ConfigRepo struct {
		Path    string `json:"path"`
		Enabled bool   `json:"enabled"`
	} `json:"config_repo"`
	Watch struct {
		TimeSleep int  `json:"time_sleep"`
		Enabled   bool `json:"enabled"`
	} `json:"watch"`
}

// LoadConfiguration Lê o arquivo de configurações
func LoadConfiguration(file string) (GlobalConfig, error) {
	var configuration GlobalConfig
	configurationFile, err := os.Open(file)
	defer configurationFile.Close()
	if err != nil {
		//fmt.Println(err.Error())
		return configuration, err
	}
	jsonParser := json.NewDecoder(configurationFile)
	jsonParser.Decode(&configuration)
	return configuration, err
}

type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

type FileList struct {
	path          string
	name          string
	fullPath      string
	size          int64
	lastWriteTime time.Time
	md5           string
}

// Implementa um tipo usado para ordenar a lista de arquivos
type ByFileSize []FileList

func (a ByFileSize) Len() int           { return len(a) }
func (a ByFileSize) Less(i, j int) bool { return a[i].size < a[j].size }
func (a ByFileSize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// FileConfig é a representacao de uma configuração de processamento de arquivo
type FileConfig struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Watch       struct {
		Path      string `json:"path"`
		RegEx     string `json:"regex"`
		Recursive bool   `json:"recursive"`
	} `json:"watch"`
	Tasks []Task `json:"tasks"`
}

// Task é o tipo de tarefa que será executada para o arquivo.
type Task struct {
	Order      int                 `json:"order"`
	Type       string              `json:"type"`
	Parameters []map[string]string `json:"parameters"`
}

// TaskTypeFileRemove Representa uma tarefa de demoção de arquivo
type TaskTypeFileRemove struct {
	file string
}

type ByOrder []Task

func (a ByOrder) Len() int           { return len(a) }
func (a ByOrder) Less(i, j int) bool { return a[i].Order < a[j].Order }
func (a ByOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Remove arquivo
func (r *TaskTypeFileRemove) Remove() error {
	err := os.Remove(r.file)
	if err != nil {
		return err
	}
	return nil
}

// TaskTypeFileCopy Representa uma tarefa de copia
type TaskTypeFileCopy struct {
	source      string
	destination string
}

// Copy files
func (c *TaskTypeFileCopy) Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// LoadConfigFile Recebe o nome de um arquivo e retorna ele em forma de Struct
func LoadConfigFile(file string) (FileConfig, error) {
	var x FileConfig
	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("- ERRO em LoadConfigFile ao ler arquivo de configuração de tarefas. Detalhes: ", err)
		return x, err
	}
	data := FileConfig{}
	err = json.Unmarshal([]byte(f), &data)
	if err != nil {
		log.Println("- ERRO em LoadConfigFile ao converter JSON para Objeto. Detalhes: ", err)
		return data, err
	}
	return data, nil
}

// GetConfigFiles
func GetConfigFiles() ([]FileConfig, error) {
	files := []FileConfig{}
	arq1, err := LoadConfigFile("D:/file2.json")
	if err != nil {
		fmt.Println("GetConfigFiles")
		fmt.Println(err)
		return files, err
	}
	files = append(files, arq1)
	arq2, err := LoadConfigFile("D:/file3.json")
	if err != nil {
		fmt.Println("GetConfigFiles")
		fmt.Println(err)
		return files, err
	}
	files = append(files, arq2)
	return files, nil
}

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

var globalConfig, _ = LoadConfiguration("../../configs/waitter.json")

func init() {
	file, err := LoadConfiguration("../../configs/waitter.json")
	if err != nil {
		fmt.Println("Erro ao carregar configurações.")
		os.Exit(1)
	}

	globalConfig = file

	log.SetOutput(&lumberjack.Logger{
		Filename:   globalConfig.Logging.Path,
		MaxSize:    1, // megabytes
		MaxBackups: 6,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

}

// QueueExec armazena todas as execuções
type QueueExec struct {
	RunFileConfig FileConfig
	RunFileList   FileList
}

func main() {

	keepRunning := globalConfig.Watch.Enabled
	timeSleep := time.Duration(globalConfig.Watch.TimeSleep)
	fmt.Println(keepRunning)
	fmt.Println(globalConfig)

	for keepRunning {

		log.Println("- INFO - Final loop. Intervalo: ", globalConfig.Watch.TimeSleep)
		time.Sleep(timeSleep)
	}

}
