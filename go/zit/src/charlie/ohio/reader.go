package ohio

import (
	"encoding/binary"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func ReadAllOrDieTrying(r io.Reader, b []byte) (n int, err error) {
	var acc int

	for n < len(b) {
		acc, err = r.Read(b[n:])
		n += acc

		if err != nil {
			break
		}
	}

	switch err {
	case io.EOF:
		if n < len(b) && n > 0 {
			err = errors.Wrapf(
				io.ErrUnexpectedEOF,
				"Expected %d, got %d",
				len(b),
				n,
			)
		}
	case nil:
	default:
		err = errors.Wrap(err)
	}

	return
}

func ReadUint8(r io.Reader) (n uint8, read int, err error) {
	cl := [1]byte{}

	read, err = ReadAllOrDieTrying(r, cl[:])
	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	clInt, clIntErr := binary.Uvarint(cl[:])

	if clIntErr <= 0 {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n = uint8(clInt)

	return
}

func ReadFixedUint8(r io.Reader) (n uint8, read int, err error) {
	cl := [1]byte{}

	read, err = ReadAllOrDieTrying(r, cl[:])
	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n = cl[0]

	return
}

func ReadInt8(r io.Reader) (n int8, read int, err error) {
	cl := [1]byte{}

	read, err = ReadAllOrDieTrying(r, cl[:])
	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	clInt, clIntErr := binary.Uvarint(cl[:])

	if clIntErr <= 0 {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n = int8(clInt)

	return
}

func ReadUint16(r io.ByteReader) (v uint16, n int64, err error) {
	var clInt uint64

	if clInt, err = binary.ReadUvarint(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	v = uint16(clInt)
	n = int64(binary.Size(v))

	return
}

func ReadInt64(r io.Reader) (n int64, read int, err error) {
	cl := [binary.MaxVarintLen64]byte{}

	read, err = ReadAllOrDieTrying(r, cl[:])
	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	n, clIntErr := binary.Varint(cl[:])

	if clIntErr <= 0 {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}

func ReadFixedInt32(r io.Reader) (n int, val int32, err error) {
	cl := [4]byte{}

	n, err = ReadAllOrDieTrying(r, cl[:])
	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	val = ByteArrayToInt32(cl)

	return
}

func ReadFixedInt64(r io.Reader) (n int, val int64, err error) {
	cl := [8]byte{}

	n, err = ReadAllOrDieTrying(r, cl[:])
	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	val = ByteArrayToInt64(cl)

	return
}

func MakeLineReaderIterateStrict(
	rffs ...interfaces.FuncSetString,
) interfaces.FuncSetString {
	si, _ := errors.MakeStackInfo(1)
	var i int64

	return func(v string) (err error) {
		if int64(len(rffs))-1 < i {
			err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
				error:  err,
				string: v,
			})

			return
		}

		if err = rffs[i](v); err != nil {
			err = si.Wrapf(err, "Value: %s", v)
			return
		}

		i++

		return
	}
}

func MakeLineReaderIterate(
	rffs ...interfaces.FuncSetString,
) interfaces.FuncSetString {
	si, _ := errors.MakeStackInfo(1)
	var i int64

	return func(v string) (err error) {
		for {
			if int64(len(rffs))-1 < i {
				err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
					error:  err,
					string: v,
				})

				return
			}

			if err = rffs[i](v); err != nil {
				i++
				err = si.Wrapf(err, "Value: %s", v)
				continue
			}

			return
		}
	}
}

func MakeLineReaderKeyValues(
	dict map[string]interfaces.FuncSetString,
) interfaces.FuncSetString {
	si, _ := errors.MakeStackInfo(1)

	return func(line string) (err error) {
		loc := strings.Index(line, " ")

		if loc == -1 {
			err = si.Errorf(
				"expected at least one space, but found none: %q",
				line,
			)
			return
		}

		key := line[:loc]
		value := line[loc+1:]

		var reader interfaces.FuncSetString
		ok := false

		if reader, ok = dict[key]; !ok {
			err = si.Errorf("key not supported: %q", key)
			return
		}

		if err = reader(value); err != nil {
			err = si.Errorf("%s: %q", err, value)
			return
		}

		return
	}
}

func MakeLineReaderKeyValue(
	key string,
	valueReader interfaces.FuncSetString,
) interfaces.FuncSetString {
	return func(line string) (err error) {
		loc := strings.Index(line, " ")

		if loc == -1 {
			err = errors.Errorf(
				"expected at least one space, but found none: %q",
				line,
			)
			return
		}

		keyActual := line[:loc]
		value := line[loc+1:]

		if keyActual != key {
			err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
				string: value,
			})

			return
		}

		if err = valueReader(value); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeLineReaderRepeat(
	in interfaces.FuncSetString,
) interfaces.FuncSetString {
	return func(line string) (err error) {
		if err = in(line); err != nil {
			err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
				error:  err,
				string: line,
			})

			return
		}

		return
	}
}

func MakeLineReaderIgnoreErrors(
	in interfaces.FuncSetString,
) interfaces.FuncSetString {
	return func(line string) (err error) {
		in(line)

		return
	}
}

func MakeLineReaderNop() interfaces.FuncSetString {
	return func(line string) (err error) {
		return
	}
}
