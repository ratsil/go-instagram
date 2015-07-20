// Copyright 2013 The go-instagram AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package instagram

import (
	"fmt"
)

// RelationshipsService handles communication with the user's relationships related
// methods of the Instagram API.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/
type RelationshipsService struct {
	client *Client
}

// Relationship represents relationship authenticated user with another user.
type Relationship struct {
	// Current user's relationship to another user. Can be "follows", "requested", or "none".
	OutgoingStatus string `json:"outgoing_status,omitempty"`

	// A user's relationship to current user. Can be "followed_by", "requested_by",
	// "blocked_by_you", or "none".
	IncomingStatus string `json:"incoming_status,omitempty"`

	// Undocumented part of the API, though was stable at least from 2012-2015
	// Informs whether the target user is a private user
	TargetUserIsPrivate bool `json:"target_user_is_private,omitempty"`
}

// Follows gets the list of users this user follows. If empty string is
// passed then it refers to `self` or curret authenticated user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#get_users_follows
func (s *RelationshipsService) Follows(userID string) ([]User, *ResponsePagination, error) {
	var u string
	if userID != "" {
		u = fmt.Sprintf("users/%v/follows", userID)
	} else {
		u = "users/self/follows"
	}

	req, err := s.client.NewRequest("GET", u, "")
	if err != nil {
		return nil, nil, err
	}

	users := new([]User)

	_, err = s.client.Do(req, users)
	if err != nil {
		return nil, nil, err
	}

	page := new(ResponsePagination)
	if s.client.Response.Pagination != nil {
		page = s.client.Response.Pagination
	}

	return *users, page, err
}

// FollowedBy gets the list of users this user is followed by. If empty string is
// passed then it refers to `self` or curret authenticated user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#get_users_followed_by
func (s *RelationshipsService) FollowedBy(userID string) (aRetVal []User, err error) {
	var u string
	if userID != "" {
		u = fmt.Sprintf("users/%v/followed-by", userID)
	} else {
		u = "users/self/followed-by"
	}
	
	aRetVal = []User{}
	for{
		req, err := s.client.NewRequest("GET", u, "")
		if err != nil {
			return nil, err
		}
	
		aUsers := []User{}
		_, err = s.client.Do(req, &aUsers)
		if err != nil {
			return nil, err
		}
		aRetVal = append(aRetVal, aUsers...)
		if s.client.Response.Pagination != nil && 0 < len(s.client.Response.Pagination.NextURL) {
			u = s.client.Response.Pagination.NextURL
			continue
		}
		break
	}
	return aRetVal, err
}

// RequestedBy lists the users who have requested this user's permission to follow.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#get_incoming_requests
func (s *RelationshipsService) RequestedBy() ([]User, *ResponsePagination, error) {
	u := "users/self/requested-by"
	req, err := s.client.NewRequest("GET", u, "")
	if err != nil {
		return nil, nil, err
	}

	users := new([]User)

	_, err = s.client.Do(req, users)
	if err != nil {
		return nil, nil, err
	}

	page := new(ResponsePagination)
	if s.client.Response.Pagination != nil {
		page = s.client.Response.Pagination
	}

	return *users, page, err
}

// Relationship gets information about a relationship to another user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#get_relationship
func (s *RelationshipsService) Relationship(userID string) (*Relationship, error) {
	return relationshipAction(s, userID, "", "GET")
}

// Follow a user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#post_relationship
func (s *RelationshipsService) Follow(userID string) (*Relationship, error) {
	return relationshipAction(s, userID, "follow", "POST")
}

// Unfollow a user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#post_relationship
func (s *RelationshipsService) Unfollow(userID string) (*Relationship, error) {
	return relationshipAction(s, userID, "unfollow", "POST")
}

// Block a user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#post_relationship
func (s *RelationshipsService) Block(userID string) (*Relationship, error) {
	return relationshipAction(s, userID, "block", "POST")
}

// Unblock a user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#post_relationship
func (s *RelationshipsService) Unblock(userID string) (*Relationship, error) {
	return relationshipAction(s, userID, "unblock", "POST")
}

// Approve a user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#post_relationship
func (s *RelationshipsService) Approve(userID string) (*Relationship, error) {
	return relationshipAction(s, userID, "approve", "POST")
}

// Deny a user.
//
// Instagram API docs: http://instagram.com/developer/endpoints/relationships/#post_relationship
func (s *RelationshipsService) Deny(userID string) (*Relationship, error) {
	return relationshipAction(s, userID, "deny", "POST")
}

func relationshipAction(s *RelationshipsService, userID, action, method string) (*Relationship, error) {
	u := fmt.Sprintf("users/%v/relationship", userID)
	if action != "" {
		action = "action=" + action
	}
	req, err := s.client.NewRequest(method, u, action)
	if err != nil {
		return nil, err
	}

	rel := new(Relationship)
	_, err = s.client.Do(req, rel)
	return rel, err
}
