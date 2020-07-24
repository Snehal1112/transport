/*
Copyright © 2020 Snehal Dangroshiya <snehaldangroshiya@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/snehal1112/transport/client"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Command which used to insert the record",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		var err  error
		var collection, url, path string
		collection, err = flags.GetString("collection")
		if err != nil {
			err = fmt.Errorf("invalid collection name %w", err)
			log.Fatal(err)
		}

		url, err = flags.GetString("url")
		if err != nil {
			err = fmt.Errorf("connection url is required %w", err)
			log.Fatal(err)
		}

		path , err = flags.GetString("data_file")
		if err != nil {
			err = fmt.Errorf("connection url is required %w", err)
			log.Fatal(err)
		}

		var strData []byte
		strData, err = ioutil.ReadFile(path)
		type data map[string]interface{}
		var datas []data

		err = json.Unmarshal(strData, &datas)
		if err != nil {
			err = fmt.Errorf("issue while unmarshaling %w", err)
			log.Fatal(err)
		}

		conn := client.NewConnection(context.Background(), url)
		for _, items := range datas {
			err = conn.CreateDocument(collection, items)
			if err != nil {
				err = fmt.Errorf("issue while unmarshaling %w", err)
				log.Fatal(err)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	insertCmd.Flags().StringP("collection", "c", "", "Collection name")
	insertCmd.Flags().StringP("url", "u", "", "Connection url")
	insertCmd.Flags().StringToStringP("data", "d", map[string]string{}, "Collection data")
	insertCmd.Flags().StringP("data_file", "f", "", "Insert data file path")
}
