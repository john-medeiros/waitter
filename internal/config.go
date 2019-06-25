package waitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/john-medeiros/waitter/internal/config"
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

// FileConfig represents an user configuration to processing a file
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

// LoadConfiguration Lê o arquivo de configurações
func LoadConfiguration(file string) (waitter.GlobalConfig, error) {
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
