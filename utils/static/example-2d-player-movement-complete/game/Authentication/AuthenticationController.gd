extends Control

onready var login_container: Container = $Login
onready var signup_container: Container = $Signup
onready var label: Label = $Label


func _ready():
	label.text = "v%s" % GameConfig.version

	var cmd_args = OS.get_cmdline_args()

	if cmd_args.size() > 1 && "email" in cmd_args[0] && "password" in cmd_args[1]:
		# await SessionManager.Signup(
		#   cmdArgs[0].Split("=")[1], null, cmdArgs[1].Split("=")[1]
		# );
		print("signed in!!")


func show_login():
	show_login_or_signup(true)


func show_signup():
	show_login_or_signup(false)


func show_login_or_signup(show_login: bool):
	if show_login:
		login_container.show()
		signup_container.hide()
	else:
		login_container.hide()
		signup_container.show()
