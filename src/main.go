package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	TEMPORARY_DIR      = "./temp"
	TOKEN_REGEX        = `$token_prefix([a-zA-Z0-9\-_\$]*)$token_suffix`
	ENV_VARIABLE_REGEX = `^\$([a-zA-Z0-9\-_]+)$`
	FILE_ERROR_ACCESS  = `file '%s' not found or without access -> `
)

var separator = string(os.PathSeparator)

func main() {
	filename := ""
	if len(os.Args) < 5 {
		if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
			printHelp()
			return
		} else if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
			printVersion()
			return
		} else {
			log.Fatalln("Insufficient arguments")
		}
	}

	tokenPrefix := string(os.Args[1])
	tokenSuffix := string(os.Args[2])
	tokensFilePath := string(os.Args[3])
	toReplaceFilePath := string(os.Args[4])
	if len(os.Args) > 5 {
		filename = os.Args[5]
	}

	log.Printf("replacing tokens on file %s", toReplaceFilePath)

	_, tokensExtension := getFileName(tokensFilePath)
	if tokensExtension != "yaml" && tokensExtension != "yml" {
		log.Fatalf(
			"Invalid file tokens, must be a yaml file (.yaml - .yml) and must be on yaml format -> file extension: %s",
			tokensExtension,
		)
	}

	tokenPrefixEscaped := escapeRegexChars(tokenPrefix)
	tokenSuffixEscaped := escapeRegexChars(tokenSuffix)

	tokenRegexS := strings.ReplaceAll(TOKEN_REGEX, "$token_prefix", tokenPrefixEscaped)
	tokenRegexS = strings.ReplaceAll(tokenRegexS, "$token_suffix", tokenSuffixEscaped)

	tokenRegex, err := regexp.Compile(tokenRegexS)
	if err != nil {
		log.Fatalln("Invalid prefix or suffix -> ", tokenPrefix, tokenSuffix)
	}
	envRegex := regexp.MustCompile(ENV_VARIABLE_REGEX)

	replaceTokens(tokenPrefix, tokenSuffix, tokensFilePath, tokensFilePath, "", tokenRegex, envRegex, true)
	replacedTokens, extension := getFileName(tokensFilePath)
	if len(extension) > 0 {
		extension = "." + extension
	}
	replacedTokens = replacedTokens + "-replaced" + extension
	replaceTokens(tokenPrefix, tokenSuffix, replacedTokens, toReplaceFilePath, filename, tokenRegex, envRegex, true)

}

func replaceTokens(prefix, suffix, tokensPath, toReplacePath, output string,
	haveTokensRegex, envRegex *regexp.Regexp, useFlag bool) {

	fileLines, err := getLines(toReplacePath)
	if err != nil {
		log.Fatal(fmt.Sprintf(FILE_ERROR_ACCESS, toReplacePath), err)
	}
	tokenLines, err := getLines(tokensPath)
	if err != nil {
		log.Fatal(fmt.Sprintf(FILE_ERROR_ACCESS, tokensPath), err)
	}

	replaced := []string{}
	var replacedLines int = 1

	escapedTokenPreffix := escapeRegexChars(prefix)
	prefixRegex := regexp.MustCompile(escapedTokenPreffix)

	for replacedLines != 0 {
		replacedLines = 0
		for _, l := range fileLines {
			if s := haveTokensRegex.FindStringSubmatch(l); s == nil {
				replaced = append(replaced, l)
				continue
			}
			value := strings.Split(strings.Split(l, prefix)[1], suffix)[0]

			splittedSuff := strings.Split(l, suffix)[1:]
			nextSplit := []string{}
			for _, w := range splittedSuff {
				if s1 := prefixRegex.FindStringSubmatch(w); s1 != nil {
					nextSplit = append(nextSplit, fmt.Sprintf("%s%s", w, suffix))
					continue
				}
				nextSplit = append(nextSplit, w)
			}

			next := strings.Join(nextSplit, "")

			name := strings.Split(l, prefix)[0]
			if envRegex.MatchString(value) {
				replaced = getValueFromEnv(name, value, next, replaced)
				replacedLines++
			} else {
				replacedLines, replaced = getValueFromTokensFile(name, value, next, replaced, tokenLines)
			}
		}
		fileLines = replaced
		replaced = []string{}
	}

	out, err := writeLines(toReplacePath, output, useFlag, fileLines)
	if err != nil {
		log.Fatal(fmt.Sprintf(FILE_ERROR_ACCESS, toReplacePath), err)
	}
	log.Printf("wrote %d lines in output file -> %s", len(fileLines), out)
}

