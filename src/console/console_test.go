package console

import(
	"testing"
	"os"
	"path"
	"fmt"
	"regexp"
	"container/vector"
	pt "spectrum/prettytest"
	"spectrum"
)

const testdataDir = "testdata"
var screenshotFn = path.Join(testdataDir, "screenshot.scr")

type testMessageOutput struct {
	strings *vector.StringVector
}

var messageOutput *testMessageOutput

func (out *testMessageOutput) Clear() {
	out.strings = new(vector.StringVector)
}

func (out *testMessageOutput) String() string {
	outputString := ""
	for _, str := range *out.strings {
		outputString += str
	}
	return outputString
}

func (out *testMessageOutput) PrintfMsg(format string, a ...interface{}) {
	out.strings.Push(fmt.Sprintf(format, a...))

	appendNewLine := false
	if (len(format) == 0) || (format[len(format)-1] != '\n') {
		appendNewLine = true
	}

	if appendNewLine {
		out.strings.Push("\n")
	}

}

func removeTestdataFiles() {
	os.Remove(screenshotFn)
}

func beforeAll(t *pt.T) {
	var err os.Error
	app := spectrum.NewApplication()
	speccy, err = spectrum.NewSpectrum48k(app, "testdata/48.rom")
	if err != nil {
		panic(err)
	}

	app.SetMessageOutput(&testMessageOutput{ new(vector.StringVector) })
	messageOutput = app.GetMessageOutput().(*testMessageOutput)
	ignoreStartupScript = true

	Init(app, speccy)
}

func before(t *pt.T) {
	messageOutput.Clear()
}

func after(t *pt.T) {
	removeTestdataFiles()
}

func testCommandHelp(t *pt.T) {
	err := run(w, "help()")
	matched, _ := regexp.MatchString("Available commands:\n *.* *\n", messageOutput.String())

	t.True(err == nil)
	t.True(matched)
}

func testCommandScreenshot(t *pt.T) {
	err := run(w, fmt.Sprintf("screenshot(\"%s\")", screenshotFn))

	t.True(err == nil)
	t.Path(screenshotFn)
}

func testCommandScript(t *pt.T) {
	t.Pending()
}

func TestCommands(t *testing.T) {
	pt.Run(
		t,
		beforeAll,
		before,
		after,
		testCommandHelp,
		testCommandScreenshot,
		testCommandScript,
	)
}