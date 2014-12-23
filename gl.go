package main

import (
	"github.com/balzaczyy/go-logging"
	_ "github.com/balzaczyy/golucene/core/codec/lucene42"
	"github.com/balzaczyy/golucene/core/index"
	"github.com/balzaczyy/golucene/core/search"
	"github.com/balzaczyy/golucene/core/store"
	"os"
)

var log = logging.MustGetLogger("test")

func main() {
	logging.SetBackend(logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stdout, "", 0),
		logging.MustStringFormatter(
			"%{color}%{time:15:04:05.000000} %{shortfunc:16.16s} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
		),
	))

	log.Info("Oepening FSDirectory...")
	path := "core/search/testdata/win8/belfrysample"
	d, err := store.OpenFSDirectory(path)
	if err != nil {
		panic(err)
	}
	log.Info("Opening DirectoryReader...")
	r, err := index.OpenDirectoryReader(d)
	if err != nil {
		panic(err)
	}
	if r == nil {
		panic("DirectoryReader cannot be opened.")
	}
	if len(r.Leaves()) < 1 {
		panic("Should have one leaf.")
	}
	for _, ctx := range r.Leaves() {
		if ctx.Parent() != r.Context() {
			panic("leaves not point to parent!")
		}
	}
	log.Info("Initializing IndexSearcher...")
	ss := search.NewIndexSearcher(r)
	log.Info("Searching...")
	docs, err := ss.SearchTop(search.NewTermQuery(index.NewTerm("content", "bat")), 10)
	if err != nil {
		panic(err)
	}
	log.Info("Hits: %v", docs.TotalHits)
	doc, err := r.Document(docs.ScoreDocs[0].Doc)
	if err != nil {
		panic(err)
	}
	log.Info("Hit 1's title: %v", doc.Get("title"))
}