func getValueFromEnv(name, value, next string, replaced []string) []string {
	value = strings.ReplaceAll(value, "$", "")
	varValue := os.Getenv(value)
	log.Printf("replacing :::: %s%s -> %s", name, value, varValue)
	replaced = append(replaced, fmt.Sprintf("%s%s%s", name, varValue, next))
	return replaced
}

func getValueFromTokensFile(name, value, next string, replaced, tokenLines []string) (int, []string) {
	replacedLines := 0
	for _, t := range tokenLines {
		if value != strings.Split(t, ":")[0] {
			continue
		}
		varValue := strings.TrimSpace(strings.Split(t, ":")[1])
		log.Printf("replacing :::: %s%s -> %s", name, value, varValue)
		replaced = append(replaced, fmt.Sprintf("%s%s%s", name, varValue, next))
		replacedLines++
	}
	return replacedLines, replaced
}

func writeLines(path, output string, useFlag bool, fileLines []string) (string, error) {
	flag := "-replaced"
	if !useFlag {
		flag = ""
	}

	filepath, extension := getFileName(path)
	filepath = filepath + flag
	if len(extension) > 0 {
		filepath = filepath + "." + extension
	}

	if len(output) > 0 {
		paths := strings.Split(filepath, separator)
		if len(paths) > 1 {
			filepath = strings.Join(append(paths[:len(paths)-2], output), separator)
		} else {
			filepath = output
		}
	}

	if _, err := os.Stat(filepath); err == nil {
		if err2 := os.Remove(filepath); err2 != nil {
			return filepath, err2
		}
	}

	f, err := os.Create(filepath)
	if err != nil {
		return filepath, err
	}

	for _, v := range fileLines {
		if _, err := f.Write([]byte(v + "\n")); err != nil {
			return filepath, err
		}
	}

	return filepath, f.Close()
}

func getLines(filePath string) ([]string, error) {
	lines := []string{}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func escapeRegexChars(toEscape string) string {
	chars := []rune{}
	especialChars := `[\\\^\$\.\|\?\*\+\(\)\[\]\{\}]`
	especialRegex := regexp.MustCompile(especialChars)

	for _, v := range toEscape {
		if especialRegex.Match([]byte(string(v))) {
			chars = append(chars, []rune("\\")...)
		}
		chars = append(chars, v)
	}

	return string(chars)
}

func getFileName(path string) (string, string) {
	spplited := strings.Split(path, ".")

	if len(spplited) < 2 {
		return spplited[0], ""
	}
	if len(spplited) > 2 {
		return spplited[len(spplited)-2], spplited[len(spplited)-1]
	}

	return spplited[0], spplited[1]
}

func printHelp() {
	helpPrompt := `Usage:
	replacetokens [PREFIX] [SUFFIX] [TOKENS FILE] [TO REPLACE] [[OUTPUT]]

PREFIX		required	Prefix used to denote a token in the file to replace.
SUFFIX		required	Suffix used to denote a token in the file to replace.
TOKENS FILE	required	Path to the file that contains the values of the tokens to be replaced. 
					Must be key-value in YAML format with just one hierarchical level. Eg. **TOKEN: value**.
TO REPLACE	required	Path to the file that contains the tokens which might be replaced by the values in the TOKENS FILE 
OUTPUT		optional	Path to the file where will be wrote the TO REPLACE file with all tokens replaced. 
					If OUTPUT is not provided, the output file will be paced in the same location oh the TO REPLACE 
					file with the same name adding a flag at the end of the filename with value of "-replaced"`
	fmt.Println(helpPrompt)
}

func printVersion() {
	versionPrompt := `replacetokens v1.0.1
Build by github.com/enaldo1709`

	fmt.Println(versionPrompt)
}
