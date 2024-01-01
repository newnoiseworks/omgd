---
title: "Debugging"
---

# Debugging

To debug any of the OMGD commands, you'll need to setup an environment flag.

`OMGD_LOG_LEVEL` can be set to `DEBUG` or `TRACE` for more verbose output of any `omgd` commands that have been run.

In Windows, to setup an environment variable for one command, try `set OMGD_LOG_LEVEL=DEBUG & omgd ... & set OMGD_LOG_LEVEL=`. This will set the log level and then unset it afterwards.

In MacOS / Unix, try `OMGD_LOG_LEVEL=DEBUG omgd ...`.
