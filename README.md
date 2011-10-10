About this package
==================

This package is designed to support better logging for Go. Specifically, this
project aims to support different levels of logging and the ability to
customize log output via custom implementations of the Logger interface.

Using this package
==================

Understanding this package
==========================

There are 3 important objects in this package.
	LogOuter: Outputs a LogMessage to (file, testing.T, network, xml, etc)
	Logger: Decides whether on not to generate output
	LevelLogger: Easier interface for Logger.

In practice, the user is encouraged to use the LevelLogger as an entrypoint into
the package. The provided Global LevelLogger is set up to have easy defaults
and to be easily configurable with flags and the AddLogFile and the
{Start,Stop}TestLogging functions. As an alternative, the user can create
package specific LevelLogger with their own presets or the default (flag based)
presets.

NOTE: The package is not quite stable. Most exported methods and types will
remain exported, but may change.
