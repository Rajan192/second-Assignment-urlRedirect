


package main

import (
	"fmt"
	"net/http"
    "log"
	"github.com/gophercises/urlshort"
    "gopkg.in/yaml.v3"
)

func MapHandler(pathTourl map[string]string, fallback http.Handler) http.HandlerFunc  {
    return func(w http.ResponseWriter,r *http.Request){
        path :=r.URL.Path;
        newUrl,ok:=pathTourl[path];
        if ok{
            http.Redirect(w,r,newUrl, http.StatusFound)
            return
        }
        
        fallback.ServeHTTP(w,r);
       // return;
    }
    

}

func YamlHandler(yamlByte []byte,fallback http.Handler)(http.HandlerFunc,error){
    pathUrls:= parseYaml(yamlByte)
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback),nil
}

type structPath struct{
    path string
    url string
}

func parseYaml(yamlData []byte) []structPath{
  var pathUrl []structPath
  err2:=yaml.Unmarshal(yamlData,&pathUrl)
  if(err2!=nil){
      log.Fatal("error in Unmashal ",err2);
  }
  return pathUrl;

}

func buildMap(pathtourl []structPath)map[string]string{
    mapOfpathToUrl:=make(map[string]string)
    for _,v:=range pathtourl{
        mapOfpathToUrl[v.path]=v.url
        fmt.Println(v.path," ",v.url)
    }
    return mapOfpathToUrl;
}







func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler,_:= urlshort.YAMLHandler([]byte(yaml), mapHandler)
	
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8000", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}