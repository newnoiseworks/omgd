local nk = require("nakama")
local DM = require("dungeon_manager")
local time = require("time")
local wallet = require("wallet")

local M = {}

local empty_room_ticks = 0
local max_empty_room_ticks = 15

local matchStartTimestamp = nk.time() / 1000
local gameDaysFromInit = 0

function M.match_init(context, setupstate)
  local gamestate = {
    presences = {},
    plot_map = {},
    avatar_map = {},
    username_map = {},
    dungeon_map = {},
    weather = {}
  }

  local tickrate = 1 -- per sec

  nk.logger_info("Realm started")

  gameDaysFromInit = time.number_of_game_days_from_daybreak_unix_timestamp(matchStartTimestamp)

  gamestate = make_random_weather_change_event(gamestate)

  return gamestate, tickrate, setupstate.label
end

function M.match_join_attempt(context, dispatcher, tick, state, presence, metadata)
  local acceptuser = true

  local presence_count = get_presence_count(state)

  if (presence_count >= DM.max_dungeon_size) then
    acceptuser = false
  else
    set_plot_id(state, presence)
    state.presences[presence.session_id] = presence
    state.avatar_map[presence.user_id] = metadata.avatarName
    state.username_map[presence.user_id] = metadata.username
    -- TODO: Initial dungeon should pipe in from metadata of realm join (this will be important when adding invite stuff)
    state.dungeon_map[presence.user_id] = metadata.avatarName .. "'s Farm"
  end

  return state, acceptuser
end

function M.match_join(context, dispatcher, tick, state, presences)
  --for _, presence in pairs(presences) do
  --end

  broadcast_state(dispatcher, state)

  return state
end

function M.match_leave(context, dispatcher, tick, state, presences)
  for _, presence in pairs(presences) do
    state.presences[presence.session_id] = nil
    state.plot_map[presence.user_id] = nil
    state.avatar_map[presence.user_id] = nil
    state.username_map[presence.user_id] = nil
  end

  broadcast_state(dispatcher, state)

  return state
end

function broadcast_state(dispatcher, state)
  dispatcher.broadcast_message(
    0,
    nk.json_encode(
      {
        plotMap = state.plot_map,
        avatarMap = state.avatar_map,
        dungeonMap = state.dungeon_map,
        usernameMap = state.username_map,
        weather = state.weather
      }
    ),
    nil
  )
end

function set_plot_id(state, presence)
  local plot_id = 0
  local taken_plots = {}

  for user, plot in pairs(state.plot_map) do
    taken_plots[plot] = plot
  end

  for i = 0, DM.max_dungeon_size, 1 do
    if (taken_plots[i] == nil) then
      plot_id = i
      break
    end
  end

  state.plot_map[presence.user_id] = plot_id
end

function get_presence_count(state)
  local presence_count = 0
  for _, presence in pairs(state.presences) do
    -- nk.logger_info(("Presence %s named %s"):format(presence.user_id, presence.username))
    presence_count = presence_count + 1
  end

  return presence_count
end

function should_room_be_open(state)
  presence_count = get_presence_count(state)

  if presence_count == 0 then
    -- nk.logger_info(("match room empty, tick count: %s"):format(empty_room_ticks))
    empty_room_ticks = empty_room_ticks + 1

    if empty_room_ticks >= max_empty_room_ticks then
      nk.logger_info("realm room empty for too long, terminating...")
      return false
    end
  elseif empty_room_ticks > 0 then
    -- nk.logger_info("match room not empty, zeroing out empty_room_ticks")
    empty_room_ticks = 0
  end

  return true
end

function M.match_loop(context, dispatcher, tick, state, messages)
  if should_room_be_open(state) == false then
    return nil
  end

  for _, message in ipairs(messages) do
    dispatcher.broadcast_message(message.op_code, message.data, nil, message.sender)
  end

  for _, message in ipairs(messages) do
    state = perform_message_validation(message, dispatcher, state)
  end

  local _gameDaysFromInit = time.number_of_game_days_from_daybreak_unix_timestamp(matchStartTimestamp)

  if (gameDaysFromInit < _gameDaysFromInit) then
    -- nk.logger_info("game day advanced")
    gameDaysFromInit = _gameDaysFromInit
    for _, presence in pairs(state.presences) do
      local avatarKey = state.avatar_map[presence.user_id]
      wallet.fund(presence.user_id, avatarKey, "CorpusCoin", 5, {}, true)
      -- nk.logger_info("daily play delivered to user id " .. presence.user_id .. " and avatarKey " .. avatarKey)
    end

    dispatcher.broadcast_message(
      2, -- Wallet update event
      nk.json_encode({}),
      nil
    )

    state = make_random_weather_change_event(state)

    dispatcher.broadcast_message(
      3, -- Weather change event
      nk.json_encode(state.weather),
      nil
    )
  end

  return state
end

function make_random_weather_change_event(state)
  local rain_level = math.random()
  local cloud_level = math.random()
  local sun_level = math.random()

  if rain_level < 0.66 and math.random() < 0.80 then rain_level = 0 end

  if rain_level == 0 and cloud_level < 0.44 then cloud_level = 0 end

  if sun_level < 0.66 and math.random() < 0.80 then sun_level = 1 end

  local weather = {
    rainLevel = rain_level,
    sunLevel = sun_level,
    cloudLevel = cloud_level
  }

  state.weather = weather

  return state
end

function perform_message_validation(message, dispatcher, state)
  if message.op_code == 1 then -- 1 == ChangeDungeon event
    local data = nk.json_decode(message.data)
    state.dungeon_map[message.sender.user_id] = data.dungeon
  end

  return state
end

function M.match_terminate(context, dispatcher, tick, state, grace_seconds)
  nk.logger_info("terminating realm")
  local message = "Server shutting down realm in " .. grace_seconds .. " seconds"
  -- TODO: implement message and opcode for server messaging
  -- dispatcher.broadcast_message(2, message)
  return nil
end

return M
