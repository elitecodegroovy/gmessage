// +build !windows
// All rights reserved.

package message

// Run starts the gMessage server. This wrapper function allows Windows to add a
// hook for running NATS as a service.
func Run(server *Server) error {
	server.Start()
	return nil
}

// isWindowsService indicates if gMessage is running as a Windows service.
func isWindowsService() bool {
	return false
}
