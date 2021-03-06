package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/lawrencecraft/terrainmodel/drawer"
	"github.com/lawrencecraft/terrainmodel/generator"
	"github.com/meatballhat/negroni-logrus"
	"net/http"
	"runtime"
	"strconv"
)

func main() {
	port := flag.Int("p", 8822, "Port to listen for connections")
	loglevel := flag.String("loglevel", "Info", "Log level - values are (Debug|Info|Warn|Fatal)")
	flag.Parse()

	switch *loglevel {
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.Fatal("Unknown level: ", *loglevel)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Info("Setting max procs to", runtime.NumCPU())

	r := mux.NewRouter()

	r.HandleFunc("/map/{x:[0-9]+}/{y:[0-9]+}/{b:8}", imageHandler)

	n := negroni.New()

	n.Use(negronilogrus.NewMiddleware())

	n.UseHandler(r)

	n.Run(fmt.Sprintf(":%d", *port))
}

func imageHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	x, errx := strconv.ParseInt(vars["x"], 10, 0)
	y, erry := strconv.ParseInt(vars["y"], 10, 0)
	b, erry := strconv.ParseInt(vars["b"], 10, 0)

	if errx != nil || erry != nil || x > 1025 || y > 1025 {
		w.Header().Set("Content-Type", "text/ascii")
		w.Write([]byte("argument error"))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	g := generator.NewDiamondSquareGenerator(1.0, int(x), int(y))
	switch b {
	}
	d := drawer.NewPngDrawer(w)
	t, _ := g.Generate()
	d.Draw(t)
}
