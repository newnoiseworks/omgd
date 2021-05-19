extends Node

var session: NakamaSession

var api_account: NakamaAPI.ApiAccount

onready var client = Nakama.create_client(
	GameConfig.nakama_key,
	GameConfig.nakama_host,
	GameConfig.nakama_port,
	"https" if GameConfig.nakama_secure else "http"
)


func rpc_async(uri: String, payload: String = "") -> NakamaAPI.ApiRpc:
	return yield(client.rpc_async(session, uri, payload), "completed")


func login(email: String, password: String):
	if client == null:
		client = Nakama.create_client(
			GameConfig.nakama_key,
			GameConfig.nakama_host,
			GameConfig.nakama_port,
			"https" if GameConfig.nakama_secure else "http"
		)

	session = yield(client.authenticate_email_async(email, password), "completed")
	yield(post_auth(), "completed")

	return session


func signup(email: String, username: String, password: String):
	if client == null:
		client = Nakama.create_client(
			GameConfig.nakama_key,
			GameConfig.nakama_host,
			GameConfig.nakama_port,
			"https" if GameConfig.nakama_secure else "http"
		)

	session = yield(client.authenticate_email_async(email, password, username, true), "completed")
	yield(post_auth(), "completed")

	return session


func post_auth():
	yield(MatchManager.connect_socket(), "completed")
