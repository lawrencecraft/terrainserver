package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/lawrencecraft/terrainmodel/drawer"
	"github.com/lawrencecraft/terrainmodel/generator"
	"github.com/meatballhat/negroni-logrus"
	"net/http"
	"runtime"
)

func main() {
	port := flag.Int("p", 8822, "Port to listen for connections")
	loglevel := flag.String("loglevel", "Info", "Log level - values are (Debug|Info|Warn|Fatal)")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Info("Setting max procs to", runtime.NumCPU())

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

	g := generator.NewDiamondSquareGenerator(1.0, 1025, 1025)

	r := http.NewServeMux()
	r.HandleFunc("/map", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		d := drawer.NewPngDrawer(w)
		t, _ := g.Generate()
		d.Draw(t)
	})

	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())

	n.UseHandler(r)

	n.Run(fmt.Sprintf(":%d", *port))
}
