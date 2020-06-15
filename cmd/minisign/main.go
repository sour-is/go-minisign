package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/sour-is/go-minisign"
	"github.com/sour-is/go-xdg"
)

func main() {

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Error:", r)
	// 		flag.Usage()
	// 		//			os.Exit(1)
	// 	}
	// }()

	args := flags()

	if err := run(args, os.Stdin, os.Stderr); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

var MinisignPath = xdg.EnvDirs(
	xdg.Env(minisign.SigDefaultConfigDirEnvVar),
	xdg.PrependDir(
		xdg.UserData,
		xdg.NewDirs(
			xdg.ParsePath(minisign.SigDefaultConfigDir),
			xdg.ParsePath("."),
			xdg.ParsePath("~"),
		)))

func run(args Args, _ io.ReadCloser, _ io.Writer) (err error) {
	switch args.Mode {
	case VerifyFiles:
		var pubkey *minisign.PublicKey
		var signature *minisign.Signature

		if args.PubKey != "" {
			if key, err := minisign.NewPublicKey(args.PubKey); err != nil {
				return errors.Wrap(err, "provided pubkey failed to parse")
			} else {
				pubkey = &key
			}
		}

		if pubkey == nil && args.PubKeyFile != "" {
			var pubkeyfile string
			if pubkeyfile, err = xdg.Find(MinisignPath, args.PubKeyFile); err != nil {
				return errors.Wrap(err, "provided pubkey file not found")
			}
			if key, err := minisign.NewPublicKeyFromFile(pubkeyfile); err != nil {
				return errors.Wrap(err, "provided pubkey file failed to parse")
			} else {
				pubkey = &key
			}
		}

		for _, f := range args.Files {
			fmt.Println("CHECK", f)
			sigFile := args.SigFile
			if sigFile == "" {
				sigFile = f + minisign.SigSuffix
			}

			if s, err := minisign.NewSignatureFromFile(sigFile); err != nil {
				return errors.Wrap(err, "provided sigfile file failed to parse")
			} else {
				signature = &s
			}

			ok, err := pubkey.VerifyFromFile(f, *signature)
			if err != nil {
				fmt.Println("RESULT", ok, err)
				os.Exit(1)
			}
			fmt.Println("RESULT", ok)
		}

	case PrintHelp:
		flag.Usage()
	default:
		return fmt.Errorf("operation not implemented")
	}

	return nil
}
