package model

type Model interface {
    Validate() (bool, map[string]string)
}