package qbc

import (
	"github.com/rskvp/qb-core/qb_coding/openssl"
	"github.com/rskvp/qb-core/qb_email"
	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_exec"
	"github.com/rskvp/qb-core/qb_exec_bucket"
	"github.com/rskvp/qb-core/qb_fnvars"
	"github.com/rskvp/qb-core/qb_generators/genusers"
	"github.com/rskvp/qb-core/qb_i18n_bundle"
	"github.com/rskvp/qb-core/qb_license"
	"github.com/rskvp/qb-core/qb_log"
	"github.com/rskvp/qb-core/qb_net"
	"github.com/rskvp/qb-core/qb_nio"
	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_shamir"
	"github.com/rskvp/qb-core/qb_state"
	"github.com/rskvp/qb-core/qb_stegano"
	"github.com/rskvp/qb-core/qb_stopwatch"
	"github.com/rskvp/qb-core/qb_structs"
	"github.com/rskvp/qb-core/qb_sys"
	"github.com/rskvp/qb-core/qb_utils"
	"github.com/rskvp/qb-core/qb_utils/phonenumbers"
	"github.com/rskvp/qb-core/qb_vcal"
	"github.com/rskvp/qb-core/qb_watchdog"
	"github.com/rskvp/qb-core/qb_xtend"
)

var Strings *qb_utils.StringsHelper
var Arrays *qb_utils.ArraysHelper
var Convert *qb_utils.ConversionHelper
var Compare *qb_utils.CompareHelper
var Dates *qb_utils.DatesHelper
var Errors *qb_utils.ErrorsHelper
var Regex *qb_utils.RegexHelper
var Rnd *qb_rnd.RndHelper
var Paths *qb_utils.PathsHelper
var IO *qb_utils.IoHelper
var Zip *qb_utils.ZipHelper
var Maps *qb_utils.MapsHelper
var JSON *qb_utils.JSONHelper
var XML *qb_utils.XMLHelper
var Reflect *qb_utils.ReflectHelper
var CSV *qb_utils.CsvHelper
var Async *qb_utils.AsyncHelper
var Coding *qb_utils.CodingHelper
var BOM *qb_utils.BOMHelper
var Sys *qb_sys.SysHelper
var Formatter *qb_utils.FormatterHelper
var MIME *qb_utils.MIMEHelper
var Exec *qb_exec.ExecHelper
var StopWatch *qb_stopwatch.StopWatchHelper
var FnVars *qb_fnvars.FnVarsHelper
var GenUsers *genusers.GenUsersHelper
var PhoneNumber *phonenumbers.PhoneNumberHelper
var VCal *qb_vcal.VCalHelper
var Structs *qb_structs.StructHelper
var Stegano *qb_stegano.SteganoHelper
var Shamir *qb_shamir.ShamirHelper
var State *qb_state.StateHelper
var OpenSSL *openssl.OpenSSLHelper
var Log *qb_log.LogHelper
var Bucket *qb_exec_bucket.BucketHelper
var Dir *qb_utils.DirHelper
var License *qb_license.LicenseHelper
var Net *qb_net.NetHelper
var Xtend *qb_xtend.XtendHelper

//-- advanced --//

var I18N *qb_i18n_bundle.I18NHelper
var Events *qb_events.EventsHelper
var NIO *qb_nio.NIOHelper
var Email *qb_email.EmailHelper
var Watchdog *qb_watchdog.WatchdogHelper

func init() {
	Strings = qb_utils.Strings
	Arrays = qb_utils.Arrays
	Convert = qb_utils.Convert
	Compare = qb_utils.Compare
	Dates = qb_utils.Dates
	Errors = qb_utils.Errors
	Regex = qb_utils.Regex
	Rnd = qb_rnd.Rnd
	Paths = qb_utils.Paths
	IO = qb_utils.IO
	Zip = qb_utils.Zip
	Maps = qb_utils.Maps
	JSON = qb_utils.JSON
	XML = qb_utils.XML
	Reflect = qb_utils.Reflect
	CSV = qb_utils.CSV
	Async = qb_utils.Async
	Coding = qb_utils.Coding
	BOM = qb_utils.BOM
	Sys = qb_sys.Sys
	Formatter = qb_utils.Formatter
	MIME = qb_utils.MIME
	Exec = qb_exec.Exec
	StopWatch = qb_stopwatch.Watch
	FnVars = qb_fnvars.FnVars
	GenUsers = genusers.GenUsers
	VCal = qb_vcal.VCal
	Structs = qb_structs.Structs
	Stegano = qb_stegano.Stegano
	Shamir = qb_shamir.Shamir
	OpenSSL = openssl.OpenSSLUtil
	State = qb_state.StateH
	Log = qb_log.Log
	Bucket = qb_exec_bucket.Bucket
	Dir = qb_utils.Dir
	PhoneNumber = phonenumbers.PhoneNumber
	License = qb_license.License
	Net = qb_net.Net
	Xtend = qb_xtend.Xtend

	// advanced
	I18N = qb_i18n_bundle.I18N
	Events = qb_events.Events
	NIO = qb_nio.NIO
	Email = qb_email.Email
	Watchdog = qb_watchdog.Watchdog
}
