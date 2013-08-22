package config
import(
	"io"
	"bufio"
	"os"
	"log"
	"idiom/util"
)

type serverInfo struct{
	RootContext string
	Port string
	IdiomFile string
	Token string
}

var ServerInfo=serverInfo{"/","8081","idiom.txt","shixiao"}

var Idiom=map[string] []string{}

func InitIdiom() {
	if util.FileExist(ServerInfo.IdiomFile){
		file,err:=os.Open(ServerInfo.IdiomFile)
		defer file.Close()
		if err==nil{
			br := bufio.NewReader(file)
			for{
				lineByte ,isPrefix, err := br.ReadLine()
				if isPrefix{}
				if err == io.EOF {
						break
				}else{
					head:=string(lineByte[0:3])
					if Idiom[head]==nil{
						Idiom[head]=[]string{}
					}
					Idiom[head]=append(Idiom[head],string(lineByte[0:len(lineByte)]))
				}
			}
		}
		log.Println("Init idiom successed")
	}else{
		log.Println("Init idiom failed, Can not found file ["+ServerInfo.IdiomFile+"]")
	}

}

var Players=map[string] string{}

func ProcArgs(args []string) {
	l:=len(args)
	if l>1{
		ServerInfo.Port=args[1]
	}else if l>2{
		ServerInfo.RootContext=args[2]
	}
}