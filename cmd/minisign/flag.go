package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sour-is/go-minisign"
)

type Mode int

const (
	GenerateKey Mode = iota
	SignFiles
	VerifyFiles
	CreatePubKey
	PrintVersion
	PrintHelp
)

type Args struct {
	Mode             Mode
	Files            []string
	PubKeyFile       string
	SecKeyFile       string
	PubKey           string
	SigFile          string
	UntrustedComment string
	TrustedComment   string
	OutputContent    bool
	PreHash          bool
	Quiet            bool
	Pretty           bool
	Force            bool
}

func flags() Args {
	var args Args
	var (
		generateKey    bool
		generatePubKey bool
		signFiles      bool
		verifyFiles    bool
		printVersion   bool

		files listString
	)

	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.Usage = func() {
		fmt.Print(`Usage:
minisign -G [-p pubkey] [-s seckey]
minisign -S [-H] [-x sigfile] [-s seckey] [-c untrusted_comment] [-t trusted_comment] -m file [file ...]
minisign -V [-x sigfile] [-p pubkeyfile | -P pubkey] [-o] [-q] -m file
minisign -R -s seckey -p pubkeyfile

`)
		flag.PrintDefaults()
	}

	flag.BoolVar(&generateKey, "G", false, "generate a new key pair")
	flag.BoolVar(&signFiles, "S", false, "sign files")
	flag.BoolVar(&verifyFiles, "V", false, "verify that a signature is valid for a given file")
	flag.Var(&files, "m", "file to sign/verify")
	flag.BoolVar(&args.OutputContent, "o", false, "combined with -V, output the file content after verification")
	flag.BoolVar(&args.PreHash, "H", false, "combined with -S, pre-hash in order to sign large files")
	flag.StringVar(&args.PubKeyFile, "p", minisign.SigDefaultPKFile, "public key file (default: ./minisign.pub)")
	flag.StringVar(&args.PubKey, "P", "", "public key, as a base64 string")
	flag.StringVar(&args.SecKeyFile, "s", minisign.SigDefaultSKFile, "secret key file (default: ~/.minisign/minisign.key)")
	flag.StringVar(&args.SigFile, "x", "", "signature file (default: <file>.minisig)")
	flag.StringVar(&args.UntrustedComment, "c", "", "add a one-line untrusted comment")
	flag.StringVar(&args.TrustedComment, "t", "", "add a one-line trusted comment")
	flag.BoolVar(&args.Quiet, "q", false, "quiet mode, suppress output")
	flag.BoolVar(&args.Pretty, "Q", false, "pretty quiet mode, only print the trusted comment")
	flag.BoolVar(&generatePubKey, "R", false, "recreate a public key file from a secret key file")
	flag.BoolVar(&args.Force, "f", false, "force. Combined with -G, overwrite a previous key pair")
	flag.BoolVar(&printVersion, "v", false, "display version number")

	err := flag.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}

	switch {
	case generateKey:
		args.Mode = GenerateKey
	case signFiles:
		args.Mode = SignFiles
	case verifyFiles:
		args.Mode = VerifyFiles
	case generatePubKey:
		args.Mode = CreatePubKey
	case printVersion:
		args.Mode = PrintVersion
	default:
		flag.Usage()
		os.Exit(0)
	}

	args.Files = files

	return args
}

// ListString from multiple calls to flag
type listString []string

func (i *listString) String() string {
	if i == nil {
		return ""
	}

	return strings.Join(*i, ",")
}

func (i *listString) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}
