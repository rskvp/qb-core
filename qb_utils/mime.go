package qb_utils

import (
	"strings"
)

type MIMEHelper struct {
}

var MIME *MIMEHelper
var csvData []map[string]string

func init() {
	MIME = new(MIMEHelper)

	options := CSV.NewCsvOptionsDefaults()
	options.FirstRowHeader = true
	options.Comma = ","
	data, err := CSV.ReadAll(mimeCsv, options)
	if nil == err {
		csvData = data
	}
}

const (
	fldName        = "Name"
	fldMIMEType    = "MIMEType"
	fldExtension   = "Extension"
	fldDescription = "Description"
)

const mimeCsv = `Name,MIMEType,Extension,Description
aac,audio/aac,.aac,AAC audio file
abw,application/x-abiword,.abw,AbiWord document
arc,application/octet-stream,.arc,Archive document (multiple files embedded)
avi,video/x-msvideo,.avi,AVI: Audio Video Interleave
azw,application/vnd.amazon.ebook,.azw,Amazon Kindle eBook format
bin,application/octet-stream,.bin,Any kind of binary data
bz,application/x-bzip,.bz,BZip archive
bz2,application/x-bzip2,.bz2,BZip2 archive
csh,application/x-csh,.csh,C-Shell script
css,text/css,.css,Cascading Style Sheets (CSS)
csv,text/csv,.csv,Comma-separated values (CSV)
doc,application/msword,.doc,Microsoft Word
epub,application/epub+zip,.epub,Electronic publication (EPUB)
gif,image/gif,.gif,Graphics Interchange Format (GIF)
htm,text/html,.htm,
html,text/html,.html,HyperText Markup Language (HTML)
ico,image/x-icon,.ico,Icon format
ics,text/calendar,.ics,iCalendar format
jar,application/java-archive,.jar,Java Archive (JAR)
jpeg,image/jpeg,.jpeg,
jpg,image/jpeg,.jpg,JPEG images
js,application/javascript,.js,JavaScript (ECMAScript)
json,application/json,.json,JSON format
mid,audio/midi,.mid,
midi,audio/midi,.midi,Musical Instrument Digital Interface (MIDI)
mpeg,video/mpeg,.mpeg,MPEG Video
mpkg,application/vnd.apple.installer+xml,.mpkg,Apple Installer Package
odp,application/vnd.oasis.opendocument.presentation,.odp,OpenDocuemnt presentation document
ods,application/vnd.oasis.opendocument.spreadsheet,.ods,OpenDocuemnt spreadsheet document
odt,application/vnd.oasis.opendocument.text,.odt,OpenDocument text document
oga,audio/ogg,.oga,OGG audio
ogv,video/ogg,.ogv,OGG video
ogx,application/ogg,.ogx,OGG
pdf,application/pdf,.pdf,Adobe Portable Document Format (PDF)
png,image/png,.png,PNG images
ppt,application/vnd.ms-powerpoint,.ppt,Microsoft PowerPoint
rar,application/x-rar-compressed,.rar,RAR archive
rtf,application/rtf,.rtf,Rich Text Format (RTF)
sh,application/x-sh,.sh,Bourne shell script
svg,image/svg+xml,.svg,Scalable Vector Graphics (SVG)
swf,application/x-shockwave-flash,.swf,Small web format (SWF) or Adobe Flash document
tar,application/x-tar,.tar,Tape Archive (TAR)
tif,image/tiff,.tif,
tiff,image/tiff,.tiff,Tagged Image File Format (TIFF)
ttf,font/ttf,.ttf,TrueType Font
vsd,application/vnd.visio,.vsd,Microsft Visio
wav,audio/x-wav,.wav,Waveform Audio Format
weba,audio/webm,.weba,WEBM audio
webm,video/webm,.webm,WEBM video
webp,image/webp,.webp,WEBP image
woff,font/woff,.woff,Web Open Font Format (WOFF)
woff2,font/woff2,.woff2,Web Open Font Format (WOFF)
xhtml,application/xhtml+xml,.xhtml,XHTML
xls,application/vnd.ms-excel,.xls,Microsoft Excel
xml,application/xml,.xml,XML
xul,application/vnd.mozilla.xul+xml,.xul,XUL
zip,application/zip,.zip,ZIP archive
3gp,video/3gpp,.3gp,3GPP audio/video container
3g2,video/3gpp2,.3g2,3GPP2 audio/video container
7z,application/x-7z-compressed,.7z,7-zip archive
docx,application/msword,.docx,Microsoft Word
txt,text/plain,.txt,Simple text`

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *MIMEHelper) GetMimeTypeByExtension(ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ext = strings.ToLower(ext)
	for _, m := range csvData {
		if ext == m[fldExtension] {
			return m[fldMIMEType]
		}
	}
	return "application/octet-stream"
}

