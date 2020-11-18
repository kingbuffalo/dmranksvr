package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	//"strings"
)

type ListReq struct {
	Page int `json:"page"`
	Num  int `json:"num"`
}

type RankRecord struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type ByScore []RankRecord

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score < a[i].Score }

var m map[string]int
var mMutex sync.Mutex
var arr []RankRecord
var aMutex sync.Mutex
var arrDiry bool

func loadDataFromDB() {
	m = make(map[string]int)
	db, _ := sql.Open("mysql", "root:root@(127.0.0.1:3306)/dmrank") // 设置连接数据库的参数
	defer db.Close()                                                //关闭数据库
	err := db.Ping()                                                //连接数据库
	if err != nil {
		panic("数据库连接失败")
		return
	}
	rows, _ := db.Query("select * from rank") //获取所有数据
	var name string
	var score int
	for rows.Next() { //循环显示所有的数据
		rows.Scan(&name, &score)
		m[name] = score
	}
	arrDiry = true
}

func ranklist(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var lr ListReq
	if arrDiry {
		arrDiry = false
		aMutex.Lock()
		mMutex.Lock()
		arr = make([]RankRecord, len(m))
		var idx int = 0
		for k, v := range m {
			arr[idx].Name = k
			arr[idx].Score = v
			idx++
		}
		mMutex.Unlock()
		sort.Sort(ByScore(arr))
		aMutex.Unlock()
	}
	if err := json.Unmarshal(body, &lr); err == nil {
		aMutex.Lock()
		defer aMutex.Unlock()
		var ret = arr[lr.Page*lr.Num : lr.Page*(lr.Num+1)]
		if str, err := json.Marshal(ret); err == nil {
			fmt.Fprintf(w, string(str))
		} else {
			fmt.Fprintf(w, err.Error())
		}
	} else {
		fmt.Fprintf(w, "[]")
	}
}

func rank(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//bodyStr := string(body)
	var u RankRecord
	if err := json.Unmarshal(body, &u); err == nil {
		mMutex.Lock()
		defer mMutex.Unlock()
		m[u.Name] = u.Score
		fmt.Fprintf(w, "ok\n")
		arrDiry = true
	} else {
		fmt.Fprintf(w, err.Error()+"\n")
	}
}

func main() {
	fmt.Println("服务端正在启动...")
	loadDataFromDB()

	http.HandleFunc("/rank", rank)
	http.HandleFunc("/ranklist", ranklist)
	http.ListenAndServe(":8090", nil)
}
