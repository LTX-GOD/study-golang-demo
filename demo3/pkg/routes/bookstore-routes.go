package routes

import (
	"zsm/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	// 将 /book/ 路径和 POST 方法映射到 controllers.CreateBook 函数。
	// 也就是说，当一个 POST 请求发送到 /book/ 时，controllers.CreateBook 函数将被调用来处理这个请求。
	router.HandleFunc("/book/", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/book/", controllers.GetBook).Methods("GET")
	router.HandleFunc("/book/{BookId}", controllers.GetBookById).Methods("GET")
	router.HandleFunc("/book/{BookId}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/book/{BookId}", controllers.DeleteBook).Methods("DELETE")
}
