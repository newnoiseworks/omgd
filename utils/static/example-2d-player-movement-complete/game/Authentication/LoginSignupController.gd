extends Container

var username: String
var email: String
var password: String

onready var version_incorrect_modal: WindowDialog = $VersionIncorrect

export var game_scene: PackedScene


func username_text_changed(new_text: String):
	username = new_text


func email_text_changed(new_text: String):
	email = new_text


func password_text_changed(new_text: String):
	password = new_text


func login_button_pressed():
	var session = yield(SessionManager.login(email, password), "completed")

	handle_post_login(session.valid)


func sign_up_button_pressed():
	# var session = yield(SessionManager.signup(email, username, password), "completed")

	var session = yield(
		SessionManager.signup("%s@wow.com" % username, username, "password"), "completed"
	)

	handle_post_login(session.valid)


func handle_post_login(login_successful: bool):
	if login_successful:
		get_parent().get_parent().queue_free()
		get_tree().get_root().call_deferred("add_child", game_scene.instance())
		print("logged in!")
	else:
		print("auth failed!")
