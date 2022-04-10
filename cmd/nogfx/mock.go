package main

import (
	"io"
	"strings"
)

type MockClient struct {
	reader io.Reader
	writer io.Writer
}

func (mock *MockClient) Read(p []byte) (int, error) {
	return mock.reader.Read(p)
}

func (mock MockClient) Write(p []byte) (int, error) {
	return mock.writer.Write(p)
}

func mockReadWriter() io.ReadWriter {
	return &MockClient{
		strings.NewReader("trololol\nqweqwrreqr\none two \033[33mthree \033[39mfour five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty twentyone twentytwo twentythree twentyfour twentyfive twentysix twentyseven twentyeight twentynine thirty thirtyone thirtytwo\nzxcxzvzxcxcxzc"),
		&strings.Builder{},
	}
}
