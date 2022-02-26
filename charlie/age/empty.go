package age

import "io"

type ageEmpty struct {
}

func MakeEmpty() (a *ageEmpty) {
	return &ageEmpty{}
}

func (a ageEmpty) Recipient() _AgeRecipient {
	return nil
}

func (a ageEmpty) Identity() _AgeIdentity {
	return nil
}

func (a ageEmpty) Decrypt(src io.Reader) (io.Reader, error) {
	return src, nil
}

type writeCloser struct {
	io.Writer
}

func (w writeCloser) Close() (err error) {
	return
}

func (a ageEmpty) Encrypt(dst io.Writer) (io.WriteCloser, error) {
	return writeCloser{dst}, nil
}
