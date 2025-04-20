package main

import (
	"flag"
	"net/http"
	"socialAPI/internal/setting"
)

func main() {
	var migrations bool
	flag.BoolVar(&migrations, "migrate", false, "Run migrations")
	flag.Parse()

	app := setting.App{}
	app.LoadConfig()
	app.InitStorages(migrations)
	app.MountServices()

	r := app.MountRouter()

	http.ListenAndServe(":8080", r)
}
