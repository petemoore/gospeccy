package formats

import (
	"testing"
	"spectrum/prettytest"
	"path"
	"io/ioutil"
	"bytes"
)

const testdataDir = "testdata"

var (
	tapCodeFn = path.Join(testdataDir, "code.tap")
	tapProgramFn = path.Join(testdataDir, "hello.tap")
)

func testReadTAP(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	n, err := tap.Read(data)

	headerBlock := tap.blocks.At(0).(*tapBlockHeader)
	dataBlock := tap.blocks.At(1).(tapBlockData)

	assert.Nil(err)
	assert.Equal(27, n)

	assert.True(headerBlock != nil)
	assert.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
	assert.Equal("ROM       ", headerBlock.filename)
	assert.Equal(uint16(2), headerBlock.length)
	assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
}

func testReadTAPError(assert *prettytest.T) {
	tap := NewTAP()
	_, err := tap.Read(nil)
	assert.False(err == nil)
}

// SAVE "ROM" CODE 0,2
func testReadTAPCodeFile(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapCodeFn)
	tap := NewTAP()
	_, err := tap.Read(data)

	assert.Nil(err)

	if !assert.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		assert.Equal(byte(TAP_FILE_CODE), headerBlock.tapType)
		assert.Equal(uint16(0), headerBlock.par1)
		assert.Equal(uint16(0x8000), headerBlock.par2)

		assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		assert.Equal(byte(0xf3), dataBlock[1])
		assert.Equal(byte(0xaf), dataBlock[2])
		assert.Equal(byte(0xa3), dataBlock[3])
	}
}

// 10 PRINT "Hello World"
// SAVE "HELLO"
func testReadTAPProgramFile(assert *prettytest.T) {
	data, _ := ioutil.ReadFile(tapProgramFn)
	tap := NewTAP()
	_, err := tap.Read(data)

	assert.Nil(err)
	
	if !assert.Failed() {
		headerBlock := tap.blocks.At(0).(*tapBlockHeader)
		dataBlock := tap.blocks.At(1).(tapBlockData)

		assert.Equal(byte(TAP_FILE_PROGRAM), headerBlock.tapType)
		assert.Equal(uint16(0x8000), headerBlock.par1)
		assert.Equal(uint16(0x14), headerBlock.par2)

		assert.Equal(byte(TAP_BLOCK_DATA), dataBlock[0])
		assert.True(bytes.Equal([]byte{
			0x00, 0x0a, 
			0x10, 0x00,
			0x20, 0xf5,
			0x22, 0x48,
			0x65, 0x6c,
			0x6c, 0x6f,
			0x20, 0x57,
			0x6f, 0x72,
			0x6c, 0x64,
			0x22, 0x0d,
			0x1d,
		}, dataBlock[1:]))
	}
}

func TestLoadTAP(t *testing.T) {
	prettytest.Run(
		t,
		testReadTAP,
		testReadTAPError,
		testReadTAPCodeFile,
		testReadTAPProgramFile,
	)
}
