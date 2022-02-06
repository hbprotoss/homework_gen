package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"
)

type ViewData map[string]interface{}

func main() {
	rand.Seed(time.Now().UnixNano())
	var choice, count, copies int
	var config *GradeConfig
	var workConfig *WorkConfig
	initEnv()
	configs := initYamlConfig()
	htmlTpl := initTemplate()

	for {
		fmt.Println("请选择年级:")
		for i, v := range configs {
			fmt.Printf("%d. %s\n", i+1, v.GradeDesc)
		}
		_, err := fmt.Scanf(input("%d"), &choice)
		if err != nil || choice < 1 || choice > len(configs) {
			println("请输入正确选项")
		} else {
			config = &configs[choice-1]
			break
		}
	}

	for {
		fmt.Println("请选择题目:")
		for i, v := range config.WorkConfigs {
			fmt.Printf("%d. %s\n", i+1, v.WorkDesc)
		}
		_, err := fmt.Scanf(input("%d"), &choice)
		if err != nil || choice < 1 || choice > len(config.WorkConfigs) {
			println("请输入正确选项")
		} else {
			workConfig = &config.WorkConfigs[choice-1]
			break
		}
	}

	for {
		fmt.Println("请输入要出几道题(如100):")
		_, err := fmt.Scanf(input("%d"), &count)
		if err != nil {
			println("请输入正确题目数量")
		} else {
			break
		}
	}

	for {
		fmt.Println("请输入份数(如2):")
		_, err := fmt.Scanf(input("%d"), &copies)
		if err != nil {
			println("请输入正确份数")
		} else {
			break
		}
	}

	var viewDatas []ViewData
	println()
	fmt.Printf("已选择: %s, %s, 共%d份, 每份%d题\n\n", config.GradeDesc, workConfig.WorkDesc, copies, count)
	for i := 0; i < copies; i++ {
		viewData := ViewData{
			"grade":  config.GradeDesc,
			"work":   workConfig.WorkDesc,
			"result": []WorkResult{},
			"index":  i + 1,
		}
		for j := 0; j < count; j++ {
			result := workConfig.Gen.Gen()
			viewData["result"] = append(viewData["result"].([]WorkResult), result)
			//fmt.Printf("%s = %d\n", result.Question, result.Answer)
		}
		fmt.Println("第", i+1, "份")

		viewDatas = append(viewDatas, viewData)
	}
	for _, v := range viewDatas {
		v["isQuestion"] = true
	}
	writeOutFile(htmlTpl, viewDatas, "口算题")
	for _, v := range viewDatas {
		v["isQuestion"] = false
	}
	writeOutFile(htmlTpl, viewDatas, "答案")

	var noUse string
	fmt.Scanf(input("%s"), &noUse)
}

func writeOutFile(htmlTpl *template.Template, viewDatas []ViewData, fileName string) {
	finalFileName := fmt.Sprintf("%s/%s.htm", OutDir, fileName)
	file, err := os.OpenFile(
		finalFileName,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		panic(err)
	}
	err = htmlTpl.Execute(file, viewDatas)
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s 已生成\n", finalFileName)
}

type GradeConfig struct {
	GradeDesc string

	WorkConfigs []WorkConfig
}

type WorkConfig struct {
	WorkDesc string

	Gen Generator
}

type YamlConfig struct {
	GradeDesc   string `yaml:"gradeDesc"`
	WorkConfigs []struct {
		WorkDesc       string   `yaml:"workDesc"`
		Min            int16    `yaml:"min"`
		Max            int16    `yaml:"max"`
		MaxResult      int16    `yaml:"maxResult"`
		Ops            []string `yaml:"ops"`
		OpCounts       []int8   `yaml:"opCounts"`
		UpgradeChecker string   `yaml:"upgradeChecker"`
		SpecialNumber  string   `yaml:"specialNumber"`
	} `yaml:"workConfigs"`
}

func initYamlConfig() []GradeConfig {
	content, err := ioutil.ReadFile("work.yml")
	if err != nil {
		panic(err)
	}
	var yamlConfigs []YamlConfig
	err = yaml.Unmarshal(content, &yamlConfigs)
	if err != nil {
		panic(err)
	}
	var configs []GradeConfig
	for _, v := range yamlConfigs {
		var workConfigs []WorkConfig
		for _, w := range v.WorkConfigs {
			workConfigs = append(workConfigs, WorkConfig{
				WorkDesc: w.WorkDesc,
				Gen:      NewWork(w.Min, w.Max, w.MaxResult, w.Ops, w.OpCounts, w.UpgradeChecker, w.SpecialNumber),
			})
		}
		configs = append(configs, GradeConfig{
			GradeDesc:   v.GradeDesc,
			WorkConfigs: workConfigs,
		})
	}
	return configs
}

func initEnv() {
	if _, err := os.Stat(OutDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(OutDir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func input(src string) string {
	if runtime.GOOS == "windows" {
		return src + "\n"
	} else {
		return src
	}
}
