package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"os/exec"

	fdk "github.com/fnproject/fdk-go"
	"github.com/gofrs/uuid"
	fauxgl "github.com/hawl1/brickgl"
)

const (
	scale = 3
	fovy  = 50
	near  = 0.1
	far   = 1000
)

var (
	eye    = fauxgl.V(-0.75, 0.85, -2)
	center = fauxgl.V(0, 0.06, 0)
	up     = fauxgl.V(0, 1, 0)
	light  = fauxgl.V(0, 6, -4).Normalize()

	def = "{\"user_id\":13,\"items\":{\"face\":0,\"hats\":[20121,0,0,0,0],\"head\":0,\"tool\":0,\"pants\":0,\"shirt\":0,\"figure\":0,\"tshirt\":0},\"colors\":{\"head\":\"eab372\",\"torso\":\"85ad00\",\"left_arm\":\"eab372\",\"left_leg\":\"37302c\",\"right_arm\":\"eab372\",\"right_leg\":\"37302c\"}}"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(HandleRenderEvent))
}

type RenderEvent struct {
	AvatarJSON string
	Size       int
}

type ImageResponse struct {
	UUID  string
	Image string
}

func HandleRenderEvent(resp ImageResponse, err error) {
	if e.AvatarJSON == "" {
		e.AvatarJSON = def
	}
	var namespace uuid.UUID
	namespace, err = uuid.FromString(os.Getenv("THUMBNAIL_UUID_NAMESPACE"))
	if err != nil {
		return
	}

	aspect := float64(e.Size) / float64(e.Size)

	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	context := fauxgl.NewContext(e.Size, e.Size, scale, shader)
	scene := fauxgl.NewScene(context)

	fileUUID := uuid.NewV5(namespace, e.AvatarJSON).String()
	objFileName := fmt.Sprintf("./%s.obj", fileUUID)
	pngFileName := fmt.Sprintf("./%s.png", fileUUID)

	cmd := exec.Command(
		"/opt/bin/exporter/avatar-exporter",
		"-s", "export_avatar_rs.gd",
		"--json", e.AvatarJSON,
		"--obj-path", objFileName,
		"--png-path", pngFileName,
	)
	err = cmd.Run()
	if err != nil {
		return
	}

	tex, _ := fauxgl.LoadTexture(pngFileName)
	mesh, _ := fauxgl.LoadOBJ(objFileName)

	scene.AddObject(&fauxgl.Object{
		Texture: tex,
		Mesh:    mesh,
		Color:   fauxgl.HexColor("777"),
	})

	shader.AmbientColor = fauxgl.HexColor("AAA")
	shader.DiffuseColor = fauxgl.HexColor("777")
	shader.SpecularPower = 0

	newMatrix := scene.FitObjectsToScene(eye, center, up, fovy, aspect, near, far)
	shader.Matrix = newMatrix
	scene.Draw()

	outImg := context.Image()
	buf := new(bytes.Buffer)
	png.Encode(buf, outImg)

	// remove files to prevent tmp dir from filling
	os.Remove(objFileName)
	os.Remove(pngFileName)

	resp = ImageResponse{UUID: uuid.NewV5(namespace, buf.String()).String(), Image: base64.StdEncoding.EncodeToString(buf.Bytes())}

	return
}
