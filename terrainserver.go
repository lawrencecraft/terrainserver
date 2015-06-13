package main

import (
	"github.com/lawrencecraft/terrainmodel/drawer"
	"github.com/lawrencecraft/terrainmodel/generator"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	g := generator.NewDiamondSquareGenerator(1.0, 1025, 1025)
	http.HandleFunc("/map", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		d := drawer.NewPngDrawer(w)
		t, _ := g.Generate()
		d.Draw(t)
	})

	err := http.ListenAndServe(":8822", nil)
	if err != nil {
		panic(err)
	}
}
