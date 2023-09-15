package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

type Channel struct {
	ChannelURL 		string `json:"channelURL"`
	ChannelName 	string `json:"channelName"`
	CoverURL 		string `json:"coverURL"`
	CreateBy		string `json:"createBy"`
	Description		string `json:"description	"`
	MemberCount		int `json:"memberCount"`
	JoinedMembers	int `json:"joinedMembers"`
	MaxMessage		int `json:"maxMessage"`
	CreateAt		string `json:"createAt"`
	IsSuper			bool `json:"isSuper"`
	IsPublic		bool `json:"isPublic"`
	IsFreeze		bool `json:"isFreeze"`
	IsEphemeral		bool `json:"isEphemeral"`
	IgnoreProfanity	bool `json:"ignoreProfanity"`
}

type GroupChannels []Channel

var groupChannels GroupChannels

func getChannels(c echo.Context) error {
	return c.JSON(http.StatusOK, groupChannels)
}

func getChannel(c echo.Context)  error {
	channelurl := c.Param("channelURL")
	for i := range groupChannels {
		if groupChannels[i].ChannelURL == channelurl {
			return c.JSON(http.StatusOK, groupChannels[i])
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func postChannel(c echo.Context) error {
	channel := Channel{}
	err := c.Bind(&channel)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}
	groupChannels = append(groupChannels, channel)
	return c.JSON(http.StatusCreated, groupChannels)
}

func deleteChannel(c echo.Context) error {
	channelurl := c.Param("channelURL")
	for i := range groupChannels {
		if groupChannels[i].ChannelURL == channelurl {
			groupChannels = append(groupChannels[:i], groupChannels[i+1:]...)
			return c.JSON(http.StatusOK, groupChannels)
		}
	}
	return c.JSON(http.StatusBadRequest, nil)
}

func main() {
    e := echo.New()
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
	e.GET("/teste", func(c echo.Context) error {
        return c.String(http.StatusOK, "Testando")
    })

	e.GET("/channels", getChannels)
	e.POST("/channels", postChannel)
	e.GET("/channels/:channelURL", getChannel)
	e.DELETE("/channels/:channelURL", deleteChannel)

    e.Logger.Fatal(e.Start(":1323"))
}