package calc

import (
	"github.com/paketo-buildpacks/libjvm/helper"
	"os"
	"strconv"
)

const (
	DefaultHeadRoom = HeadRoom(helper.DefaultHeadroom)
	FlagHeadRoom    = "head-room"
	EnvHeadRoom     = "BPL_JVM_HEAD_ROOM"
	UsageHeadRoom   = "percentage of total memory available which will be left unallocated to cover JVM overhead"
)

type HeadRoom int

func NewHeadRoom() *HeadRoom {
	hr := DefaultHeadRoom
	if val, ok := os.LookupEnv(EnvHeadRoom); ok {
		f, _ := strconv.Atoi(val)
		hr = HeadRoom(f)
	}
	return &hr
}

func (hr *HeadRoom) Set(s string) error {
	f, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*hr = HeadRoom(f)
	return nil
}

func (hr *HeadRoom) Type() string {
	return "int"
}

func (hr *HeadRoom) String() string {
	return strconv.FormatInt(int64(*hr), 10)
}

func (hr *HeadRoom) Contribute() error {
	if err := os.Setenv(EnvHeadRoom, hr.String()); err != nil {
		return err
	}
	return nil
}
