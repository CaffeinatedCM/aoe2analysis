package parser

type RecordedGame struct {
	Header Header
}

/*
Header information found in aoe2record.

Note: Specifically uses int32 for integers to be consistent with the file format since int appears to be more flexible
in go with it's definition of "at least 32 bits"
 */
type Header struct {
	length	int32
	Version string
}
