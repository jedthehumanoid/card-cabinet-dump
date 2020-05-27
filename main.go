package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jedthehumanoid/card-cabinet"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Src    string            `toml:"src"`
	Colors map[string]string `toml:"colors"`
}

func loadConfig(file string) Config {
	var config Config
	err := loadToml(file, &config)
	if err != nil {
		fmt.Println("Couldn't load configuration file")
	}
	if config.Src == "" {
		config.Src = "."
	}
	config.Src = filepath.Clean(config.Src) + "/"
	return config
}

func loadToml(file string, i interface{}) error {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(d), i)
	return err
}

func ToJSON(in interface{}) string {
	b, err := json.MarshalIndent(in, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func main() {
	config := loadConfig("cabinet.toml")

	files := cardcabinet.FindFiles(config.Src)
	cards := cardcabinet.ReadCards(files)
	boards := cardcabinet.ReadBoards(files)

	data := struct {
		Config Config              `json:"config"`
		Cards  []cardcabinet.Card  `json:"cards"`
		Boards []cardcabinet.Board `json:"boards"`
	}{
		config,
		cards,
		boards,
	}

	err := ioutil.WriteFile("data.json", []byte(ToJSON(data)+"\n"), 0644)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("data.js", []byte("data = "+ToJSON(data)+"\n"), 0644)
	if err != nil {
		panic(err)
	}

}
