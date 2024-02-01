package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"

	// "github.com/gofrs/uuid"
	fdk "github.com/fnproject/fdk-go"
	fauxgl "github.com/hawl1/brickgl"
)

const (
	scale = 3
	fovy  = 50
	near  = 0.1
	far   = 1000
)

var (
	eye    = fauxgl.V(-2, 0.85, 0.75)
	center = fauxgl.V(0, 0.06, 0)
	up     = fauxgl.V(0, 1, 0)
	light  = fauxgl.V(-4, 6, 0).Normalize()
)

// RenderEvent input data to lambda to return an ImageResponse
type RenderEvent struct {
	AvatarJSON string `json:"avatarJSON"`
	Size       int    `json:"size"`
}

// Avatar avatar
type Avatar struct {
	UserID int                    `json:"user_id"`
	Items  map[string]interface{} `json:"items"`
	Colors map[string]string      `json:"colors"`
}

// ImageResponse lambda response for a base64 encoded render
type ImageResponse struct {
	// gonna fix this, can stay for now UUID  string `json:"uuid"`
	Image string `json:"image"`
}

// LoadMeshFromURL loads mesh from url
func LoadMeshFromURL(url string) *fauxgl.Mesh {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	mesh, _ := fauxgl.LoadOBJFromReader(resp.Body)

	return mesh
}

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

		texture := data[0]["texture"]
		texture = texture[len("asset://"):]

		mesh := data[0]["mesh"]
		mesh = mesh[len("asset://"):]

		scene.AddObject(&fauxgl.Object{
			Mesh:    LoadMeshFromURL("https://api.brick-hill.com/v1/assets/get/" + mesh),
			Texture: LoadTexture("https://api.brick-hill.com/v1/assets/get/" + texture),
		})
	}
}

func main() {
	fdk.Handle(fdk.HandlerFunc(HandleRenderEvent))
}

// HandleRenderEvent function to process the rendering
func HandleRenderEvent(ctx context.Context, in io.Reader, out io.Writer) {
	var e RenderEvent
	err := json.NewDecoder(in).Decode(&e)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var avatarJSON string
	if e.AvatarJSON == "" {
		// Use the default JSON string if AvatarJSON is empty
		avatarJSON = "{\"user_id\":13,\"items\":{\"face\":0,\"hats\":[20121,0,0,0,0],\"head\":0,\"tool\":0,\"pants\":0,\"shirt\":364985,\"figure\":0,\"tshirt\":0},\"colors\":{\"head\":\"eab372\",\"torso\":\"85ad00\",\"left_arm\":\"eab372\",\"left_leg\":\"37302c\",\"right_arm\":\"eab372\",\"right_leg\":\"37302c\"}}"
	} else {
		avatarJSON = e.AvatarJSON
	}

	avatar := Avatar{}
	err = json.Unmarshal([]byte(avatarJSON), &avatar)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	aspect := float64(e.Size) / float64(e.Size)
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	context := fauxgl.NewContext(e.Size, e.Size, scale, shader)
	scene := fauxgl.NewScene(context)

	shirt := fauxgl.NewImageTexture(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	if shirtValue, ok := avatar.Items["shirt"].(float64); ok && shirtValue != 0 {
		resp, err := http.Get(fmt.Sprintf("https://api.brick-hill.com/v1/assets/getPoly/1/%d", int(shirtValue)))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var data []map[string]string
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			panic(err)
		}

		texture := data[0]["texture"]
		texture = texture[len("asset://"):]
		shirt = LoadTexture("https://api.brick-hill.com/v1/assets/get/" + texture)
	}

	pants := fauxgl.NewImageTexture(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	if pantsValue, ok := avatar.Items["pants"].(float64); ok && pantsValue != 0 {
		resp, err := http.Get(fmt.Sprintf("https://api.brick-hill.com/v1/assets/getPoly/1/%d", int(pantsValue)))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var data []map[string]string
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			panic(err)
		}

		texture := data[0]["texture"]
		texture = texture[len("asset://"):]
		pants = LoadTexture("https://api.brick-hill.com/v1/assets/get/" + texture)
	}

	mesh := LoadMeshFromURL("https://hawli.pages.dev/obj/Torso.obj")
	scene.AddObject(&fauxgl.Object{
		Mesh:    mesh,
		Color:   fauxgl.HexColor("777"),
		Texture: shirt,
	})

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/Head.obj")
	if faceValue, ok := avatar.Items["face"].(float64); ok && faceValue != 0 {
		resp, err := http.Get(fmt.Sprintf("https://api.brick-hill.com/v1/assets/getPoly/1/%d", int(faceValue)))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var data []map[string]string
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			panic(err)
		}

		if len(data) > 0 {
			texture := data[0]["texture"]
			texture = texture[len("asset://"):]
			scene.AddObject(&fauxgl.Object{
				Mesh:    mesh,
				Color:   fauxgl.HexColor(avatar.Colors["head"]),
				Texture: LoadTexture("https://api.brick-hill.com/v1/assets/get/" + texture),
			})
		}
	} else {
		// Load the default face texture
		scene.AddObject(&fauxgl.Object{
			Mesh:    mesh,
			Color:   fauxgl.HexColor(avatar.Colors["head"]),
			Texture: LoadTexture("https://hawli.pages.dev/Face.png"),
		})
	}

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/LeftArm.obj")
	scene.AddObject(&fauxgl.Object{
		Mesh:    mesh,
		Color:   fauxgl.HexColor(avatar.Colors["left_arm"]),
		Texture: shirt,
	})

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/LeftLeg.obj")
	scene.AddObject(&fauxgl.Object{
		Mesh:    mesh,
		Color:   fauxgl.HexColor(avatar.Colors["left_leg"]),
		Texture: pants,
	})

	if toolValue, ok := avatar.Items["tool"].(float64); ok && toolValue != 0 {
		mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/ArmHold.obj")
		scene.AddObject(&fauxgl.Object{
			Mesh:    mesh,
			Color:   fauxgl.HexColor(avatar.Colors["right_arm"]),
			Texture: shirt,
		})
		LoadItem(int(toolValue), scene)
	} else {
		mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/RightArm.obj")
		scene.AddObject(&fauxgl.Object{
			Mesh:    mesh,
			Color:   fauxgl.HexColor(avatar.Colors["right_arm"]),
			Texture: shirt,
		})
	}

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/RightLeg.obj")
	scene.AddObject(&fauxgl.Object{
		Mesh:    mesh,
		Color:   fauxgl.HexColor(avatar.Colors["right_leg"]),
		Texture: pants,
	})

	shader.AmbientColor = fauxgl.HexColor("AAA")
	shader.DiffuseColor = fauxgl.HexColor("777")
	shader.SpecularPower = 0

	newMatrix := scene.FitObjectsToScene(eye, center, up, fovy, aspect, near, far)
	shader.Matrix = newMatrix
	scene.Draw()

	outImg := context.Image()
	buf := new(bytes.Buffer)
	err = png.Encode(buf, outImg)
	if err != nil {
		fmt.Fprintln(out, "Error:", err)
	}

	resp := ImageResponse{
		Image: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}

	json.NewEncoder(out).Encode(resp)
}
