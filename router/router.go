package router

import(
	"net/http"
	"fmt"
	"time"
	"encoding/xml"
	"idiom/util"
	"log"
	"idiom/config"
	"io/ioutil"
	"reflect"
	"math/rand"
)

type InputMessage struct{
	XMLName xml.Name `xml:"xml"`
	ToUserName string
	FromUserName string
	CreateTime int
	MsgType string
	Content string
	MsgId string
	PicUrl string
	Location_X string
	Location_Y string
	Scale string
	Label string
	Title string
	Description string
	Url string
	Event string
	EventKey string
}
func (m *InputMessage) ToXml() ([]byte,error){
	return xml.MarshalIndent(m,"","  ")
}
func (m *InputMessage) FromXml(b []byte){
	xml.Unmarshal(b,m);
}


type ArticleItem struct{
	XMLName xml.Name `xml:"item"`
	Title string
	Description string
	PicUrl string
	Url string
}
type OutputMessage struct{
	XMLName xml.Name `xml:"xml"`
	ToUserName string
	FromUserName string
	CreateTime int
	MsgType string
	Content string
	MusicUrl string
	HQMusicUrl string
	ArticleCount int
	Articles [] ArticleItem
}
func (m *OutputMessage) ToXml()([]byte,error){
	return xml.MarshalIndent(m,"","  ")
}

var routerMap=map[string] func(w http.ResponseWriter, r *http.Request,m *InputMessage){}
func AddRouter(key string,f func(w http.ResponseWriter, r *http.Request,m *InputMessage)){
	routerMap[key]=f
	log.Println("Add router ["+key+"] successed")

}

func Router(w http.ResponseWriter, r *http.Request){
	path:=r.URL.Path[len(config.ServerInfo.RootContext):]
	if(util.FileExist(path)){
		filecontent,err:=ioutil.ReadFile(path)
		if err==nil{
			fmt.Fprintf(w,"%s",filecontent)
		}else{
			log.Printf("%s",err);
		}

	}else if len(r.URL.RawQuery)>0{
		weibosignin(w,r);
	}else{
		content,err:=ioutil.ReadAll(r.Body)
		if err==nil{
			log.Printf("%s\n",content)
			m:=new(InputMessage)
			m.FromXml(content)
			f:=routerMap[m.MsgType]
			if f!=nil{
				f(w,r,m)
			}else{
				fmt.Fprintf(w,"%s","Can not find router[MsgType="+m.MsgType+"]")
			}
		}else{
			log.Printf("%s",err);
		}
	}
}

func weibosignin(w http.ResponseWriter, r *http.Request){
	queryData:=r.URL.Query()
	signature:=util.GetParameter(queryData,"signature")
	timestamp:=util.GetParameter(queryData,"timestamp")
	nonce:=util.GetParameter(queryData,"nonce")
	echostr:=util.GetParameter(queryData,"echostr")
	log.Printf("signature=%s timestamp=%s nonce=%s echostr=%s",signature,timestamp,nonce,echostr)
	fmt.Fprintf(w,"%s",echostr)
}

func Router_Text(w http.ResponseWriter, r *http.Request,m *InputMessage){
	ret:=""
	pWord:=config.Players[m.FromUserName]
	switch m.Content{
	case "?":
		ret=help(m)
	case "2":
		ret=newWord(m)
	default:
		if checkIdiom(m){
			if checkHead(m){
				bottom:=m.Content[len(m.Content)-3:len(m.Content)]
				pWord=ContinueIdiom(bottom)
				if len(pWord)==0{
					ret="你赢了"
					config.Players[m.FromUserName]=""
				}else{
					config.Players[m.FromUserName]=pWord
					ret="请接["+pWord+"]"
				}
			}else{
				ret="请接["+pWord+"]"
			}
		}else{
			ret="["+m.Content+"]不是成语"
		}
	}

	retXml,err:=GetOutputMessage(m,ret).ToXml()
	if err!=nil{
		retXml=[]byte{0x00}
	}
	fmt.Fprintf(w,"%s",retXml)
	
}

func Router_Image(w http.ResponseWriter, r *http.Request,m *InputMessage){
	fmt.Fprintf(w,"%s",m.MsgType)
}

func Router_Location(w http.ResponseWriter, r *http.Request,m *InputMessage){
	fmt.Fprintf(w,"%s",m.MsgType)
}

func Router_Link(w http.ResponseWriter, r *http.Request,m *InputMessage){
	fmt.Fprintf(w,"%s",m.MsgType)
}

func Router_Event(w http.ResponseWriter, r *http.Request,m *InputMessage){
	fmt.Fprintf(w,"%s",m.MsgType)
}

func GetOutputMessage(m *InputMessage,c string)*OutputMessage {
	ret:=new(OutputMessage)
	ret.ToUserName=m.FromUserName
	ret.FromUserName=m.ToUserName
	ret.CreateTime=time.Now().Nanosecond()
	ret.MsgType="text"
	ret.Content=c
	return ret
}

func RandomIdiom()string {
	v1:=reflect.ValueOf(config.Idiom)
	num1:=len(v1.MapKeys())
	ran1:=rand.Intn(num1)
	v2:=v1.MapIndex(v1.MapKeys()[ran1])
	num2:=v2.Len()
	ran2:=rand.Intn(num2)
	ret:=v2.Index(ran2).String()
	return ret
}

func ContinueIdiom(head string) string {
	ret:=""
	v:=config.Idiom[head]
	if v!=nil{
		num:=len(v)
		ran:=rand.Intn(num)
		ret=v[ran]
	}
	return ret
}

func help(m *InputMessage)string{
	ret:="发送‘2’开始新的成语接龙"
	pWord:=config.Players[m.FromUserName]
	if len(pWord)>0{
		ret+=",请接["+pWord+"]"
	}
	return ret
}

func newWord(m *InputMessage)string{
	pWord:=RandomIdiom()
	config.Players[m.FromUserName]=pWord
	ret:="请接["+pWord+"]"
	return ret
}

func checkIdiom(m *InputMessage) bool {
	found:=false
	head:=m.Content[0:3]
	if config.Idiom[head]!=nil{
		for i:=0;i<len(config.Idiom[head]);i++{
			if m.Content==config.Idiom[head][i]{
				found=true
				break
			}
		}
	}
	return found
}

func checkHead(m *InputMessage)bool {
	ret:=false
	pWord:=config.Players[m.FromUserName]
	if len(pWord)>3{
		bottom:=pWord[len(pWord)-3:len(pWord)]
		head:=m.Content[0:3]
		ret=head==bottom
	}else{
		ret=true
	}
	return ret
}