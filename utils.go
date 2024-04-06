package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	fauxgl "github.com/hawl1/brickgl"
)

// LoadMeshFromURL loads mesh from url
func LoadMeshFromURL(url string) *fauxgl.Mesh {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	mesh, _ := fauxgl.LoadOBJFromReader(resp.Body)

	return mesh
}

// LoadTexture loads texture from URL
func LoadTexture(url string) fauxgl.Texture {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return fauxgl.TexFromBytes(body)
}

// LoadItem loads item from url
func LoadItem(item int, scene *fauxgl.Scene) {
	if item != 0 {
		resp, err := http.Get(fmt.Sprintf("https://api.brick-hill.com/v1/assets/getPoly/1/%d", item))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var data []map[string]string
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			panic(err)
		}

		mesh := data[0]["mesh"]
		mesh = mesh[len("asset://"):]

		texture := data[0]["texture"]
		texture = texture[len("asset://"):]

		scene.AddObject(&fauxgl.Object{
			Mesh:    LoadMeshFromURL("https://api.brick-hill.com/v1/assets/get/" + mesh),
			Texture: LoadTexture("https://api.brick-hill.com/v1/assets/get/" + texture),
		})
	}
}