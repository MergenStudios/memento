package scripts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"memento/structs"
	"memento/utils"
	"os"
	"strings"
)

func Add() {
	config, err := utils.LoadConfig()
	if utils.Handle(err) != nil { return }

	var id string
	var trueName string
	var extentionsString string
	var dated string
	var determineTime string

	fmt.Println("Check memento types add --help for further specification")

	input := bufio.NewReader(os.Stdin)

	fmt.Print("Enter ID (SCREENSHOT): ")
	id, err = input.ReadString('\n')
	if utils.Handle(err) != nil {
		return
	}
	id = utils.CleanString(id)

	fmt.Print("Enter display name (Screenshot): ")
	trueName, err = input.ReadString('\n')
	if utils.Handle(err) != nil { return }
	trueName = utils.CleanString(trueName)

	fmt.Print("Enter allowed extentions seperated by commas: ")
	extentionsString, err = input.ReadString('\n')
	if utils.Handle(err) != nil { return }
	extentionsString = utils.CleanString(extentionsString)

	extentions := strings.Split(extentionsString, ",")

	fmt.Print("Enter the type of time data represented by the type (point|range): ")
	dated, err = input.ReadString('\n')
	if utils.Handle(err) != nil { return }
	dated = utils.CleanString(dated)

	for !utils.InList(dated, []string{"point", "range"}) {
		fmt.Print("The time data must be point or range: ")
		dated, _ = input.ReadString('\n')
	}

	fmt.Print("Enter how the time should be determined (mtime|video|mc-log): ")
	determineTime, err = input.ReadString('\n')
	if utils.Handle(err) != nil { return }
	determineTime = utils.CleanString(determineTime)

	for !utils.InList(determineTime, []string{"mtime", "video", "mc-logs"}) {
		fmt.Print("The method to determine the time must be mtime, video or mc-logs: ")
		determineTime, _ = input.ReadString('\n')
	}

	config[id] = structs.TypeEnum{
		TrueName:      trueName,
		Extensions:    extentions,
		Dated:         dated,
		DetermineTime: determineTime,
	}

	jsonString, err := json.Marshal(config)
	if utils.Handle(err) != nil { return }
	os.WriteFile("./config/typesEnum.json", jsonString, os.ModePerm)
}

func Remove() {
	var id string

	input := bufio.NewReader(os.Stdin)
	config, err := utils.LoadConfig()
	if utils.Handle(err) != nil { return }

	fmt.Print("Enter the ID of the type you wish to delete: ")
	id, _ = input.ReadString('\n')
	id = utils.CleanString(id)

	for true {
		if _, ok := config[id]; !ok {
			fmt.Print("Enter a valid ID (use memento types list to list all types): ")
			id, _ = input.ReadString('\n')
			id = utils.CleanString(id)
		} else {
			break
		}
	}

	delete(config, id)
	jsonString, _ := json.Marshal(config)
	os.WriteFile("./config/typesEnum.json", jsonString, os.ModePerm)
}

func List() {
	config, err := utils.LoadConfig()
	if utils.Handle(err) != nil { return }

	for key, value := range config {
		fmt.Println(key, value.TrueName, value.Extensions, value.Dated, value.DetermineTime)
	}
}