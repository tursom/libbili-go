/*
 * Copyright (c) 2022 tursom. All rights reserved.
 * Use of this source code is governed by a GPL-3
 * license that can be found in the LICENSE file.
 */

package live

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tursom/GoCollections/exceptions"
)

const (
	sendDanmuURL   = "https://api.live.bilibili.com/msg/send"
	danmuColorsURL = "https://api.live.bilibili.com/xlive/web-room/v1/dM/GetDMConfigByGroup?room_id=%d"
)

type (
	/*
	 * Room interface to access an specific live room
	 */
	Room interface {
		// ID get live room id
		ID() uint32
		// SetCookie set user cookie
		SetCookie(cookie string)
		// Send send danmu
		Send(msg string) (*DanmuResp, exceptions.Exception)
		// SendDanmu send danmu with more infos
		SendDanmu(danmu *Danmu) (*DanmuResp, exceptions.Exception)
		// GetDanmuColors get available colors of danmu
		GetDanmuColors() (*DanmuColors, exceptions.Exception)
	}

	roomImpl struct {
		id     uint32
		cookie string
	}

	/*
	 * Danmu danmu send info
	 */
	Danmu struct {
		Bubble    int32  `json:"bubble,omitempty"`
		Msg       string `json:"msg,omitempty"`
		Color     string `json:"color,omitempty"`
		Mode      int32  `json:"mode,omitempty"`
		Fontsize  int32  `json:"fontsize,omitempty"`
		Rnd       int64  `json:"rnd,omitempty"`
		RoomId    uint32 `json:"roomid,omitempty"`
		Csrf      string `json:"csrf,omitempty"`
		CsrfToken string `json:"csrf_token,omitempty"`
	}

	/*
	 * Danmu danmu send response
	 */
	DanmuResp struct {
		Code int `json:"code"`
		Data struct {
			ModeInfo struct {
				Mode           int    `json:"mode"`
				ShowPlayerType int    `json:"show_player_type"`
				Extra          string `json:"extra"`
			} `json:"mode_info"`
		} `json:"data"`
		Message string `json:"message"`
		Msg     string `json:"msg"`
	}

	/*
	 * Danmu available danmu color
	 */
	DanmuColors struct {
		Code int `json:"code"`
		Data struct {
			Group []struct {
				Name  string `json:"name"`
				Sort  int    `json:"sort"`
				Color []struct {
					Name     string `json:"name"`
					Color    string `json:"color"`
					ColorHex string `json:"color_hex"`
					Status   int    `json:"status"`
					Weight   int    `json:"weight"`
					ColorId  int    `json:"color_id"`
					Origin   int    `json:"origin"`
				} `json:"color"`
			} `json:"group"`
			Mode []struct {
				Name   string `json:"name"`
				Mode   int    `json:"mode"`
				Type   string `json:"type"`
				Status int    `json:"status"`
			} `json:"mode"`
		} `json:"data"`
		Message string `json:"message"`
		Msg     string `json:"msg"`
	}
)

var (
	client = resty.New()
)

func NewRoom(id uint32) Room {
	return &roomImpl{
		id: id,
	}
}

func (r *roomImpl) SetCookie(cookie string) {
	r.cookie = cookie
}

func (r *roomImpl) ID() uint32 {
	return r.id
}

func (r *roomImpl) Send(msg string) (*DanmuResp, exceptions.Exception) {
	// check cookie status
	if r.cookie == "" {
		// return err that no cookie set
		return nil, exceptions.NewIllegalParameterException("cookie not set", nil)
	}

	// send danmu
	return r.SendDanmu(&Danmu{
		Bubble:    0,
		Msg:       msg,
		Color:     "16777215",
		Mode:      1,
		Fontsize:  25,
		Rnd:       time.Now().Unix(),
		RoomId:    r.id,
		Csrf:      "c1b21617a15daf838f505271ff8f5204",
		CsrfToken: "c1b21617a15daf838f505271ff8f5204",
	})
}

func (r *roomImpl) SendDanmu(danmu *Danmu) (*DanmuResp, exceptions.Exception) {
	// check cookie status
	if r.cookie == "" {
		// return err that no cookie set
		return nil, exceptions.NewIllegalParameterException("cookie not set", nil)
	}

	// build request object
	request := client.R()

	// build danmu param with multipart/form-data
	form, boundary, exception := multipartForm(danmu)
	if exception != nil {
		return nil, exception
	}

	request.SetBody(form)

	// build request object
	//request, err := http.NewRequest("POST", sendDanmuURL, form)
	//if err != nil {
	//	return nil, exceptions.Package(err)
	//}

	// build request headers
	//request.Header.Add("Accept", "*/*")
	request.Header.Set("Cookie", r.cookie)
	//request.Header.Set("Origin", "https://live.bilibili.com")
	//request.Header.Set("Referer", fmt.Sprintf("https://li|ve.bilibili.com/%d?spm_id_from=444.41.live_users.item.click", r.id))
	//request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	request.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
	//request.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	//request.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	//request.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	//request.Header.Set("Sec-Fetch-Dest", "empty")
	//request.Header.Set("Sec-Fetch-Mode", "cors")
	//request.Header.Set("Sec-Fetch-Site", "same-site")
	//request.Header.Set("content-length", strconv.Itoa(len(form)))

	// send request
	//do, err := http.DefaultClient.Do(request)
	do, err := request.Post(sendDanmuURL)
	if err != nil {
		return nil, exceptions.Package(err)
	}

	// check response http status
	if do.StatusCode() != 200 {
		// return err
		return nil, exceptions.NewPackageException(fmt.Sprintf("send response status failed: %d", do.StatusCode()),
			exceptions.Cfg().SetCause(do.StatusCode))
	}

	// parse response data
	var resp DanmuResp
	err = json.Unmarshal(do.Body(), &resp)
	if err != nil {
		return nil, exceptions.Package(err)
	}

	return &resp, nil
}

// multipartForm build multipart/form-data by send danmu of body data
func multipartForm(danmu *Danmu) (formData []byte, boundary string, exception exceptions.Exception) {
	formBuffer := bytes.NewBuffer(nil)
	formWriter := multipart.NewWriter(formBuffer)

	err := formWriter.WriteField("bubble", strconv.Itoa(int(danmu.Bubble)))
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("msg", danmu.Msg)
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("color", danmu.Color)
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("mode", strconv.Itoa(int(danmu.Mode)))
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("fontsize", strconv.Itoa(int(danmu.Fontsize)))
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("rnd", strconv.FormatInt(danmu.Rnd, 10))
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("roomid", strconv.Itoa(int(danmu.RoomId)))
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("csrf", danmu.Csrf)
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.WriteField("csrf_token", danmu.CsrfToken)
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	err = formWriter.Close()
	if err != nil {
		return nil, "", exceptions.Package(err)
	}

	return formBuffer.Bytes(), formWriter.Boundary(), nil
}

// GetDanmuColors get available colors of danmu sending
func (r *roomImpl) GetDanmuColors() (*DanmuColors, exceptions.Exception) {
	// build request URL
	url := fmt.Sprintf(danmuColorsURL, r.id)

	// build request object
	request := client.R()

	// check cookie status
	if r.cookie != "" {
		request.Header.Set("Cookie", r.cookie)
	}

	// send http request
	response, err := request.Get(url)
	if err != nil {
		return nil, exceptions.Package(err)
	}

	// parse http response
	var colors DanmuColors
	err = json.Unmarshal(response.Body(), &colors)
	if err != nil {
		return nil, exceptions.Package(err)
	}

	return &colors, nil
}
