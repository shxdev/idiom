package odb

import(
//	"reflect"
	"encoding/json"
	"fmt"
)

type U struct{
	Name string
	Age int
}
func Do(){
	u:=U{"shixiao",36}
	Save(u)
}

func Save(v interface{}) {
/*	
	vt:=reflect.TypeOf(v)
	vv:=reflect.ValueOf(v)
	num:=vt.NumField()
	for i:=0;i<num;i++{
		k:=vt.Field(i).Type.Kind()
		n:=vt.Field(i).Name
		nv:=vv.FieldByName(n)
		switch k {
			case reflect.Int:
				fmt.Printf("%s=%d\n",n,nv.Int())
			default:
				fmt.Printf("%s=%s\n",n,nv.String())
		}
	}
*/
	b,err:=json.Marshal(v)
	if(err==nil){
		fmt.Printf("%s\n",b)
	}else{

	}
}

func Find() {
	
}

func Delete() {
	
}