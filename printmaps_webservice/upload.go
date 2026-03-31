// Upload handler

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/printmaps/printmaps/pd"
)

/*
uploadUserdata allows the upload of an user data file (e.g. gpx file).
*/
func uploadUserdata(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var pmData pd.PrintmapsData
	var pmErrorList pd.PrintmapsErrorList

	id := params.ByName("id")

	// verify ID
	_, err := uuid.FromString(id)
	if err != nil {
		appendError(&pmErrorList, "4001", "error = "+err.Error(), "")
	}

	// map directory must exist
	if len(pmErrorList.Errors) == 0 {
		if !pd.IsExistMapDirectory(id) {
			appendError(&pmErrorList, "4002", "requested ID not found: "+id, id)
		}
	}

	if len(pmErrorList.Errors) == 0 {
		if err := pd.ReadMetadata(&pmData, id); err != nil {
			if os.IsNotExist(err) {
				appendError(&pmErrorList, "4002", "requested ID not found: "+id, id)
			} else {
				message := fmt.Sprintf("error <%v> at readMetadata(), id = <%s>", err, id)
				http.Error(writer, message, http.StatusInternalServerError)
				log.Printf("Response %d - %s", http.StatusInternalServerError, message)
				return
			}
		}
	}

	userfileName := ""
	userfileSize := int64(-1)
	filelimit := int64(224 * 1024 * 1024)

	if len(pmErrorList.Errors) == 0 {
		// Limit the size of the request body to prevent disk exhaustion (DoS).
		// An additional 1 MB (1<<20) overhead is granted for multipart form boundaries and headers.
		request.Body = http.MaxBytesReader(writer, request.Body, filelimit+(1<<20))

		// input file
		file, header, err := request.FormFile("file")
		if err != nil {
			message := fmt.Sprintf("upload failed or exceeded max upload size (%d bytes): %v", filelimit, err)
			log.Printf("uploadUserdata(): %s", message)
			appendError(&pmErrorList, "7001", fmt.Sprintf("max upload size = %d bytes", filelimit), id)
		} else {
			defer func() { _ = file.Close() }()

			_, userfileName = filepath.Split(header.Filename)

			filename := filepath.Join(pd.PathWorkdir, pd.PathMaps, pmData.Data.ID, userfileName)
			out, err := os.Create(filename)
			if err != nil {
				message := fmt.Sprintf("error <%v> at os.Create(), file = <%s>", err, filename)
				http.Error(writer, message, http.StatusInternalServerError)
				log.Printf("Response %d - %s", http.StatusInternalServerError, message)
				return
			}

			// write content from POST to file
			bytesWritten, err := io.Copy(out, file)
			_ = out.Close()

			if err != nil {
				message := fmt.Sprintf("error <%v> at io.Copy(), file = <%s> (exceeds max size?)", err, filename)
				log.Printf("uploadUserdata(): %s", message)
				// Remove partially written file
				_ = os.Remove(filename)
				appendError(&pmErrorList, "7001", fmt.Sprintf("max upload size = %d bytes", filelimit), id)
			} else {
				userfileSize = bytesWritten

				// verify security of uploaded file
				err := verifyUploadedFileMimetype(filename)
				if err != nil {
					log.Printf("insecure user file <%s> rejected: %v", filename, err)
					appendError(&pmErrorList, "7002", "only data or image files are accepted", id)
					errRemove := os.Remove(filename)
					if errRemove != nil {
						log.Printf("unexpected error <%s> os.Remove(), file = <%s>", errRemove, filename)
					}
				}
			}
		}
	}

	if len(pmErrorList.Errors) == 0 {
		// upload request ok (user data file created)
		writer.WriteHeader(http.StatusCreated)
		message := fmt.Sprintf("file <%s, %d bytes> successfully uploaded", userfileName, userfileSize)
		_, _ = writer.Write([]byte(message))
		log.Printf("uploadUserdata(): %s", message)
	} else {
		// request not ok, response with error list
		content, err := json.MarshalIndent(pmErrorList, pd.IndentPrefix, pd.IndexString)
		if err != nil {
			message := fmt.Sprintf("error <%v> at json.MarshalIndent()", err)
			http.Error(writer, message, http.StatusInternalServerError)
			log.Printf("Response %d - %s", http.StatusInternalServerError, message)
			return
		}
		writer.Header().Set("Content-Type", pd.JSONAPIMediaType)
		writer.Header().Set("Content-Length", strconv.Itoa(len(content)))
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write(content)
	}
}

/*
verifyUploadedFileMimetype verifies if a file is 'secure'.
*/
func verifyUploadedFileMimetype(filename string) error {
	mtype, err := mimetype.DetectFile(filename)
	if err != nil {
		return fmt.Errorf("error detecting mimetype: %w", err)
	}

	for m := mtype; m != nil; m = m.Parent() {
		if m.String() == "application/x-executable" || m.String() == "application/x-mach-binary" || m.String() == "application/x-dosexec" {
			return errors.New("executable files are not permitted")
		}
	}

	return nil
}
