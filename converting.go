package main

import "strconv"

const minute = 60
const hour = minute * 60
const day = hour * 24

type left struct {
	Days    uint64
	Hours   uint64
	Minutes uint64
	Seconds uint64
}

func convertSecondsToTime(seconds uint64) left {
	ret := left{Days: 0, Hours: 0, Minutes: 0, Seconds: 0}

	if seconds == 0 {
		return ret
	}

	if seconds > day {
		ret.Days = seconds / day
		seconds = seconds - (seconds/day)*day
		if seconds == 0 {
			return ret
		}
	}
	if seconds > hour {
		ret.Hours = seconds / hour
		seconds = seconds - (seconds/hour)*hour
		if seconds == 0 {
			return ret
		}
	}
	if seconds > minute {
		ret.Minutes = seconds / minute
		seconds = seconds - (seconds/minute)*minute
		if seconds == 0 {
			return ret
		}
	}
	ret.Seconds = seconds

	return ret
}

func (lf left) ToString() string {
	timeLeftString := ""

	if lf.Days != 0 {
		timeLeftString += strconv.FormatUint(lf.Days, 10) + " Days "
	}
	if lf.Hours != 0 {
		timeLeftString += strconv.FormatUint(lf.Hours, 10) + " Hours "
	}
	if lf.Minutes != 0 {
		timeLeftString += strconv.FormatUint(lf.Minutes, 10) + " Minutes "
	}
	if lf.Seconds != 0 {
		timeLeftString += strconv.FormatUint(lf.Seconds, 10) + " Seconds"
	}

	return timeLeftString
}

const kiloByte = 1024
const MegaByte = kiloByte * 1024
const GigaByte = MegaByte * 1024

type Size struct {
	GigaBytes int
	MegaBytes int
	KiloBytes int
	Bytes     int
}

func convertBytesToSize(bytesSize int) Size {
	ret := Size{GigaBytes: 0, MegaBytes: 0, KiloBytes: 0, Bytes: 0}

	if bytesSize == 0 {
		return ret
	}

	if bytesSize > GigaByte {
		ret.GigaBytes = bytesSize / GigaByte
		bytesSize = bytesSize - (bytesSize/GigaByte)*GigaByte
		if bytesSize == 0 {
			return ret
		}
	}
	if bytesSize > MegaByte {
		ret.MegaBytes = bytesSize / MegaByte
		bytesSize = bytesSize - (bytesSize/MegaByte)*MegaByte
		if bytesSize == 0 {
			return ret
		}
	}
	if bytesSize > kiloByte {
		ret.KiloBytes = bytesSize / kiloByte
		bytesSize = bytesSize - (bytesSize/kiloByte)*kiloByte
		if bytesSize == 0 {
			return ret
		}
	}
	ret.Bytes = bytesSize

	return ret
}

func (sz Size) ToString() string {
	sizeString := ""

	if sz.GigaBytes != 0 {
		sizeString += strconv.Itoa(sz.GigaBytes) + " Gb "
	}
	if sz.MegaBytes != 0 {
		sizeString += strconv.Itoa(sz.MegaBytes) + " Mb "
	}
	if sz.KiloBytes != 0 {
		sizeString += strconv.Itoa(sz.KiloBytes) + " Kb "
	}
	if sz.Bytes != 0 {
		sizeString += strconv.Itoa(sz.Bytes) + " Bytes"
	}

	return sizeString
}
