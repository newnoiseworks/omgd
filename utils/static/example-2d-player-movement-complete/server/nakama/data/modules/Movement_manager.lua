local nk = require("nakama")

local M = {}

M.max_Movement_size = 8

local function find_or_create_nk_match(context, payload, type)
  local limit = 1
  local isauthoritative = true
  local label = payload
  local min_size = 0
  local max_size = M.max_Movement_size - 1
  local matches = nk.match_list(1, isauthoritative, label, min_size, max_size)

  local matchid = nil

  for _, match in ipairs(matches) do
    matchid = match.match_id
  end

  if (matchid == nil) then
    matchid =
      nk.match_create(
      type,
      {
        label = label
      }
    )
  end

  return nk.json_encode(matchid)
end

local function find_or_create_Movement(context, payload)
  return find_or_create_nk_match(context, payload, "Movement")
end

nk.register_rpc(find_or_create_Movement, "find_or_create_Movement")

return M
