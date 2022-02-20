package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
    			<head>
        			<meta charset="UTF-8">
        			<title>チャット</title>
    			</head>
    			<body>
        			チャットをしよう！
    			</body>
			</html>
		`))
	})
	//Webサーバーを起動
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
