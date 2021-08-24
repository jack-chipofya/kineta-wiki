package main

import (
    "log"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net/http"
    "encoding/json"
    "strconv"
    "github.com/gorilla/mux"
    )

type Article struct {
    Id int `json:"Id"`
    Title string `json:"Title"`
    Body string `json:"Body"`
    Author string `json:"Author"`
    Created_on string `json:"Created_on"`
    Last_updated string `json:"Last_updated"`
}

var Articles []Article
func initDb() {
    db,err:=sql.Open("mysql","root:sarley@#7@tcp(127.0.0.1)/wiki")
    if err!=nil {
        panic(err.Error())
    }
    defer db.Close()
}

func returnSingleArticle(w http.ResponseWriter,r* http.Request) {
    vars:=mux.Vars(r)
    key:=vars["id"]
    db,err:=sql.Open("mysql","root:sarley@#7@tcp(127.0.0.1)/wiki")
    if err!=nil {
        panic(err.Error())
    }
    defer db.Close()
    var article Article
    err=db.QueryRow("SELECT * FROM articles where id="+key).Scan(&article.Id,&article.Title,&article.Body,&article.Author,&article.Created_on,&article.Last_updated)
    if err!=nil {
        panic(err.Error())
    }
    json.NewEncoder(w).Encode(article)
}
func returnAllArticles(w http.ResponseWriter,r *http.Request) {
    prev_article:=len(Articles)
    db,err:=sql.Open("mysql","root:sarley@#7@tcp(127.0.0.1)/wiki")
    if err!=nil {
        panic(err.Error())
    }
    defer db.Close()
    
    results,err:=db.Query("SELECT * FROM articles where id>"+strconv.Itoa(prev_article))
    if err!=nil {
        panic(err.Error())
    }
    for results.Next() {
        var article Article
        err:=results.Scan(&article.Id,&article.Title,&article.Body,&article.Author,&article.Created_on,&article.Last_updated)
        if err!=nil {
            panic(err.Error())
        }
        Articles=append(Articles,article)
    }
    json.NewEncoder(w).Encode(Articles)
}

func handleRequest() {
    myRouter:=mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/articles/",returnAllArticles)
    myRouter.HandleFunc("/articles/{id}",returnSingleArticle)
    log.Fatal(http.ListenAndServe(":8080",myRouter))
}
func main() {
    handleRequest()
}