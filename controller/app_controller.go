package controller

import (
	"net/http"
	"strconv"
)

type AppController struct{}

func SetNextPage(r *http.Request) (initPage, nextPage int) {
	// ページネーションを行うためのデータを取得
	queryVal := r.URL.Query()
	page := queryVal.Get("page")
	var pageNum int
	if page != "" {
		var err error
		pageNum, err = strconv.Atoi(page)
		if err != nil {
			panic(err)
		}
	}

	// 初期ページのデータを下にoffsetの数値を取得
	initPage = 20
	nextPage = initPage * pageNum

	return
}

func getUserID(r *http.Request) (userID int, ok bool) {
	param := r.Context().Value("userID")
	id, ok := param.(float64)

	userID = int(id)

	return
}
