package core

import (
	"strings"
	"zdb/redis/internal/resp"
)

func EvalCommand(args []resp.Value, store *Store) *resp.Value {
	if len(args) == 0 {
		return &resp.Value{
			Type: "error",
			Str:  "ERR Empty command",
		}
	}

	command := strings.ToUpper(args[0].Bulk)

	switch command {
	case "PING":
		return evalPing(args)
	case "ECHO":
		return evalEcho(args)
	case "GET":
		return evalGet(args, store)
	case "SET":
		return evalSet(args, store)
	default:
		return &resp.Value{
			Type: "error",
			Str:  "ERR Unknown command",
		}
	}
}

func evalPing(args []resp.Value) *resp.Value {
	if len(args) == 2 {
		param := args[1].Bulk
		return &resp.Value{
			Type: "bulk",
			Bulk: param,
		}
	}

	return &resp.Value{
		Type: "string",
		Str:  "PONG",
	}
}

func evalEcho(args []resp.Value) *resp.Value {
	if len(args) != 2 {
		return &resp.Value{
			Type: "error",
			Str:  "ERR wrong number of arguments for 'echo' command",
		}
	}

	param := args[1].Bulk

	return &resp.Value{
		Type: "bulk",
		Bulk: param,
	}
}

func evalGet(args []resp.Value, store *Store) *resp.Value {
	if len(args) != 2 {
		return &resp.Value{
			Type: "error",
			Str:  "ERR wrong number of arguments for 'get' command",
		}
	}

	key := args[1].Bulk

	obj := store.Get(key)

	if obj != nil && obj.Type == ObjectTypeString {
		return &resp.Value{
			Type: "bulk",
			Bulk: obj.Value.(string),
		}
	}

	return &resp.Value{
		Type: "null",
	}
}

func evalSet(args []resp.Value, store *Store) *resp.Value {
	if len(args) != 3 {
		return &resp.Value{
			Type: "error",
			Str:  "ERR wrong number of arguments for 'set' command",
		}
	}

	key := args[1].Bulk
	value := args[2].Bulk

	store.Put(key, &RedisObject{
		Type:      ObjectTypeString,
		Value:     value,
		ExpiresAt: -1,
	})
	return &resp.Value{
		Type: "string",
		Str:  "OK",
	}
}
