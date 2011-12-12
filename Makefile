include $(GOROOT)/src/Make.inc

TARG=gogress

GOFILES=gogress.go

all: package 

include $(GOROOT)/src/Make.pkg
