package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//---------------------------structures---------------------------//

type Monsters []struct {
	Name        string
	Type        string
	Species     string
	Elements    []interface{}
	Description string
	Url         string
	Icon        string
	Weaknesses  []struct {
		Element string
		Stars   int
	}
	Rewards []struct {
		Rarity      int
		Value       int
		Name        string
		Description string
	}
	Locations []struct {
		ZoneCount int
		Name      string
	}
}

type Items []struct {
	Rarity      int
	Value       int
	Name        string
	Description string
}

type Weapons []struct {
	Name   string
	Type   string
	Rarity int
	Attack struct {
		Display int
		Raw     int
	}
	Assets struct {
		Icon  string
		Image string
	}
}

type Armors []struct {
	Name    string
	Type    string
	Rarity  int
	Defense struct {
		Base      int
		Max       int
		Augmented int
	}
	Resistances struct {
		Fire    int
		Water   int
		Ice     int
		Thunder int
		Dragon  int
	}
	Assets struct {
		ImageMale   string
		ImageFemale string
	}
}

type Foods []struct {
	Name        string
	Description string
	Recovery    struct {
		Actions []string
		Items   []struct {
			Name string
		}
	}
	Protection struct {
		Items []struct {
			Name string
		}
		Skills []struct {
			Name        string
			Description string
		}
	}
}

type Skills []struct {
	Name        string
	Description string
	Ranks       []struct {
		Level       int
		Description string
		Skillname   string
	}
}

//---------------------------variables---------------------------//
var index int
var iresponse = Items{}
var wresponse = Weapons{}
var aresponse = Armors{}
var fresponse = Foods{}
var sresponse = Skills{}
var mresponse = Monsters{}
var tmplerror = template.Must(template.ParseFiles("../pages/Error/Error_page.html"))

var url string
var mot string

func main() {
	connection()
	home()
	monsters()
	items()
	weapons()
	armors()
	foods()
	skills()
	http.ListenAndServe(":80", nil)
}

//---------------------------conexion au JS et au dossier Images---------------------------//
func connection() {
	cssglo := http.FileServer(http.Dir("../css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssglo))

	img := http.FileServer(http.Dir("../pictures"))
	http.Handle("/pictures/", http.StripPrefix("/pictures/", img))

	js := http.FileServer(http.Dir("../js"))
	http.Handle("/js/", http.StripPrefix("/js/", js))
}

//---------------------------Remplacer/Ajouter des images---------------------------//
func unknow() {
	for i := 0; i < len(wresponse); i++ {
		if wresponse[i].Assets.Icon == "" {
			wresponse[i].Assets.Icon = "../pictures/logo/None.png"
		}
		if wresponse[i].Assets.Image == "" {
			wresponse[i].Assets.Image = "../pictures/logo/None.png"
		}
	}
	for i := 0; i < len(aresponse); i++ {
		if aresponse[i].Assets.ImageMale == "" {
			aresponse[i].Assets.ImageMale = "../pictures/logo/unknown.png"
		}
		if aresponse[i].Assets.ImageFemale == "" {
			aresponse[i].Assets.ImageFemale = "../pictures/logo/None.png"
		}
	}
}

//---------------------------Erreur 404---------------------------//
func error404(chemin string, w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != chemin {
		return true
	}
	return false
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		tmplerror.Execute(w, nil)
	}
}

//---------------------------Page Home---------------------------//
func home() {
	tmpl := template.Must(template.ParseFiles("../pages/Home/Home.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if error404("/", w, r) {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			tmpl.Execute(w, nil)
		}
	})
}

//---------------------------conexion API personnelle---------------------------//
func monsters() {
	url = "https://thomasquadro.github.io/Data/api_image_monster.json"
	httpClient := http.Client{
		Timeout: time.Second * 15, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &mresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//---------------------------routes---------------------------//
	tmpl := template.Must(template.ParseFiles("../pages/Monsters/Monsters_index.html"))
	http.HandleFunc("/search_monsters", func(w http.ResponseWriter, r *http.Request) {
		error404("/search_monsters", w, r)
		tmpl.Execute(w, mresponse)
	})
	http.HandleFunc("/monsters/", func(w http.ResponseWriter, r *http.Request) {
		mot = strings.ReplaceAll(r.URL.Path, "/monsters/", "")
		monster0()
		tmpl := template.Must(template.ParseFiles("../pages/Monsters/Monsters_page.html"))
		if error404("/monsters/"+mresponse[index].Name, w, r) {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			tmpl.Execute(w, mresponse[index])
		}
	})
}

//---------------------------prendre l'index du mot associé---------------------------//
func monster0() {
	var tab []string
	for _, character := range mresponse {
		tab = append(tab, character.Name)
	}
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(mot) {
			index = i
		}
	}
}

