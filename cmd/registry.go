package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path/filepath"
)

var username string
var password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login",
	Long:  "Login to registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(username) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}
		if len(password) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}
		user, err := user.Current()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfgPath := filepath.Join(user.HomeDir, ".docker", "config.json")
		fmt.Printf("WARNING! Your password will be stored unencrypted in %s. \n", cfgPath)

		var cfgMap = make(map[string]interface{})

		// read origin config
		_, err = os.Stat(cfgPath)
		if err == nil {
			content, err := os.ReadFile(cfgPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = json.Unmarshal(content, &cfgMap)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// encode registry auth
		bytes := []byte(username + ":" + password)
		auth := base64.StdEncoding.EncodeToString(bytes)

		// write to config
		auths := cfgMap["auths"].(map[string]interface{})
		auths[args[0]] = auth

		newCfg, err := json.Marshal(cfgMap)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file, err := os.OpenFile(cfgPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		write := bufio.NewWriter(file)
		_, err = write.Write(newCfg)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = write.Flush()
		if err != nil {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	dockerCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
}
