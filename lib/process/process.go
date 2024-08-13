package process

import (
  "strings"
  "slices"
  "github.com/h2so5/goback/regexp"
)

func ProcessHTML(line string) string{
  // Remove all HTML tags
  chars := strings.Split(line, "")
  var begin int
  var end int
  for {
    if !slices.Contains(chars, "<") && !slices.Contains(chars, ">"){
      break
    }
    for i := 0; i < len(chars); i++{
      if chars[i] == "<"{
        begin = i
      }
      if chars[i] == ">"{
        end = i
        chars = slices.Delete(chars, begin, end + 1)
        break
      }
    }
  }
  return strings.Join(chars, "")
}

func ReplaceBlockCharsR(line string) string {
  blockChars := regexp.MustCompile("&#xF0(?!20;).{3}")
  result := blockChars.ReplaceAll([]byte(line), []byte("#"))
  blackBlockChars := regexp.MustCompile("&#xF020;")
  result = blackBlockChars.ReplaceAll(result, []byte(" "))
  return string(result)
}

func ReplaceSpecialChars(line string) string{
  // https://www.html.am/reference/html-special-characters.cfm ISO 8859-1 section
  conversionMap := map[string]string {
    "&Agrave;": "À",
    "&Aacute;": "Á",
    "&Acirc;": "Â",
    "&Atilde;": "Ã",
    "&Auml;": "Ä",
    "&Aring;": "Å",
    "&AElig;": "Æ",
    "&Ccedil;": "Ç",
    "&Egrave;": "È",
    "&Eacute;": "É",
    "&Ecirc;": "Ê",
    "&Euml;": "Ë",
    "&Igrave;": "Ì",
    "&Iacute;": "Í",
    "&Icirc;": "Î",
    "&Iuml;": "Ï",
    "&ETH;": "Ð",
    "&Ntilde;": "Ñ",
    "&Ograve;": "Ò",
    "&Oacute;": "Ó",
    "&Ocirc;": "Ô",
    "&Otilde;": "Õ",
    "&Ouml;": "Ö",
    "&Oslash;": "Ø",
    "&Ugrave;": "Ù",
    "&Uacute;": "Ú",
    "&Ucirc;": "Û",
    "&Uuml;": "Ü",
    "&Yacute;": "Ý",
    "&THORN;": "Þ",
    "&szlig;": "ß",
    "&agrave;": "à",
    "&aacute;": "á",
    "&acirc;": "â",
    "&atilde;": "ã",
    "&auml;": "ä",
    "&aring;": "å",
    "&aelig;": "æ",
    "&ccedil;": "ç",
    "&egrave;": "è",
    "&eacute;": "é",
    "&ecirc;": "ê",
    "&euml;": "ë",
    "&igrave;": "ì",
    "&iacute;": "í",
    "&icirc;": "î",
    "&iuml;": "ï",
    "&eth;": "ð",
    "&ntilde;": "ñ",
    "&ograve;": "ò",
    "&oacute;": "ó",
    "&ocirc;": "ô",
    "&otilde;": "õ",
    "&ouml;": "ö",
    "&oslash;": "ø",
    "&ugrave;": "ù",
    "&uacute;": "ú",
    "&ucirc;": "û",
    "&uuml;": "ü",
    "&yacute;": "ý",
    "&thorn;": "þ",
    "&yuml;": "ÿ",
  }

  for k, v := range conversionMap {
    if strings.Contains(line, k){
      line = strings.ReplaceAll(line, k, v) 
    }
  }
  return line
}
