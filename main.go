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

	data := map[string]interface{}{}
	
	data["config"] = config
	//fmt.Println(boards)
	retboards := []map[string]interface{}{}

	//fmt.Println(cards)
	
	for _, b := range boards {
		board := map[string]interface{}{}
		//	fmt.Println(b.Name)
		board["name"] = b.Name
		//	fmt.Println(b.Cards(cards))
		board["decks"] = []map[string]interface{}{}
		for _, d := range b.Decks {
			deck := map[string]interface{}{}
			deck["name"] = d.Name
			deck["cards"] = d.Get(b.Cards(cards))
			board["decks"] = append(board["decks"].([]map[string]interface{}), deck)
		}
		retboards = append(retboards, board)
	}

	data["boards"] = retboards

	
	err := ioutil.WriteFile("data.json", []byte(ToJSON(data)+"\n"), 0644)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("data.js", []byte("data = "+ToJSON(data)+"\n"), 0644)
	if err != nil {
		panic(err)
	}

}
