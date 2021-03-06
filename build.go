package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/thijzert/go-resemble"
)

type job func(ctx context.Context) error

func main() {
	if _, err := os.Stat("qm9k/assets"); err != nil {
		log.Fatalf("Error: cannot find quizmeister9000 assets directory. (error: %s)\nAre you running this from the repository root?", err)
	}

	devBuild := false
	quickBuild := false
	watch := false
	run := false
	flag.BoolVar(&devBuild, "development", false, "Create a development build")
	flag.BoolVar(&quickBuild, "quick", false, "Create a development build")
	flag.BoolVar(&watch, "watch", false, "Watch source tree for changes")
	flag.BoolVar(&run, "run", false, "Run quizmeister9000 upon successful compilation")
	flag.Parse()

	if devBuild && quickBuild {
		log.Printf("")
		log.Printf("You requested a quick build. This will assume")
		log.Printf(" you have a version of  `gulp watch`  running")
		log.Printf(" in a separate process.")
		log.Printf("")
	}

	var theJob job

	if run {
		theJob = func(ctx context.Context) error {
			err := compile(ctx, devBuild, quickBuild)
			if err != nil {
				return err
			}
			runArgs := append([]string{"./quizmeister9000"}, flag.Args()...)
			return passthru(ctx, runArgs...)
		}
	} else {
		theJob = func(ctx context.Context) error {
			return compile(ctx, devBuild, quickBuild)
		}
	}

	if watch {
		theJob = watchSourceTree([]string{"."}, []string{"*.go"}, theJob)
	}

	err := theJob(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func compile(ctx context.Context, devBuild, quickBuild bool) error {
	// Check for local gulp
	fi, err := os.Stat("node_modules/.bin/gulp")
	if (!devBuild && !quickBuild) || err != nil || fi.Mode()&0x1 != 1 {
		err = passthru(ctx, "npm", "install")
		if err != nil {
			return errors.WithMessage(err, "error installing local gulp")
		}
	}

	// Compile static assets
	production := "--production"
	if devBuild {
		production = "--development"
	}
	if !devBuild || !quickBuild {
		err = passthru(ctx, "node_modules/.bin/gulp", "compile", production)
		if err != nil {
			return errors.WithMessage(err, "error compiling assets")
		}
	}

	// Embed static assets
	if err := os.Chdir("qm9k/assets"); err != nil {
		return errors.Errorf("Error: cannot find quizmeister9000 assets directory. (error: %s)\nAre you *sure* you're running this from the repository root?", err)
	}
	var emb resemble.Resemble
	emb.OutputFile = "../assets.go"
	emb.PackageName = "qm9k"
	emb.Debug = devBuild
	emb.AssetPaths = []string{
		".",
	}
	if err := emb.Run(); err != nil {
		return errors.WithMessage(err, "error running 'resemble'")
	}

	os.Chdir("../..")

	// Build main executable
	gofiles, err := filepath.Glob("cmd/quizmeister9000/*.go")
	if err != nil || gofiles == nil {
		return errors.WithMessage(err, "error: cannot find any go files to compile.")
	}
	compileArgs := append([]string{
		"go", "build", "-o", "quizmeister9000",
	}, gofiles...)
	err = passthru(ctx, compileArgs...)
	if err != nil {
		return errors.WithMessage(err, "compilation failed")
	}

	if devBuild && !quickBuild {
		log.Printf("")
		log.Printf("Development build finished. For best results,")
		log.Printf(" run  `node_modules/.bin/gulp watch`  in a")
		log.Printf(" separate process.")
		log.Printf("")
	} else {
		log.Printf("Compilation finished.")
	}

	return nil
}

func passthru(ctx context.Context, argv ...string) error {
	c := exec.CommandContext(ctx, argv[0], argv[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func watchSourceTree(paths []string, fileFilter []string, childJob job) job {
	return func(ctx context.Context) error {
		var mu sync.Mutex
		for {
			lastHash := sourceTreeHash(paths, fileFilter)
			current := lastHash
			cctx, cancel := context.WithCancel(ctx)
			go func() {
				mu.Lock()
				err := childJob(cctx)
				if err != nil {
					log.Printf("child process: %s", err)
				}
				mu.Unlock()
			}()

			for lastHash == current {
				time.Sleep(250 * time.Millisecond)
				current = sourceTreeHash(paths, fileFilter)
			}

			log.Printf("Source change detected - rebuilding")
			cancel()
		}
	}
}

func sourceTreeHash(paths []string, fileFilter []string) string {
	h := sha1.New()
	for _, d := range paths {
		h.Write(directoryHash(0, d, fileFilter))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func directoryHash(level int, filePath string, fileFilter []string) []byte {
	h := sha1.New()
	h.Write([]byte(filePath))

	fi, err := os.Stat(filePath)
	if err != nil {
		return h.Sum(nil)
	}
	if fi.IsDir() {
		base := filepath.Base(filePath)
		if level > 0 {
			if base == ".git" || base == ".." || base == "node_modules" {
				return []byte{}
			}
		}
		// recurse
		var names []string
		f, err := os.Open(filePath)
		if err == nil {
			names, err = f.Readdirnames(-1)
		}
		if err == nil {
			for _, name := range names {
				if name == "" || name[0] == '.' {
					continue
				}
				h.Write(directoryHash(level+1, path.Join(filePath, name), fileFilter))
			}
		}
	} else {
		if fileFilter != nil {
			found := false
			for _, pattern := range fileFilter {
				if ok, _ := filepath.Match(pattern, filePath); ok {
					found = true
				} else if ok, _ := filepath.Match(pattern, filepath.Base(filePath)); ok {
					found = true
				}
			}
			if !found {
				return []byte{}
			}
		}
		f, err := os.Open(filePath)
		if err == nil {
			io.Copy(h, f)
			f.Close()
		}
	}
	return h.Sum(nil)
}
