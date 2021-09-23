package translate

import (
	"bytes"
	"github.com/jmhodges/gocld3/cld3"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	CLDConfidenceThreshold = 0.3
	BergamotAppPath = "/root/bergamot-translator/build-native/app/"
	AllModelsPath = "/root/firefox-translations-models/models/"
	ConfigFolderPath = "/root/app/"
)

func DetectLanguage(text string) (string, error) {
	langId, cldErr := cld3.NewLanguageIdentifier(0, 512)
	if cldErr != nil {
		return "", cldErr
	}
	defer cld3.FreeLanguageIdentifier(langId)
	res := langId.FindLanguage(text)
	if res.Probability >= CLDConfidenceThreshold {
		return res.Language, nil
	} else {
		return "unknown", nil
	}
}

func TranslateTexts(texts []string, from string, to string) ([]string, error) {
	bergamotExecPath := BergamotAppPath + "bergamot"
	configFilePath := ConfigFolderPath + "config.yml"

	modelFolder := from + to + "/"
	modelDataPath := AllModelsPath + modelFolder

	if _, err := os.Stat(modelDataPath); os.IsNotExist(err) {
		if to == "en" {
			return []string{}, err
		}
		firstPhasePath := AllModelsPath + from + "en/"
		if _, err := os.Stat(firstPhasePath); os.IsNotExist(err) {
			return []string{}, err
		}
		secondPhasePath := AllModelsPath + "en" + to + "/"
		if _, err := os.Stat(secondPhasePath); os.IsNotExist(err) {
			return []string{}, err
		}
		firstPhaseResultTexts, err := TranslateTexts(texts, from, "en")
		if err != nil {
			return []string{}, err
		}
		return TranslateTexts(firstPhaseResultTexts, "en", to)
	}

	files, err := ioutil.ReadDir(modelDataPath)
	if err != nil {
		return []string{}, err
	}

	shortListFileName := ""
	vocabsFileName := ""
	modelsFileName := ""
	for _, file := range files {
		if !file.IsDir() {
			if strings.HasPrefix(file.Name(), "lex") {
				shortListFileName = modelDataPath + file.Name()
			}
			if strings.HasPrefix(file.Name(), "vocab") {
				vocabsFileName = modelDataPath + file.Name()
			}
			if strings.HasPrefix(file.Name(), "model") {
				modelsFileName = modelDataPath + file.Name()
			}
		}
	}

	if (shortListFileName == "") || (vocabsFileName == "") || (modelsFileName == "") {
		return []string{}, err
	}

	originalText := ""
	for ind, text := range texts {
		if ind > 0 {
			originalText = originalText + "\n"
		}
		originalText = originalText + text
	}

	bergamotProcess := exec.Command(bergamotExecPath,
									"--config", configFilePath,
									"--bergamot-mode", "native",
									"--models", modelsFileName,
									"--vocabs", vocabsFileName,
									"--vocabs", vocabsFileName,
									"--shortlist", shortListFileName,
									"--shortlist", "false")

	var outb, errb bytes.Buffer
	bergamotProcess.Stdout = &outb
	bergamotProcess.Stderr = &errb

	bergamotStdin, err := bergamotProcess.StdinPipe()
	if err != nil {
		return []string{}, err
	}

	err = bergamotProcess.Start()
	if err != nil {
		return []string{}, err
	}

	io.WriteString(bergamotStdin, originalText)
	bergamotStdin.Close()
	bergamotProcess.Wait()

	translatedText := outb.String()
	translations := strings.Split(translatedText, "\n")

	return translations, nil
}