//---------------------------conexion API items---------------------------//
func items() {
	url = "https://mhw-db.com/items"
	httpClient := http.Client{
		Timeout: time.Second * 15, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &iresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//---------------------------routes---------------------------//
	tmpl := template.Must(template.ParseFiles("../pages/Items/Items_index.html"))
	http.HandleFunc("/search_items", func(w http.ResponseWriter, r *http.Request) {
		error404("/search_items", w, r)
		tmpl.Execute(w, iresponse)
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		mot = strings.ReplaceAll(r.URL.Path, "/items/", "")
		items0()
		tmpl := template.Must(template.ParseFiles("../pages/Items/Items_page.html"))
		if error404("/items/"+iresponse[index].Name, w, r) {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			tmpl.Execute(w, iresponse[index])
		}
	})
}

//---------------------------prendre l'index du mot associé---------------------------//
func items0() {
	var tab []string
	for _, character := range iresponse {
		tab = append(tab, character.Name)
	}

	for i := 0; i < len(tab); i++ {
		if tab[i] == string(mot) {
			index = i
		}
	}
}

//---------------------------conexion API weapons---------------------------//
func weapons() {
	url = "https://mhw-db.com/weapons"
	httpClient := http.Client{
		Timeout: time.Second * 15, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &wresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//---------------------------routes---------------------------//
	tmpl := template.Must(template.ParseFiles("../pages/Weapons/Weapons_index.html"))
	http.HandleFunc("/search_weapons", func(w http.ResponseWriter, r *http.Request) {
		unknow()
		error404("/search_weapons", w, r)
		tmpl.Execute(w, wresponse)
	})

	http.HandleFunc("/weapons/", func(w http.ResponseWriter, r *http.Request) {
		unknow()
		mot = strings.ReplaceAll(r.URL.Path, "/weapons/", "")
		weapons0()
		tmpl := template.Must(template.ParseFiles("../pages/Weapons/Weapons_page.html"))
		if error404("/weapons/"+wresponse[index].Name, w, r) {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			tmpl.Execute(w, wresponse[index])
		}
	})
}

//---------------------------prendre l'index du mot associé---------------------------//
func weapons0() {
	var tab []string
	for _, character := range wresponse {
		tab = append(tab, character.Name)
	}
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(mot) {
			index = i
		}
	}
}

//---------------------------conexion API armors---------------------------//
func armors() {
	url = "https://mhw-db.com/armor"
	httpClient := http.Client{
		Timeout: time.Second * 15, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &aresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//---------------------------routes---------------------------//
	tmpl := template.Must(template.ParseFiles("../pages/Armors/Armors_index.html"))
	http.HandleFunc("/search_armors", func(w http.ResponseWriter, r *http.Request) {
		unknow()
		error404("/search_armors", w, r)
		tmpl.Execute(w, aresponse)
	})

	http.HandleFunc("/armors/", func(w http.ResponseWriter, r *http.Request) {
		unknow()
		mot = strings.ReplaceAll(r.URL.Path, "/armors/", "")
		armors0()
		tmpl := template.Must(template.ParseFiles("../pages/Armors/Armors_page.html"))
		if error404("/armors/"+aresponse[index].Name, w, r) {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			tmpl.Execute(w, aresponse[index])
		}
	})
}

//---------------------------prendre l'index du mot associé---------------------------//
func armors0() {
	var tab []string
	for _, character := range aresponse {
		tab = append(tab, character.Name)
	}
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(mot) {
			index = i
		}
	}
}

//---------------------------conexion API foods---------------------------//
func foods() {
	url = "https://mhw-db.com/ailments"
	httpClient := http.Client{
		Timeout: time.Second * 15, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &fresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//---------------------------routes---------------------------//
	tmpl := template.Must(template.ParseFiles("../pages/Foods/Foods_index.html"))
	http.HandleFunc("/search_foods", func(w http.ResponseWriter, r *http.Request) {
		error404("/search_foods", w, r)
		tmpl.Execute(w, fresponse)
	})

	http.HandleFunc("/foods/", func(w http.ResponseWriter, r *http.Request) {
		mot = strings.ReplaceAll(r.URL.Path, "/foods/", "")
		foods0()
		tmpl := template.Must(template.ParseFiles("../pages/Foods/Foods_page.html"))
		if error404("/foods/"+fresponse[index].Name, w, r) {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			tmpl.Execute(w, fresponse[index])
		}
	})
}

//---------------------------prendre l'index du mot associé---------------------------//
func foods0() {
	var tab []string
	for _, character := range fresponse {
		tab = append(tab, character.Name)
	}
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(mot) {
			index = i
		}
	}
}

//---------------------------conexion API skills---------------------------//
func skills() {
	url = "https://mhw-db.com/skills"
	httpClient := http.Client{
		Timeout: time.Second * 15, // define timeout
	}

	//create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	//make api call
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	//parse response
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &sresponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//---------------------------routes---------------------------//
	tmpl := template.Must(template.ParseFiles("../pages/Skills/Skills_index.html"))
	http.HandleFunc("/search_skills", func(w http.ResponseWriter, r *http.Request) {
		error404("/search_skills", w, r)
		tmpl.Execute(w, sresponse)
	})

	http.HandleFunc("/skills/", func(w http.ResponseWriter, r *http.Request) {
		mot = strings.ReplaceAll(r.URL.Path, "/skills/", "")
		skills0()
		tmpl := template.Must(template.ParseFiles("../pages/Skills/Skills_page.html"))
		if error404("/skills/"+sresponse[index].Name, w, r) {
			errorHandler(w, r, http.StatusNotFound)
		} else {
			tmpl.Execute(w, sresponse[index])
		}
	})
}

//---------------------------prendre l'index du mot associé---------------------------//
func skills0() {
	var tab []string
	for _, character := range sresponse {
		tab = append(tab, character.Name)
	}
	for i := 0; i < len(tab); i++ {
		if tab[i] == string(mot) {
			index = i
		}
	}
}
