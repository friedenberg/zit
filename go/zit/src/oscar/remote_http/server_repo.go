package remote_http

import "io"

func (server *Server) writeInventoryList(
	request Request,
) (response Response) {
	server.Repo.GetEnv().GetUI().Print("would write")
	io.Copy(server.Repo.GetEnv().GetUIFile(), request.Body)
	return
}
