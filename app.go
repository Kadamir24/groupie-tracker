package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Data struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}

type Page struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	LocationURL  string   `json:"locations"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relation     string   `json:"relations"`
	RelationData Relations
	Location     LocationData
}
type Relations struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationData struct {
	Locations []string `json:"locations"`
	index     int
}

var page []Page

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	api, err := http.Get("https://groupietrackers.herokuapp.com/api")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		// panic(err)
	}
	defer api.Body.Close()
	tmp, err := ioutil.ReadAll(api.Body)
	if err != nil {
		panic(err)
	}
	var data Data
	json.Unmarshal(tmp, &data)

	log.Println(data.Artists)
	resp, err := http.Get(data.Artists)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	dataArt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dataArt, &page)
	if err != nil {
		panic(err)
	}

	ch := make(chan LocationData)
	for k, v := range page {
		go func(ch chan LocationData, v Page, k int) {
			jjj, err := http.Get(v.LocationURL)
			if err != nil {
				panic(err)
			}
			temp, err := ioutil.ReadAll(jjj.Body)
			var qqq LocationData
			json.Unmarshal(temp, &qqq)
			qqq.index = k
			ch <- qqq
		}(ch, v, k)
	}
	for range page {
		select {
		case qqq := <-ch:
			page[qqq.index].Location = qqq
		}
	}

	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		errorHandler(w, r, 500)
		return
		//panic(err)
	}

	e := t.Execute(w, page)
	if e != nil {
		fmt.Fprint(w, e)
	}
	// searchbox

}

func artHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("template/art.html")
	if err != nil {
		errorHandler(w, r, 500)
		// panic(err)
	}
	text := r.FormValue("username")

	for i, v := range page {
		if v.Name == text {
			jjj, err := http.Get(v.Relation)

			if err != nil {
				panic(err)
			}
			temp, err := ioutil.ReadAll(jjj.Body)
			var qqq Relations
			json.Unmarshal(temp, &qqq)

			page[i].RelationData = qqq
			t.Execute(w, page[i])
		}

	}

}

func oshibochka(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/oshibka.html")
	if err != nil {
		errorHandler(w, r, 500)
		// panic(err)
	}
	search := r.FormValue("artist")

	count := false
	for _, v := range page {

		// str := []string{v.Name}

		if strings.ToLower(v.Name) == strings.ToLower(search) {
			count = true
			res := Search{
				Input:  v.Name,
				Output: "Est takaya artist ji est",
			}

			t.Execute(w, res)

		} else if strconv.Itoa(v.CreationDate) == strings.ToLower(search) {
			count = true
			res := Search{
				Input:  v.Name,
				Output: search + " eto ihniy(egonniy) creation date ji est",
			}

			t.Execute(w, res)

		} else if strings.ToLower(v.FirstAlbum) == strings.ToLower(search) {
			count = true
			res := Search{
				Input:  v.Name,
				Output: search + " eto ihniy(egonniy) first album ji est",
			}

			t.Execute(w, res)

		}
		for _, k := range v.Members {
			if strings.ToLower(k) == strings.ToLower(search) {
				count = true
				res := Search{
					Input:  v.Name,
					Output: search + " eto ihniy(egonniy) member ji est",
				}

				t.Execute(w, res)
			}
		}
		for _, k := range v.Location.Locations {
			if strings.ToLower(k) == strings.ToLower(search) {
				count = true
				res := Search{
					Input:  v.Name,
					Output: search + " eto ihniy(egonniy) location concerta ji est",
				}

				t.Execute(w, res)
			}
		}
	}
	if !count {
		oshibka := Search{
			Input: "NET TAKOGO JI EST",
		}
		t.Execute(w, oshibka)
	}

}

type Search struct {
	Input  interface{}
	Output interface{}
}

func main() {
	fs := http.FileServer(http.Dir("template"))
	http.Handle("/template/", http.StripPrefix("/template/", fs))
	http.HandleFunc("/", handler)
	http.HandleFunc("/art", artHandler)
	http.HandleFunc("/oshibka", oshibochka)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	// err, e := ascii.AsciiCall(strconv.Itoa(status)+" "+http.StatusText(status), "standard")
	if status == 500 {
		fmt.Fprint(w, "500 Internal Server Error")
		// tmpl, _ := template.ParseFiles("template/error.html")
		// tmpl.Execute(w, status)
		return
	} else if status == 404 {
		// fmt.Fprint(w, "404 Not Found")
		tmpl, _ := template.ParseFiles("template/error.html")
		tmpl.Execute(w, status)
		//tmpl, _ := template.ParseFiles("template/error.html")
		//tmpl.Execute(w, status)

	} else if status == 400 {
		fmt.Fprint(w, "400 Bad Request")
		return
	}
}
