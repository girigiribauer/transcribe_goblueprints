package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL は Avatar インスタンスがアバターの URL を返すことができない
// 場合に発生するエラーです
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません。")

// Avatar はユーザーのプロフィール画像を表す型です。
type Avatar interface {
	// GetAvatarURL は指定されたクライアントのアバターのURLを返します。
	// 問題が発生した場合にエラーを返します。特に、URLを取得できなかった
	// 場合にはErrNoAvatarURLを返します。
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars は複数の Avatar 実装を順に試すための型です。
type TryAvatars []Avatar

// GetAvatarURL は TryAvatars のスライスが持つ実装を順に試します
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

// AuthAvatar は認証サービス用の Avatar です
type AuthAvatar struct{}

// UseAuthAvatar は AuthAvatar のインスタンスです
var UseAuthAvatar AuthAvatar

// GetAvatarURL はアバター用の画像URLを返します
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar は gravatar.com 用の Avatar です
type GravatarAvatar struct{}

// UseGravatar は GravatarAvatar のインスタンスです
var UseGravatar GravatarAvatar

// GetAvatarURL は gravatar.com 用の画像URLを返します
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar は自らアップロードした場合の Avatar です
type FileSystemAvatar struct{}

// UseFileSystemAvatar は FileSystemAvatar のインスタンスです
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL は自らアップロードした画像のURLを返します
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}
