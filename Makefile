include $(GOROOT)/src/Make.inc

TARG=golog
GOFILES=\
	log_outer.go\
	golog.go\
	internal.go\
	vmodules.go\
	logger.go\
	level_logger.go\
	multi_log_outer.go\

include $(GOROOT)/src/Make.pkg
