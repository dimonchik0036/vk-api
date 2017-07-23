package vkapi

import (
	"bytes"
	"encoding/json"
	"github.com/technoweenie/multipartstreamer"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// FileBytes contains information about a set of bytes to upload
// as a File.
type FileBytes struct {
	Name  string
	Bytes []byte
}

// FileReader contains information about a reader to upload as a File.
// If Size is -1, it will read the entire Reader into memory to
// calculate a Size.
type FileReader struct {
	Name   string
	Reader io.Reader
	Size   int64
}

type ServerResponse struct {
	Response
	Server int64  `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

func (client *Client) UploadFile(url string, fieldname string, file interface{}) (ServerResponse, *Error) {
	ms := multipartstreamer.New()

	switch f := file.(type) {
	case string:
		fileHandle, err := os.Open(f)
		if err != nil {
			return ServerResponse{}, NewError(ErrBadCode, err.Error())
		}
		defer fileHandle.Close()

		fi, err := os.Stat(f)
		if err != nil {
			return ServerResponse{}, NewError(ErrBadCode, err.Error())
		}

		ms.WriteReader(fieldname, fileHandle.Name(), fi.Size(), fileHandle)
	case FileBytes:
		buf := bytes.NewBuffer(f.Bytes)
		ms.WriteReader(fieldname, f.Name, int64(len(f.Bytes)), buf)
	case FileReader:
		if f.Size != -1 {
			ms.WriteReader(fieldname, f.Name, f.Size, f.Reader)

			break
		}

		data, err := ioutil.ReadAll(f.Reader)
		if err != nil {
			return ServerResponse{}, NewError(ErrBadCode, err.Error())
		}

		buf := bytes.NewBuffer(data)

		ms.WriteReader(fieldname, f.Name, int64(len(data)), buf)
	default:
		return ServerResponse{}, NewError(ErrBadCode, "Bad input")
	}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return ServerResponse{}, NewError(ErrBadCode, err.Error())
	}

	ms.SetupRequest(req)

	res, err := client.apiClient.httpClient.Do(req)
	if err != nil {
		return ServerResponse{}, NewError(ErrBadCode, err.Error())
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ServerResponse{}, NewError(ErrBadCode, err.Error())
	}
	client.apiClient.logPrintf("upload %s: %s", fieldname, string(bytes))

	response := ServerResponse{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return ServerResponse{}, NewError(ErrBadCode, err.Error())
	}

	return response, nil
}

type UploadServer struct {
	UploadURL string `json:"upload_url"`
}

type UploadServerPhoto struct {
	UploadServer
	UserID  int64 `json:"user_id"`
	AlbumID int64 `json:"album_id"`
}

func (client *Client) GetMessagesUploadServer() (server UploadServerPhoto, err *Error) {
	res, err := client.Do(NewRequest("photos.getMessagesUploadServer", "", nil))
	if err != nil {
		return UploadServerPhoto{}, err
	}

	if err := res.To(&server); err != nil {
		return UploadServerPhoto{}, NewError(ErrBadCode, err.Error())
	}

	return
}

func (client *Client) SaveMessagesPhoto(response ServerResponse) (Photo, *Error) {
	values := url.Values{}
	values.Add("photo", response.Photo)
	values.Add("server", ConcatInt64ToString(response.Server))
	values.Add("hash", response.Hash)

	res, err := client.Do(NewRequest("photos.saveMessagesPhoto", "", values))
	if err != nil {
		return Photo{}, err
	}

	var photo []Photo
	if err := res.To(&photo); err != nil {
		return Photo{}, NewError(ErrBadCode, err.Error())
	}

	return photo[0], nil
}
