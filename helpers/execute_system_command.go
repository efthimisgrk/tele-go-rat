package helpers

import (
	"os/exec"
	"strings"
)

func ExecuteSystemCommand(command string) (string, error) {

	//Get the executable (e.g. From: cmd.exe /c dir, get: cmd.exe)
	executablePart := strings.Fields(command)[0]

	//Get the arguments (e.g. From: cmd.exe /c dir, get: {"/c", "dir", "\Users\"})
	argumentParts := strings.Fields(command)[1:]

	arguments := make([]string, 0)
	var quoteContent strings.Builder

	//Treat any content enclosed within "double quotes" as a single argument,
	//because it is probably a text/file/directory name containing spaces
	for _, part := range argumentParts {
		if strings.HasPrefix(part, "\"") && strings.HasSuffix(part, "\"") {
			//Case: "argument" - a single word is enclosed within double quotes
			arguments = append(arguments, strings.Trim(part, "\""))
		} else if strings.HasPrefix(part, "\"") {
			//Case: "start - double quotes open
			quoteContent.WriteString(strings.TrimPrefix(part, "\""))
		} else if strings.HasSuffix(part, "\"") {
			//Case: end" - double cotes close
			quoteContent.WriteString(" " + strings.TrimSuffix(part, "\""))
			arguments = append(arguments, quoteContent.String())
			//Reset the Builder for the next quoteContent
			quoteContent.Reset()
		} else if quoteContent.Len() > 0 {
			//Case: Concatenate words within double quotes
			quoteContent.WriteString(" " + part)
		} else {
			//Case: Regular argument
			arguments = append(arguments, part)
		}
	}

	//Executing system command
	output, err := exec.Command(executablePart, arguments...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
