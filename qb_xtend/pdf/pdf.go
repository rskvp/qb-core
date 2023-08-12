package pdf

type IPdfExtension interface {
	PDF2TextFile(source string) (outputFile string, err error)
	PDF2Text(inputPath string) (string, error)
	PDFProtect(inputPath string, outputPath string, userPassword, ownerPassword string) error
}
