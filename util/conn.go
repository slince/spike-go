package util

import "io"

func Pipe(to io.Writer, from io.Reader, bytesCopied *int64) error{
	var err error

	*bytesCopied, err = io.Copy(to, from)

	if err != nil {
		return err
	}

}
