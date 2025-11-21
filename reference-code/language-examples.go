const (
	// Programming languages

	Bash       = "bash"
	C          = "c"
	CPP        = "cpp"
	CSharp     = "csharp"
	Go         = "go"
	Java       = "java"
	JavaScript = "javascript"
	JSON       = "json"
	Kotlin     = "kotlin"
	PHP        = "php"
	Python     = "python"
	Ruby       = "ruby"
	Rust       = "rust"
	Scala      = "scala"
	Shell      = "shell"
	Swift      = "swift"
	Text       = "text"
	TypeScript = "typescript"
	Undefined  = "undefined"
	XML        = "xml"
	YAML       = "yaml"

	// File extensions

	BashExtension       = ".sh"
	CExtension          = ".c"
	CPPExtension        = ".cpp"
	CSharpExtension     = ".cs"
	GoExtension         = ".go"
	JavaExtension       = ".java"
	JavaScriptExtension = ".js"
	JSONExtension       = ".json"
	KotlinExtension     = ".kt"
	PHPExtension        = ".php"
	PythonExtension     = ".py"
	RubyExtension       = ".rb"
	RustExtension       = ".rs"
	ScalaExtension      = ".scala"
	ShellExtension      = ".sh"
	SwiftExtension      = ".swift"
	TextExtension       = ".txt"
	TypeScriptExtension = ".ts"
	UndefinedExtension  = ".txt"
	XMLExtension        = ".xml"
	YAMLExtension       = ".yaml"
)

var CanonicalLanguages = []string{Bash, C, CPP,
	CSharp, Go, Java, JavaScript,
	JSON, Kotlin, PHP, Python,
	Ruby, Rust, Scala, Shell,
	Swift, Text, TypeScript, Undefined, XML, YAML,
}

func GetNormalizedLanguageFromString(language string) string {
	normalizeLanguagesMap := make(map[string]string)

	// Add the canonical languages and their values
	normalizeLanguagesMap[common.Bash] = common.Bash
	normalizeLanguagesMap[common.C] = common.C
	normalizeLanguagesMap[common.CPP] = common.CPP
	normalizeLanguagesMap[common.CSharp] = common.CSharp
	normalizeLanguagesMap[common.Go] = common.Go
	normalizeLanguagesMap[common.Java] = common.Java
	normalizeLanguagesMap[common.JavaScript] = common.JavaScript
	normalizeLanguagesMap[common.JSON] = common.JSON
	normalizeLanguagesMap[common.Kotlin] = common.Kotlin
	normalizeLanguagesMap[common.PHP] = common.PHP
	normalizeLanguagesMap[common.Python] = common.Python
	normalizeLanguagesMap[common.Ruby] = common.Ruby
	normalizeLanguagesMap[common.Rust] = common.Rust
	normalizeLanguagesMap[common.Scala] = common.Scala
	normalizeLanguagesMap[common.Shell] = common.Shell
	normalizeLanguagesMap[common.Swift] = common.Swift
	normalizeLanguagesMap[common.Text] = common.Text
	normalizeLanguagesMap[common.TypeScript] = common.TypeScript
	normalizeLanguagesMap[common.Undefined] = common.Undefined
	normalizeLanguagesMap[common.XML] = common.XML
	normalizeLanguagesMap[common.YAML] = common.YAML

	// Add variations and map to canonical values
	normalizeLanguagesMap[""] = common.Undefined
	normalizeLanguagesMap["console"] = common.Shell
	normalizeLanguagesMap["cs"] = common.CSharp
	normalizeLanguagesMap["golang"] = common.Go
	normalizeLanguagesMap["http"] = common.Text
	normalizeLanguagesMap["ini"] = common.Text
	normalizeLanguagesMap["js"] = common.JavaScript
	normalizeLanguagesMap["none"] = common.Undefined
	normalizeLanguagesMap["sh"] = common.Shell
	normalizeLanguagesMap["json\\n :copyable: false"] = common.JSON
	normalizeLanguagesMap["json\\n :copyable: true"] = common.JSON

	canonicalLanguage, exists := normalizeLanguagesMap[language]
	if exists {
		return canonicalLanguage
	} else {
		return common.Undefined
	}
}

func GetFileExtensionFromStringLang(language string) string {
	langExtensionMap := make(map[string]string)

	// Add the canonical languages and their extensions
	langExtensionMap[common.Bash] = common.BashExtension
	langExtensionMap[common.C] = common.CExtension
	langExtensionMap[common.CPP] = common.CPPExtension
	langExtensionMap[common.CSharp] = common.CSharpExtension
	langExtensionMap[common.Go] = common.GoExtension
	langExtensionMap[common.Java] = common.JavaExtension
	langExtensionMap[common.JavaScript] = common.JavaScriptExtension
	langExtensionMap[common.JSON] = common.JSONExtension
	langExtensionMap[common.Kotlin] = common.KotlinExtension
	langExtensionMap[common.PHP] = common.PHPExtension
	langExtensionMap[common.Python] = common.PythonExtension
	langExtensionMap[common.Ruby] = common.RubyExtension
	langExtensionMap[common.Rust] = common.RustExtension
	langExtensionMap[common.Scala] = common.ScalaExtension
	langExtensionMap[common.Shell] = common.ShellExtension
	langExtensionMap[common.Swift] = common.SwiftExtension
	langExtensionMap[common.Text] = common.TextExtension
	langExtensionMap[common.TypeScript] = common.TypeScriptExtension
	langExtensionMap[common.Undefined] = common.UndefinedExtension
	langExtensionMap[common.XML] = common.XMLExtension
	langExtensionMap[common.YAML] = common.YAMLExtension

	// Add variations and map to canonical values
	langExtensionMap[""] = common.UndefinedExtension
	langExtensionMap["console"] = common.ShellExtension
	langExtensionMap["cs"] = common.CSharpExtension
	langExtensionMap["golang"] = common.GoExtension
	langExtensionMap["http"] = common.TextExtension
	langExtensionMap["ini"] = common.TextExtension
	langExtensionMap["js"] = common.JavaScriptExtension
	langExtensionMap["none"] = common.UndefinedExtension
	langExtensionMap["sh"] = common.ShellExtension
	langExtensionMap["json\\n :copyable: false"] = common.JSONExtension
	langExtensionMap["json\\n :copyable: true"] = common.JSONExtension

	extension, exists := langExtensionMap[language]
	if exists {
		return extension
	} else {
		return common.UndefinedExtension
	}
}
