package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
