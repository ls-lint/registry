package main

type Publish struct {
	Token   string `json:"token" binding:"required"`
	Package string `json:"package" binding:"required"`
	Public  bool   `json:"public" binding:"required"`
	Version string `json:"version"`
	Data    []byte `json:"data"`
}
