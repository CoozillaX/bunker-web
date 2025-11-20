# COPIED FROM https://github.com/CMA2401PT/neomega/blob/main/omega_launcher/Makefile

# define go build flags
GO_CGO_FLAGS_COMMON:=CGO_ENABLED=0
GO_BUILD_FLAGS_COMMON:=-trimpath -ldflags "-s -w" -tags "jsoniter"
# end define go build flags

ROOT_DIR:= $(shell pwd)
SRCS_GO := $(foreach dir, $(shell find $(ROOT_DIR) -type d), $(wildcard $(dir)/*.go $(dir)/*.c))

RELEASE_DIR:=$(shell pwd)/build
OUTPUT_DIR:=${RELEASE_DIR}
EXEC_PREFIX:=phoenixauth-

TYPE:=linux
${TYPE}_EXEC:=${OUTPUT_DIR}/${EXEC_PREFIX}${TYPE}
${${TYPE}_EXEC}_TRIPLE:=GOOS=linux GOARCH=amd64
${TYPE}: ${${TYPE}_EXEC}
EXECS:=${EXECS} ${${TYPE}_EXEC}

TYPE:=linux-arm64
${TYPE}_EXEC:=${OUTPUT_DIR}/${EXEC_PREFIX}${TYPE}
${${TYPE}_EXEC}_TRIPLE:=GOOS=linux GOARCH=arm64
${TYPE}: ${${TYPE}_EXEC}
EXECS:=${EXECS} ${${TYPE}_EXEC}

TYPE:=windows-x86
${TYPE}_EXEC:=${OUTPUT_DIR}/${EXEC_PREFIX}${TYPE}.exe
${${TYPE}_EXEC}_TRIPLE:=GOOS=windows GOARCH=386
${TYPE}: ${${TYPE}_EXEC}
EXECS:=${EXECS} ${${TYPE}_EXEC}

TYPE:=windows
${TYPE}_EXEC:=${OUTPUT_DIR}/${EXEC_PREFIX}${TYPE}.exe
${${TYPE}_EXEC}_TRIPLE:=GOOS=windows GOARCH=amd64
${TYPE}: ${${TYPE}_EXEC}
EXECS:=${EXECS} ${${TYPE}_EXEC}

TYPE:=macos
${TYPE}_EXEC:=${OUTPUT_DIR}/${EXEC_PREFIX}${TYPE}
${${TYPE}_EXEC}_TRIPLE:=GOOS=darwin GOARCH=amd64
${TYPE}: ${${TYPE}_EXEC}
EXECS:=${EXECS} ${${TYPE}_EXEC}

TYPE:=macos-arm64
${TYPE}_EXEC:=${OUTPUT_DIR}/${EXEC_PREFIX}${TYPE}
${${TYPE}_EXEC}_TRIPLE:=GOOS=darwin GOARCH=arm64
${TYPE}: ${${TYPE}_EXEC}
EXECS:=${EXECS} ${${TYPE}_EXEC}


${OUTPUT_DIR}:
	@echo make output dir $@
	@mkdir -p $@
	
info:
	echo ${EXECS}
.PHONY: ${EXECS}
${EXECS}: ${OUTPUT_DIR}/${EXEC_PREFIX}%: ${OUTPUT_DIR} ${SRCS_GO}
	@${GO_CGO_FLAGS_COMMON} ${$@_TRIPLE}  go build ${GO_BUILD_FLAGS_COMMON} -o $@ ${ROOT_DIR}
	@echo build $@

execs:${EXECS}

all: ${EXECS}

clean:
	rm -f ${OUTPUT_DIR}/*