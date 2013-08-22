package util
import 	(
	"net/url"
	"os"
)
func GetParameter(v url.Values,k string) string{
	if v[k]==nil {
		return ""
	}else{
		return v[k][0]
	}
}

func FileExist(path string) bool {
    _, err := os.Stat(path)
    if err != nil &&  os.IsNotExist(err) {
        return false
    }
    return true
}
