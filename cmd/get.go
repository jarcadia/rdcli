/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"strings"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get PATH",
	Short: "get all (or some) dao values",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get called")
	},
	//ValidArgs: []string{"what", "the"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completePath(toComplete)
	},
}

func completePath(toComplete string) ([]string, cobra.ShellCompDirective) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	idx := strings.Index(toComplete, "/")
	if idx == -1 {
		return scanTypes(rdb, toComplete), cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	} else {
		objType := toComplete[0:idx]
		toComplete = toComplete[idx + 1:]
		return scanIds(rdb, objType, toComplete), cobra.ShellCompDirectiveNoFileComp
	}
}

func scanTypes(rdb *redis.Client, toComplete string) []string {
	var cursor uint64
	var keys []string
	var scanKeys []string
	var err error
	for {
		scanKeys, cursor, err = rdb.SScan(context.Background(), "rd/types", cursor, toComplete + "*", 10).Result()
		if err != nil {
			break
		}
		for _, key := range scanKeys {
			keys = append(keys, key + "/")
		}

		if cursor == 0 {
			break
		}
	}
	return keys
}

func scanIds(rdb *redis.Client, objType string, toComplete string) []string {
	var cursor uint64
	var keys []string
	var scanKeys []string
	var err error
	for {
		scanKeys, cursor, err = rdb.ZScan(context.Background(), objType, cursor, toComplete + "*", 10).Result()
		if err != nil {
			break
		}
		for i, key := range scanKeys {
			if i % 2 == 0 {
				keys = append(keys, objType + "/" + key)
			}
		}
		if cursor == 0 {
			break
		}
	}
	return keys
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
