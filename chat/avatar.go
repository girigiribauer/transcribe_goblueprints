package main

import (
	"errors"
)

// ErrNoAvatarURL は Avatar インスタンスがアバターの URL を返すことができない
// 場合に発生するエラーです
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatar はユーザーのプロフィール画像を表す型です。
type Avatar interface {
	// GetAvatarURL は指定されたクライアントのアバターのURLを返します。
	// 問題が発生した場合にエラーを返します。特に、URLを取得できなかった
	// 場合にはErrNoAvatarURLを返します。
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar は認証サービス用の Avatar です
type AuthAvatar struct{}

// UseAuthAvatar は AuthAvatar のインスタンスです
var UseAuthAvatar AuthAvatar

// GetAvatarURL はアバター用の画像URLを返します
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar は gravatar.com 用の Avatar です
type GravatarAvatar struct{}

// UseGravatar は GravatarAvatar のインスタンスです
var UseGravatar GravatarAvatar

// GetAvatarURL は gravatar.com 用の画像URLを返します
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// FileSystemAvatar は自らアップロードした場合の Avatar です
type FileSystemAvatar struct{}

// UseFileSystemAvatar は FileSystemAvatar のインスタンスです
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL は自らアップロードした画像のURLを返します
func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
