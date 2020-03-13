package main

type Publish struct {
	Token   string `json:"token" binding:"required"`
	Package string `json:"package" binding:"required"`
	Public  bool   `json:"public" binding:"required"`
	Tag     string `json:"tag"`
	Data    []byte `json:"data" binding:"required"`
}
