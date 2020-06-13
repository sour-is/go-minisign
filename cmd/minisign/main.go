package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adrg/xdg"

	"github.com/pkg/errors"
	"github.com/sour-is/go-minisign"
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
			os.Setenv("XDG_CONFIG_DIRS",
				joinPaths(
					".",
					os.Getenv("SigDefaultConfigDirEnvVar"),
					minisign.SigDefaultConfigDir,
					os.Getenv("XDG_CONFIG_DIRS"),
				),
			)
			xdg.Reload()

			var pubkeyfile string
			if pubkeyfile, err = xdg.SearchConfigFile(args.PubKeyFile); err != nil {
				return errors.Wrap(err, "provided pubkey file not found")
			}

			if *pubkey, err = minisign.NewPublicKeyFromFile(pubkeyfile); err != nil {
				return errors.Wrap(err, "provided pubkey file failed to parse")
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
			fmt.Println(ok, err)
		}

	case PrintHelp:
		flag.Usage()
	default:
		return fmt.Errorf("operation not implemented")
	}

	return nil
}

func joinPaths(path ...string) string {
	var arr = make([]string, len(path))

	for _, s := range path {
		if s != "" {
			arr = append(arr, s)
		}
	}

	return strings.Join(path, ":")
}
