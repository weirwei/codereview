package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/weirwei/codereview/utils"
)

func TestLoadConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigName(".codereview-example")
	v.SetConfigType("yml")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	data, _ := v.Get(KNOWLEDGE).(map[string]any)
	var myStruct Knowledge
	err := mapstructure.Decode(data, &myStruct)
	if err != nil {
		fmt.Println("mapstructure decode error:", err)
		return
	}
	t.Log(utils.ToJson(myStruct))
}

func TestLoadUserConfig(t *testing.T) {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("ini")
	viper.SetConfigName(".codereview")

	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	data, _ := viper.Get("llm").(map[string]any)
	var myStruct LLMConf
	config := &mapstructure.DecoderConfig{
		Result:  &myStruct,
		TagName: "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		t.Fatal(err)
	}

	err = decoder.Decode(data)
	if err != nil {
		fmt.Println("mapstructure decode error:", err)
		return
	}
	t.Log(utils.ToJson(myStruct))
}
