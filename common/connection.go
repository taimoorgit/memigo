package common

import (
	"bufio"
	"net"

	"github.com/rs/zerolog/log"
)

func StartListener(store *Store) {
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

		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store *Store) {
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

		res, err := runExpression(message, store)
		if err != nil {
			log.Error().Str("remote_addr", conn.RemoteAddr().String()).Err(err).Msg("Error trying to run expression")
		}

		// TODO: send error messages to client like when getting a key that does not exist

		conn.Write(res)
	}
}
