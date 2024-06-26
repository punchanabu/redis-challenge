package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/config"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func HandleCommand(command string, argument []string, store *store.Store, config *config.ReplicaConfig) string {

	lowerCommand := strings.ToLower(command)

	switch lowerCommand {
	case "ping":
		return handlePing()
	case "echo":
		return handleEcho(argument)
	case "get":
		return handleGet(argument, store)
	case "set":
		return handleSet(argument, store)
	case "info":
		return handleInfo(argument, config)
	case "replconf":
		return handleReplicationConfig()
	case "psync":
		return handlePsync(config)
	default:
		return "-ERR unknown command"
	}
}

func handlePing() string {
	return "+PONG"
}

func handleEcho(argument []string) string {
	if len(argument) == 0 {
		return "-ERR no argument provided"
	}
	return "+" + argument[0]
}

func handleGet(argument []string, store *store.Store) string {
	if len(argument) == 0 {
		return "-ERR no argument provided"
	}
	value, ok := store.Get(argument[0])
	/*
		If there is no value returns an empty string
		as it will be encoded as a $-1 response.
	*/
	if !ok {
		return ""
	}
	return "+" + value
}

func handleSet(argument []string, store *store.Store) string {
	// Check if there are at least key and value arguments
	if len(argument) < 2 {
		return "-ERR not enough arguments"
	}

	var expiryMillis int64 = 0 // Default of no expiry time
	// Check if the optional expiration time is provided
	if len(argument) > 2 {
		// Check if the 'PX' expiration time is provided
		fmt.Println(strings.ToUpper(argument[2]), " ", argument[3])
		if len(argument) == 4 && strings.ToUpper(argument[2]) == "PX" {
			var err error
			expiryMillis, err = strconv.ParseInt(argument[3], 10, 64)
			if err != nil {
				return "-ERR invalid expiration time"
			}
		} else {
			return "-ERR wrong number of arguments for 'set' command or wrong syntax"
		}
	}

	// Perform the Set operation
	store.Set(argument[0], argument[1], expiryMillis)
	return "+OK"
}

func handleInfo(argument []string, config *config.ReplicaConfig) string {

	if len(argument) == 0 {
		return "-ERR no argument provided"
	}

	section := strings.ToLower(argument[0])
	switch section {
	case "replication":
		return formatReplicationInfo(config)
	default:
		return "-ERR unsupported INFO section"
	}
}

func formatReplicationInfo(config *config.ReplicaConfig) string {
	infoLines := []string{
		"role:" + config.Role,
		"connected_slaves:" + strconv.Itoa(config.ConnectedSlaves),
		"master_replid:" + config.MasterReplID,
		"master_repl_offset:" + strconv.Itoa(config.MasterReplOffset),
		"second_repl_offset:" + strconv.Itoa(config.SecondReplOffset),
		"repl_backlog_active:" + strconv.Itoa(config.ReplBacklogActive),
		"repl_backlog_size:" + strconv.Itoa(config.ReplBacklogSize),
		"repl_backlog_first_byte_offset:" + strconv.Itoa(config.ReplBacklogFirstByteOffset),
		"repl_backlog_histlen:" + strconv.Itoa(config.ReplBacklogHistLen),
	}
	info := strings.Join(infoLines, "\r\n")
	return "$" + strconv.Itoa(len(info)) + "\r\n" + info + "\r\n"
}

func handleReplicationConfig() string {
	return "+OK"
}

func handlePsync(config *config.ReplicaConfig) string {
	return fmt.Sprintf("+FULLRESYNC %s 0", config.MasterReplID)
}
