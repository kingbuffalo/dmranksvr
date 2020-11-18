package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	//"strings"
)

type RankRecord struct {
	Name  string `json:"name"`
	score int    `json:"score"`
}

func loadDataFromDB() {
	db, _ := sql.Open("mysql", "root:root@(127.0.0.1:3306)/dmrank") // 设置连接数据库的参数
	defer db.Close()                                                //关闭数据库
	err := db.Ping()                                                //连接数据库
	if err != nil {
		panic("数据库连接失败")
		return
	}
	rows, _ := db.Query("select * from rank") //获取所有数据
	var name string
	var score int64
	for rows.Next() { //循环显示所有的数据
		rows.Scan(&name, &score)
	}
}

func rank(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//bodyStr := string(body)
	var u RankRecord
	if err := json.Unmarshal(body, &u); err == nil {

	}
	fmt.Fprintf(w, "hello\n")
}

func main() {
	fmt.Println("服务端正在启动")
	loadDataFromDB()

	http.HandleFunc("/rank", rank)
	http.ListenAndServe(":8090", nil)
}
