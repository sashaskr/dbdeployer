// Copyright © 2011-12 Qtrac Ltd.
//
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
/*
 Code adapted and enhanced from examples to the book:
 Programming in Go by Mark Summerfield
 http://www.qtrac.eu/gobook.html

 Original author: Mark Summerfield
 Converted to package by Giuseppe Maxia in 2018

 The original code was a stand-alone program, and it
 had a few bugs:
 * when extracting from a tar file: when there
 isn't a separate item for each directory, the
 extraction fails.
 * The attributes of the files were not reproduced
 in the extracted files.
 This code fixes those problems and introduces a
 destination directory and verbosity
 levels for the extraction

*/

package unpack

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	SILENT  = iota // No output
	VERBOSE        // Minimal feedback about extraction operations
	CHATTY         // Full details of what is being extracted
)

var Verbose int

func cond_print(s string, nl bool, level int) {
	if Verbose >= level {
		if nl {
			fmt.Println(s)
		} else {
			fmt.Printf(s)
		}
	}
}

func validSuffix(filename string) bool {
	for _, suffix := range []string{".tgz", ".tar", ".tar.gz"} {
		if strings.HasSuffix(filename, suffix) {
			return true
		}
	}
	return false
}

func UnpackTar(filename string, destination string, verbosity_level int) (err error) {
	Verbose = verbosity_level
	f, err := os.Stat(destination)
	if os.IsNotExist(err) {
		return fmt.Errorf("Destination directory '%s' does not exist", destination)
	}
	filemode := f.Mode()
	if filemode.IsDir() == false {
		return fmt.Errorf("Destination '%s' is not a directory", destination)
	}
	if !validSuffix(filename) {
		return fmt.Errorf("unrecognized archive suffix")
	}
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		return err
	}
	defer file.Close()
	os.Chdir(destination)
	var fileReader io.Reader = file
	var decompressor *gzip.Reader
	if strings.HasSuffix(filename, ".gz") {
		if decompressor, err = gzip.NewReader(file); err != nil {
			return err
		}
		defer decompressor.Close()
	}
	var reader *tar.Reader
	if decompressor != nil {
		reader = tar.NewReader(decompressor)
	} else {
		reader = tar.NewReader(fileReader)
	}
	return unpackTarFiles(reader)
}

func unpackTarFiles(reader *tar.Reader) (err error) {
	var header *tar.Header
	var count int = 0

	for {
		if header, err = reader.Next(); err != nil {
			if err == io.EOF {
				cond_print("Files ", false, CHATTY)
				cond_print(strconv.Itoa(count), true, 1)
				return nil // OK
			}
			return err
		}
		//fmt.Printf("%#v\n", header)
		filemode := os.FileMode(header.Mode)
		filename := sanitizedName(header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filename, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			fileDir := path.Dir(filename)
			if _, err := os.Stat(fileDir); os.IsNotExist(err) {
				if err = os.MkdirAll(fileDir, 0755); err != nil {
					return err
				}
				cond_print(" + "+fileDir+" ", true, CHATTY)
			}
			if err = unpackTarFile(filename, header.Name, reader); err != nil {
				return err
			}
			os.Chmod(filename, filemode)
			count++
			cond_print(filename, true, CHATTY)
			if count%10 == 0 {
				mark := "."
				if count%100 == 0 {
					mark = strconv.Itoa(count)
				}
				if Verbose < CHATTY {
					cond_print(mark, false, 1)
				}
			}
		}
	}
	return nil
}

func unpackTarFile(filename, tarFilename string,
	reader *tar.Reader) (err error) {
	var writer *os.File
	if writer, err = os.Create(filename); err != nil {
		return err
	}
	defer writer.Close()
	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	return nil
}

func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' {
		filename = filename[2:]
	}
	filename = strings.TrimLeft(filename, "\\/.")
	filename = strings.Replace(filename, "../", "", -1)
	return strings.Replace(filename, "..\\", "", -1)
}
