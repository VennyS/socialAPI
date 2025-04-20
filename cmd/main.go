package main

import (
	"net/http"
	"socialAPI/internal/setting"
)

func main() {
	app := setting.App{}
	app.LoadConfig()
	app.ConnectDB()
	app.MountServices()

	r := app.MountRouter()

	http.ListenAndServe(":8080", r)
}
