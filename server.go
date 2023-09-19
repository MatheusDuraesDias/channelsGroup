package main

import (
	"database/sql"
    "log"
    "net/http"

    "github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var db *sql.DB

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
	UserIDs []int `json:"userIDs"`
}

type GroupChannels []Channel

var groupChannels GroupChannels

// func getChannels(c echo.Context) error {
// 	return c.JSON(http.StatusOK, groupChannels)
// }

func getChannels(c echo.Context) error {
    rows, err := db.Query("SELECT * FROM channel")
    if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }
    defer rows.Close()

    var channels GroupChannels
    for rows.Next() {
        var channel Channel
        if err := rows.Scan(&channel.ChannelURL, &channel.ChannelName, &channel.CoverURL, &channel.CreateBy, &channel.Description, &channel.MemberCount, &channel.JoinedMembers, &channel.MaxMessage, &channel.CreateAt, &channel.IsSuper, &channel.IsPublic, &channel.IsFreeze, &channel.IsEphemeral, &channel.IgnoreProfanity); err != nil {
            log.Println(err)
            return echo.NewHTTPError(http.StatusInternalServerError)
        }
        channels = append(channels, channel)
    }

    return c.JSON(http.StatusOK, channels)
}

func getChannelsUsers(c echo.Context) error {
    rows, err := db.Query("SELECT * FROM channel_users")
    if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }
    defer rows.Close()

    var channels GroupChannels
    for rows.Next() {
        var channel Channel
        if err := rows.Scan(&channel.ChannelURL, &channel.ChannelName, &channel.CoverURL, &channel.CreateBy, &channel.Description, &channel.MemberCount, &channel.JoinedMembers, &channel.MaxMessage, &channel.CreateAt, &channel.IsSuper, &channel.IsPublic, &channel.IsFreeze, &channel.IsEphemeral, &channel.IgnoreProfanity); err != nil {
            log.Println(err)
            return echo.NewHTTPError(http.StatusInternalServerError)
        }
        channels = append(channels, channel)
    }

    return c.JSON(http.StatusOK, channels)
}


// func getChannel(c echo.Context)  error {
// 	channelurl := c.Param("channelURL")
// 	for i := range groupChannels {
// 		if groupChannels[i].ChannelURL == channelurl {
// 			return c.JSON(http.StatusOK, groupChannels[i])
// 		}
// 	}
// 	return c.JSON(http.StatusBadRequest, nil)
// }

func getChannel(c echo.Context) error {
    channelURL := c.Param("channelURL")
    var channel Channel
    
    err := db.QueryRow("SELECT * FROM channel WHERE channelURL = $1", channelURL).Scan(&channel.ChannelURL, &channel.ChannelName, &channel.CoverURL, &channel.CreateBy, &channel.Description, &channel.MemberCount, &channel.JoinedMembers, &channel.MaxMessage, &channel.CreateAt, &channel.IsSuper, &channel.IsPublic, &channel.IsFreeze, &channel.IsEphemeral, &channel.IgnoreProfanity)
    if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusNotFound)
    }

    return c.JSON(http.StatusOK, channel)
}


// func postChannel(c echo.Context) error {
// 	channel := Channel{}
// 	err := c.Bind(&channel)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusUnprocessableEntity)
// 	}
// 	groupChannels = append(groupChannels, channel)
// 	return c.JSON(http.StatusCreated, groupChannels)
// }

// func postChannel(c echo.Context) error {
//     channel := Channel{}
//     err := c.Bind(&channel)
//     if err != nil {
//         return echo.NewHTTPError(http.StatusUnprocessableEntity)
//     }

//     _, err = db.Exec("INSERT INTO channel (channelURL, channelName, coverURL, createBy, description, memberCount, joinedMembers, maxMessage, createAt, isSuper, isPublic, isFreeze, isEphemeral, ignoreProfanity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
//         channel.ChannelURL, channel.ChannelName, channel.CoverURL, channel.CreateBy, channel.Description, channel.MemberCount, channel.JoinedMembers, channel.MaxMessage, channel.CreateAt, channel.IsSuper, channel.IsPublic, channel.IsFreeze, channel.IsEphemeral, channel.IgnoreProfanity)
//     if err != nil {
//         log.Println(err)
//         return echo.NewHTTPError(http.StatusInternalServerError)
//     }

//     return c.JSON(http.StatusCreated, channel)
// }

func postChannel(c echo.Context) error {
    channel := Channel{}
    err := c.Bind(&channel)
    if err != nil {
        return echo.NewHTTPError(http.StatusUnprocessableEntity)
    }

    _, err = db.Exec("INSERT INTO channel (channelURL, channelName, coverURL, createBy, description, memberCount, joinedMembers, maxMessage, createAt, isSuper, isPublic, isFreeze, isEphemeral, ignoreProfanity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
        channel.ChannelURL, channel.ChannelName, channel.CoverURL, channel.CreateBy, channel.Description, channel.MemberCount, channel.JoinedMembers, channel.MaxMessage, channel.CreateAt, channel.IsSuper, channel.IsPublic, channel.IsFreeze, channel.IsEphemeral, channel.IgnoreProfanity)
    if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    var channelID string
    err = db.QueryRow("SELECT channelURL FROM channel WHERE channelURL = $1", channel.ChannelURL).Scan(&channelID)
    if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    for _, userID := range channel.UserIDs {
        _, err = db.Exec("INSERT INTO channel_users (channel_id, user_id) VALUES ($1, $2)", channelID, userID)
        if err != nil {
            log.Println(err)
            return echo.NewHTTPError(http.StatusInternalServerError)
        }
    }

    return c.JSON(http.StatusOK, "Channel created successfully")
}



// func deleteChannel(c echo.Context) error {
// 	channelurl := c.Param("channelURL")
// 	for i := range groupChannels {
// 		if groupChannels[i].ChannelURL == channelurl {
// 			groupChannels = append(groupChannels[:i], groupChannels[i+1:]...)
// 			return c.JSON(http.StatusOK, groupChannels)
// 		}
// 	}
// 	return c.JSON(http.StatusBadRequest, nil)
// }

func deleteChannel(c echo.Context) error {
    channelURL := c.Param("channelURL")
    
    _, err := db.Exec("DELETE FROM channel WHERE channelURL = $1", channelURL)
    if err != nil {
        log.Println(err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    return c.JSON(http.StatusOK, "Channel deleted")
}

func main() {

    connStr := "user=postgres password=1234 dbname=dbChannels host=localhost port=5432 sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
		log.Println(err)
    }
    defer db.Close()

    e := echo.New()

	e.GET("/channels", getChannels)
	e.POST("/channels", postChannel)
	e.GET("/channels/:channelURL", getChannel)
	e.DELETE("/channels/:channelURL", deleteChannel)

    e.Logger.Fatal(e.Start(":1323"))
}