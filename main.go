package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

var version string
var filename string

func main() {
	go refreshVersion()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "https://spark.adobe.com/page/DWWqEQxXZZg8j/")
	})

	e.GET("/api/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, version)
	})
	e.GET("/api/filename", func(c echo.Context) error {
		return c.JSON(http.StatusOK, filename)
	})

	e.Static("/files", "public")
	e.Logger.Fatal(e.Start(":1323"))
}

func refreshVersion() {
	// it's late and I can't think anymore
	for {
		files, err := ioutil.ReadDir("./public")
		if err != nil {
			fmt.Println("Error reading files", err)
			return
		}
		if len(files) < 1 {
			fmt.Println("No files present")
			os.Exit(1)
		}
		sort.Slice(files, func(i, j int) bool {
			iname := files[i].Name()
			jname := files[j].Name()
			iver, err := strconv.ParseFloat(iname[0:len(iname)-len(filepath.Ext(iname))], 64) // what the heck
			if err != nil {
				fmt.Println("Error parsing int", err)
				return false
			}
			jver, err := strconv.ParseFloat(jname[0:len(jname)-len(filepath.Ext(jname))], 64)
			if err != nil {
				fmt.Println("Error parsing int", err)
				return false
			}
			return iver > jver
		})
		filename = files[0].Name()
		version = filename[0 : len(filename)-len(filepath.Ext(filename))]

		time.Sleep(time.Second * 120)
	}
}
