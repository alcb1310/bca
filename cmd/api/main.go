package main

func main() {
	db := database.New()
	server, _ := server.NewServer(db)

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("cannot start server", "error", err)
		panic("cannot start server")
	}
}
