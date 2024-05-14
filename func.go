package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"net/http"

	"github.com/nfnt/resize"

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
	eye    = fauxgl.V(0.75, 0.85, 2)
	center = fauxgl.V(0, 0.06, 0)
	up     = fauxgl.V(0, 1, 0)
	light  = fauxgl.V(0, 6, 4).Normalize()
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
		avatarJSON = "{\"user_id\":13,\"items\":{\"face\":0,\"hats\":[20121,0,0,0,0],\"head\":0,\"tool\":6929,\"pants\":0,\"shirt\":368750,\"figure\":0,\"tshirt\":187683},\"colors\":{\"head\":\"f3b700\",\"torso\":\"929292\",\"left_arm\":\"f3b700\",\"left_leg\":\"e6e6e6\",\"right_arm\":\"f3b700\",\"right_leg\":\"e6e6e6\"}}"
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

	shirt := image.NewRGBA(image.Rect(0, 0, 1, 1))
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

		if len(data) == 0 {
			panic("API response does not contain expected data")
		}

		texture := data[0]["texture"]
		texture = texture[len("asset://"):]

		resp, err = http.Get("https://api.brick-hill.com/v1/assets/get/" + texture)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		img, _, err := image.Decode(bytes.NewReader(body))
		if err != nil {
			panic(err)
		}

		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		shirt = rgba
	}

	tshirt := image.NewRGBA(image.Rect(0, 0, 836, 836))
	if tshirtValue, ok := avatar.Items["tshirt"].(float64); ok && tshirtValue != 0 {
		resp, err := http.Get(fmt.Sprintf("https://api.brick-hill.com/v1/assets/getPoly/1/%d", int(tshirtValue)))
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

		resp, err = http.Get("https://api.brick-hill.com/v1/assets/get/" + texture)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		img, _, err := image.Decode(bytes.NewReader(body))
		if err != nil {
			panic(err)
		}

		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		tshirt = rgba
		resizedTshirt := resize.Resize(126, 126, tshirt, resize.Lanczos3)

		blankTemplate := image.NewRGBA(image.Rect(0, 0, 836, 836))

		draw.Draw(blankTemplate, image.Rect(183, 87, 309, 213), resizedTshirt, image.Point{}, draw.Src)

		tshirt = blankTemplate
	}

	combinedWidth := shirt.Bounds().Max.X
	combinedHeight := shirt.Bounds().Max.Y
	combined := image.NewRGBA(image.Rect(0, 0, combinedWidth, combinedHeight))

	draw.Draw(combined, shirt.Bounds(), shirt, image.Point{}, draw.Over)

	draw.Draw(combined, tshirt.Bounds(), tshirt, image.Point{}, draw.Over)

	var tshirtBuf bytes.Buffer
	if err := png.Encode(&tshirtBuf, combined); err != nil {
		panic(err)
	}

	encoded := base64.StdEncoding.EncodeToString(tshirtBuf.Bytes())

	fmt.Println(encoded)

	combinedShirt := fauxgl.NewImageTexture(combined)

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

	mesh := LoadMeshFromURL("https://hawli.pages.dev/lunarhill/Torso.obj")
	scene.AddObject(&fauxgl.Object{
		Mesh:    mesh,
		Color:   fauxgl.HexColor(avatar.Colors["torso"]),
		Texture: combinedShirt,
	})

	if headValue, ok := avatar.Items["head"].(float64); ok && headValue != 0 {
		resp, err := http.Get(fmt.Sprintf("https://api.brick-hill.com/v1/assets/getPoly/1/%d", int(headValue)))
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
			meshdat := data[0]["mesh"]
			meshdat = meshdat[len("asset://"):]
			mesh = LoadMeshFromURL("https://api.brick-hill.com/v1/assets/get/" + meshdat)
		}
	} else {
		mesh = LoadMeshFromURL("https://hawli.pages.dev/lunarhill/Head.obj")
	}

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
			Texture: LoadTexture("https://hawli.pages.dev/lunarhill/defaultFace.png"),
		})
	}

	scene.AddObject(&fauxgl.Object{
		Mesh:    LoadMeshFromURL("https://hawli.pages.dev/lunarhill/LeftArm.obj"),
		Color:   fauxgl.HexColor(avatar.Colors["left_arm"]),
		Texture: combinedShirt,
	})

	scene.AddObject(&fauxgl.Object{
		Mesh:    LoadMeshFromURL("https://hawli.pages.dev/lunarhill/LeftLeg.obj"),
		Color:   fauxgl.HexColor(avatar.Colors["left_leg"]),
		Texture: pants,
	})

	if toolValue, ok := avatar.Items["tool"].(float64); ok && toolValue != 0 {
		mesh = LoadMeshFromURL("https://hawli.pages.dev/lunarhill/ArmHold.obj")
		scene.AddObject(&fauxgl.Object{
			Mesh:    mesh,
			Color:   fauxgl.HexColor(avatar.Colors["right_arm"]),
			Texture: combinedShirt,
		})
		LoadItem(int(toolValue), scene)
	} else {
		mesh = LoadMeshFromURL("https://hawli.pages.dev/lunarhill/RightArm.obj")
		scene.AddObject(&fauxgl.Object{
			Mesh:    mesh,
			Color:   fauxgl.HexColor(avatar.Colors["right_arm"]),
			Texture: combinedShirt,
		})
	}

	mesh = LoadMeshFromURL("https://hawli.pages.dev/lunarhill/RightLeg.obj")
	scene.AddObject(&fauxgl.Object{
		Mesh:    mesh,
		Color:   fauxgl.HexColor(avatar.Colors["right_leg"]),
		Texture: pants,
	})

	hats, ok := avatar.Items["hats"].([]interface{})
	if ok {
		for _, hatValue := range hats {
			if hat, ok := hatValue.(float64); ok && hat != 0 {
				LoadItem(int(hat), scene)
			}
		}
	}

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
