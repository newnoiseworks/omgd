local nk = require("nakama")
local DM = require("player_manager")

local M = {}

local empty_room_ticks = 0
local max_empty_room_ticks = 150

function M.match_init(context, setupstate)
  local gamestate = {
    presences = {},
  }

  local tickrate = 15 -- per sec

  nk.logger_info("player started")

  return gamestate, tickrate, setupstate.label
end

function M.match_join_attempt(context, dispatcher, tick, state, presence, metadata)
  local acceptuser = true

  local presence_count = get_presence_count(state)

  if (presence_count >= DM.max_player_size) then
    acceptuser = false
  else
    state.presences[presence.session_id] = presence
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
  end

  broadcast_state(dispatcher, state)

  return state
end

function broadcast_state(dispatcher, state)
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
      nk.logger_info("match room empty for too long, terminating...")
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

  return state
end

function perform_message_validation(message, dispatcher, state)
  -- if message.op_code == 0 then
  -- end

  return state
end

function M.match_terminate(context, dispatcher, tick, state, grace_seconds)
  nk.logger_info("terminating match")
  local message = "Server shutting down match in " .. grace_seconds .. " seconds"
  return nil
end

function M.match_signal(context, dispatcher, tick, state, data)
  return nil
end

return M
