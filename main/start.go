package main

import(
	"net/http"
	"log"
	"idiom/router"
	"idiom/config"
	"os"
)

func webHandle(w http.ResponseWriter, r *http.Request){
	router.Router(w,r)
}

func main() {
	config.ProcArgs(os.Args)

	router.AddRouter("text",router.Router_Text)
	router.AddRouter("image",router.Router_Image)
	router.AddRouter("location",router.Router_Location)
	router.AddRouter("link",router.Router_Link)
	router.AddRouter("event",router.Router_Event)

	config.InitIdiom()
	//log.Println(router.RandomIdiom())

	http.HandleFunc(config.ServerInfo.RootContext,webHandle)
	log.Println("Server will start at \""+config.ServerInfo.RootContext+"\" using port "+config.ServerInfo.Port);
	err:=http.ListenAndServe(":"+config.ServerInfo.Port,nil)
	if(err!=nil){
		log.Fatal("ListenAndServe: ", err)
	}
}