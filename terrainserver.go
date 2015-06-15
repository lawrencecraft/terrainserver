package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/jingweno/negroni-gorelic"
	"github.com/lawrencecraft/terrainmodel/drawer"
	"github.com/lawrencecraft/terrainmodel/generator"
	"net/http"
	"runtime"
)

func main() {
	port := flag.Int("p", 8822, "Port to listen for connections")
	nrkey := flag.String("nrkey", "", "NewRelic license key")
	flag.Parse()
	log.SetLevel(log.InfoLevel)
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Info("Setting max procs to", runtime.NumCPU())

	g := generator.NewDiamondSquareGenerator(1.0, 1025, 1025)

	r := http.NewServeMux()
	r.HandleFunc("/map", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		d := drawer.NewPngDrawer(w)
		t, _ := g.Generate()
		d.Draw(t)
	})

	n := negroni.New()

	if *nrkey != "" {
		log.Info("Using NewRelic")
		n.Use(negronigorelic.New(*nrkey, "terrainworker", true))
	}

	n.UseHandler(r)

	n.Run(fmt.Sprintf(":%d", *port))
}
