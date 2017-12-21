package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	_ "mysql"
	"net/http"
	"time"
	//"github.com/google/jsonapi"
	"strconv"
)

var db *sql.DB

var reqCount uint64

type Category struct {
	Id           int64     `json:"id"`
	Category     string    `json:"category"`
	Status       int       `json:"status"`
	Created_time time.Time `json:"created_time"`
}

type SuccessResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Profile struct {
	Name    string
	Hobbies []string
}

type Question struct {
	Id           uint64    `json:"id"`
	Question     string    `json:"question"`
	Right_ask    string    `json:"right_ask"`
	Ask1         string    `json:"ask1"`
	Ask2         string    `json:"ask2"`
	Ask3         string    `json:"ask3"`
	Category     string    `json:"category"`
	Created_time time.Time `json:"created_time"`
	Created_by   string    `json:"created_by"`
	System       string    `json:"system"`
	Updated_by   string    `json:"updated_by"`
	Priority     int       `json:"priority"`
}

//define your personal json format output
func (t *Question) MarshalJSON() ([]byte, error) {
	type Alias Question
	return json.Marshal(&struct {
		Created_time string `json:"created_time"`
		*Alias
	}{
		Created_time: t.Created_time.Format("2006-01-02 15:04:05"),
		Alias:        (*Alias)(t),
	})

}

func (t *Question) UnmarshalJSON(data []byte) error {
	type Alias Question
	aux := &struct {
		Created_time string `json:"created_time"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	form := "2006-01-02 15:04:05"
	time, err := time.Parse(form, aux.Created_time)
	if err != nil {
		log.Fatal("Parsing time err ", err)
		return err
	}
	t.Created_time = time
	return nil
}

func init() {
	var err error
	db, err = sql.Open("mysql", "swisse:swisse@tcp(10.50.115.114:16052)/swisse?charset=utf8&&readTimeout=3s&timeout=3s&loc=Local&autocommit=true&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func QueryQuestionById(w http.ResponseWriter, r *http.Request) {

	id := 23
	// query
	rows, err := db.Query("SELECT * FROM bw_askQuestion WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	q := &Question{}
	for rows.Next() {
		err := rows.Scan(&q.Id, &q.Question, &q.Right_ask, &q.Ask1, &q.Ask2, &q.Ask3, &q.Category, &q.Created_time, &q.Created_by, &q.System, &q.Updated_by, &q.Priority)
		if err != nil {
			fmt.Println(err)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("get question by ", q)
	if err != nil {
		log.Fatal("js parsing error", err)
	}

	js, err := json.Marshal(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("X-NSQ-Content-Type", "nsq; version=1.0")
	w.WriteHeader(200)
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	rows, err := db.Query("SELECT * FROM bw_category")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	categoryData := make([]*Category, 0)
	for rows.Next() {
		category := new(Category)
		err := rows.Scan(&category.Id, &category.Category, &category.Status, &category.Created_time)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		categoryData = append(categoryData, category)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	reqCount++
	log.Println("request count : ", reqCount)
	//for _, c := range categories {
	//	fmt.Sprintf("%d, %s, %s, £%.2f\n", c .id, c.category, c.status, c.created_time)
	//}
	js, err := json.Marshal(categoryData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("X-NSQ-Content-Type", "nsq; version=1.0")
	w.WriteHeader(200)
	w.Write(js)
}

func GoJson(w http.ResponseWriter, r *http.Request) {
	//question := json.Unmarshal(r.Body)

	jsonString := `{"created_time":"2017-03-17 14:10:27","id":23,"question":"公司福利奶粉规定，员工出勤满（）年，才可以享受公司福利奶粉","right_ask":"1年","ask1":"2年","ask2":"3年","ask3":"半年","category":"合生元","created_by":"1234","system":"合生元妈妈100微信","updated_by":"777771","priority":1}`
	question := &Question{}
	err := json.Unmarshal([]byte(jsonString), question)
	if err != nil {
		log.Fatal("parsing json string error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(" question : ", question)
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func InsertCategory(w http.ResponseWriter, r *http.Request) {
	c := Category{}
	json.NewDecoder(r.Body).Decode(&c)

	stmt, err := db.Prepare(`INSERT bw_category (category,status, created_time) values (?,?,?)`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(c.Category, c.Status, c.Created_time)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(id)
	js, err := json.Marshal(SuccessResp{Code: 200, Message: "success"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	c := Category{}
	json.NewDecoder(r.Body).Decode(&c)

	stmt, err := db.Prepare("UPDATE bw_category SET category=?  , status = ? , created_time = ? where id =?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(c.Category, c.Status, c.Created_time, c.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(SuccessResp{Code: 200, Message: "success"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	sId := r.URL.Query().Get("id")
	if len(sId) == 0 {
		http.Error(w, "req Id is empty", http.StatusInternalServerError)
		return
	}
	id, err := strconv.ParseInt(sId, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stmt, err := db.Prepare("delete from bw_category where id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	affected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(affected)
	js, err := json.Marshal(SuccessResp{Code: 200, Message: "row(s) " + fmt.Sprintf("%d", affected) + " affected"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	//Insert()
	http.HandleFunc("/", Index)
	http.HandleFunc("/json", GoJson)
	http.HandleFunc("/question", QueryQuestionById)
	http.HandleFunc("/insertCategory", InsertCategory)
	http.HandleFunc("/updateCategory", UpdateCategory)
	http.HandleFunc("/deleteCategory", DeleteCategory)
	log.Println("start server 9990")
	http.ListenAndServe(":9990", nil)
}
