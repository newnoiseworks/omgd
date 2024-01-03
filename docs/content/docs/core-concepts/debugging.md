---
title: "Debugging"
---

# Debugging

To debug any of the OMGD commands, you'll need to setup an environment flag.

`OMGD_LOG_LEVEL` can be set to `DEBUG` or `TRACE` for more verbose output of any `omgd` commands that have been run.

In Windows, to setup an environment variable for one command, try `set OMGD_LOG_LEVEL=DEBUG & omgd ... & set OMGD_LOG_LEVEL=`. This will set the log level and then unset it afterwards.

In MacOS / Unix, try `OMGD_LOG_LEVEL=DEBUG omgd ...`.

## All OMGD_LOG_LEVEL Values

| OMGD_LOG_LEVEL | Intent |
| --- | --- |
| TRACE | The highest level of debugging. Will output as much information as possible. |
| DEBUG | Will output information for debugging, as well as output from `docker`, `terraform`, and other interdependent CLI commands. |
| INFO | The default log level. Outputs basic information on commands. |
| WARN | Will only output warnings. |
| ERROR | Will only output error information. |
| FATAL | Will only provide output on fatal program exits. |
