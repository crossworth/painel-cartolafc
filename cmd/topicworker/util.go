package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	vkapi "github.com/himidori/golang-vk-api"
)

func downloadTopic(client *vkapi.VKClient, groupID int, topicID int) (Topic, error) {
	var topic Topic

	params := url.Values{}
	params.Set("topic_ids", strconv.Itoa(topicID))
	params.Set("extended", "1")

	topicResult, err := client.BoardGetTopics(groupID, 1, params)
	if err != nil {
		return topic, err
	}

	profilesUsers := mapUsers(topicResult.Profiles)

	topic.ID = topicID
	topic.Title = topicResult.Topics[0].Title
	topic.IsClosed = intToBool(topicResult.Topics[0].IsClosed)
	topic.IsFixed = intToBool(topicResult.Topics[0].IsFixed)
	topic.CreatedAt = topicResult.Topics[0].Created
	topic.CreatedBy = vkUserToProfile(profilesUsers[topicResult.Topics[0].CreatedBy])
	topic.UpdatedAt = topicResult.Topics[0].Updated
	topic.UpdatedBy = vkUserToProfile(profilesUsers[topicResult.Topics[0].UpdatedBy])
	topic.Profiles = make(map[int]Profile)

	return topic, err
}

func downloadComments(client *vkapi.VKClient, groupID int, topicID int, startCommentID int) ([]Comment, int, []Profile, Poll, error) {
	var commentsResult []Comment
	var profilesResult []Profile
	var pollResult Poll

	commentsParams := url.Values{}
	commentsParams.Set("extended", "1")
	commentsParams.Set("need_likes", "1")

	if startCommentID != 0 {
		commentsParams.Set("start_comment_id", strconv.Itoa(startCommentID))
	}

	comments, err := client.BoardGetComments(groupID, topicID, 100, commentsParams)
	if err != nil {
		return commentsResult, 0, profilesResult, pollResult, err
	}

	if comments.Poll != nil {
		pollResult = Poll{
			ID:       comments.Poll.ID,
			Question: comments.Poll.Question,
			Votes:    comments.Poll.Votes,
			Answers:  vkPollAnswerToPollAnswer(comments.Poll.Answers),
			Multiple: comments.Poll.Multiple,
			EndDate:  comments.Poll.EndDate,
			Closed:   comments.Poll.Closed,
		}
	}

	for i := range comments.Profiles {
		profilesResult = append(profilesResult, vkUserToProfile(*comments.Profiles[i]))
	}

	for i := range comments.Comments {
		commentsResult = append(commentsResult, vkCommentToComment(*comments.Comments[i]))
	}

	return commentsResult, comments.Count, profilesResult, pollResult, nil
}

func intToBool(i int) bool {
	if i > 0 {
		return true
	}

	return false
}

func vkUserToProfile(user vkapi.User) Profile {
	return Profile{
		ID:         user.UID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		ScreenName: user.ScreenName,
		Photo:      user.Photo100,
	}
}

func mapUsers(profiles []*vkapi.User) map[int]vkapi.User {
	users := make(map[int]vkapi.User)

	for i := range profiles {
		users[profiles[i].UID] = *profiles[i]
	}

	return users
}

func vkCommentToComment(comment vkapi.TopicComment) Comment {
	cmt := Comment{
		ID:         comment.ID,
		FromID:     comment.FromID,
		Date:       comment.Date,
		Text:       comment.Text,
		ReplyToUID: comment.ReplyToUID,
		ReplyToCID: comment.ReplyToUID,
	}

	if comment.Likes != nil {
		cmt.Likes = comment.Likes.Count
	}

	for i := range comment.Attachments {
		switch comment.Attachments[i].Type {
		case "photo":
			cmt.Attachments = append(cmt.Attachments, getBestPhoto(*comment.Attachments[i].Photo))
		case "sticker":
			cmt.Attachments = append(cmt.Attachments, getBestSticker(*comment.Attachments[i].Sticker))
		case "video":
			cmt.Attachments = append(cmt.Attachments, fmt.Sprintf("https://vk.com/video?z=video%d_%d%%2F%s", comment.Attachments[i].Video.OwnerID, comment.Attachments[i].Video.ID, comment.Attachments[i].Video.AccessKey))
		case "audio":
			// NOTE(Pedro): we save the JSON
			// since we dont have a good way to "make" a audio link
			output, err := json.Marshal(comment.Attachments[i].Audio)
			if err != nil {
				continue
			}

			cmt.Attachments = append(cmt.Attachments, string(output))
		}
	}

	return cmt
}

func getBestPhoto(attachment vkapi.AttachmentPhoto) string {
	best := attachment.Sizes[0]

	for i := range attachment.Sizes {
		s := attachment.Sizes[i].Width * attachment.Sizes[i].Height
		b := best.Width * best.Height
		if s > b {
			best = attachment.Sizes[i]
		}
	}

	return best.Url
}

func getBestSticker(attachment vkapi.AttachmentSticker) string {
	best := attachment.Images[0]

	for i := range attachment.Images {
		s := attachment.Images[i].Width * attachment.Images[i].Height
		b := best.Width * best.Height
		if s > b {
			best = attachment.Images[i]
		}
	}

	return best.Url
}

func vkPollAnswerToPollAnswer(answers []*vkapi.PollAnswer) []PollAnswer {
	var result []PollAnswer

	for i := range answers {
		result = append(result, PollAnswer{
			ID:    answers[i].ID,
			Text:  answers[i].Text,
			Votes: answers[i].Votes,
			Rate:  answers[i].Rate,
		})
	}

	return result
}
