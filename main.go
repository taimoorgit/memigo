package main

import (
	"bufio"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting server...")

	// Start listening on TCP port
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
	defer listener.Close()
	log.Info().Msg("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Failed to accept connection")
			continue
		}
		log.Info().Str("remote_addr", conn.RemoteAddr().String()).Msg("New connection accepted")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Info().Str("remote_addr", conn.RemoteAddr().String()).Msg("Handling connection")

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Warn().Str("remote_addr", conn.RemoteAddr().String()).Err(err).Msg("Connection closed or error reading")
			return
		}
		log.Info().Str("remote_addr", conn.RemoteAddr().String()).Str("message", message).Msg("Received message")

		response := "Message received: " + message
		conn.Write([]byte(response))
	}
}
